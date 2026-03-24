// Copyright IBM Corp. 2017, 2025
// SPDX-License-Identifier: MPL-2.0

package testing

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"text/scanner"

	"github.com/mitchellh/go-testing-interface"
)

//go:embed data
var content embed.FS

// pluginMap is the list of built plugin binaries keyed by plugin path.
// This plugin path should be canonicalized via PluginPath.
var pluginBuildDir string
var pluginMap = map[string]string{}
var pluginErr = map[string]error{}

// TestPluginCase is a single test case for configuring TestPlugin.
type TestPluginCase struct {
	// Source is a policy to execute. This should be a full program ending
	// in `main = ` and an assignment. For example `main = subject.foo`.
	Source string

	// This is the configuration that will be sent to the plugin. This
	// must serialize to JSON since the JSON will be used to pass the
	// configuration.
	Config map[string]interface{}

	// This is extra data to inject into the global scope of the policy
	// execution
	Global map[string]interface{}

	// Mock is mocked plugin data
	Mock map[string]map[string]interface{}

	// PluginPath is the path to a Go package on your GOPATH containing
	// the plugin to test. If this is blank, the test case uses heuristics
	// to extract the GOPATH and use the current package for testing.
	// This package is expected to expose a "New" function which adheres to
	// the sdk/rpc.PluginFunc signature.
	//
	// This should usually be blank. This maximizes portability of the
	// plugin if it were to be forked or moved.
	//
	// For a given plugin path, the test binary will be built exactly once
	// per test run.
	PluginPath string

	// PluginName allows passing a custom name for the plugin to be used in
	// test cases. By default, the plugin is simply named "subject". The
	// plugin name is what is used within this policy's source to access
	// functionality provided by the plugin.
	PluginName string

	// A string containing any expected runtime error during evaluation. If
	// this field is non-empty, a runtime error is expected to occur, and
	// the Sentinel output is searched for the string given here. If the
	// output contains the string, the test passes. If it does not contain
	// the string, the test will fail.
	//
	// More advanced matches can be done with regular expression patterns.
	// If the Error string is delimited by slashes (/), the string is
	// compiled as a regular expression and the Sentinel output is matched
	// against the resulting pattern. If a match is found, the test passes.
	// If it does not match, the tests will fail.
	Error string
}

// LoadTestPluginCase is used to load a TestPluginCase from a Sentinel policy
// file. Certain test case pragmas are supported in the top-most comment body.
// The following is a completely valid example:
//
//	//config: {"option1": "value1"}
//	//error: failed to do the thing
//	main = rule { true }
//
// The above would load a TestPlugin case using the specified options. The
// config is loaded as a JSON string and unmarshaled into the Config field.
// The error field is loaded as a string into the Error field. Pragmas *must*
// be at the very top of the file, starting at line one. When a non-pragma
// line is encountered, parsing will end and any further pragmas are discarded.
//
// This makes boilerplate very simple for a large number of Sentinel tests,
// and allows an entire test to be captured neatly into a single file which
// also happens to be the policy being tested.
func LoadTestPluginCase(t testing.T, path string) TestPluginCase {
	fh, err := os.Open(path)
	if err != nil {
		t.Fatalf("error opening policy: %v", err)
	}
	defer fh.Close()

	var s scanner.Scanner
	s.Init(fh)
	s.Mode ^= scanner.SkipComments

	var errMatch string
	var configStr string

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		raw := s.TokenText()
		content := strings.TrimPrefix(raw, "//")

		// Make sure we are still in the top comments.
		if raw == content {
			break
		}

		parts := strings.SplitN(content, ":", 2)
		if len(parts) < 2 {
			continue
		}

		switch parts[0] {
		case "error":
			errMatch = strings.TrimSpace(parts[1])
		case "config":
			configStr = strings.TrimSpace(parts[1])
		default:
			break // Require magic comments to be at the top.
		}
	}

	if _, err := fh.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	policyBytes, err := ioutil.ReadAll(fh)
	if err != nil {
		t.Fatal(err)
	}

	tc := TestPluginCase{
		Source: string(policyBytes),
		Error:  errMatch,
	}

	if configStr != "" {
		tc.Config = make(map[string]interface{})
		if err := json.Unmarshal([]byte(configStr), &tc.Config); err != nil {
			t.Fatalf("error decoding configuration: %v", err)
		}
	}

	return tc
}

// TestPluginDir iterates over files in a directory, calls
// LoadTestPluginCase on each file suffixed with ".sentinel", and executes all
// of the plugin tests.
func TestPluginDir(t testing.T, path string, customize func(*TestPluginCase)) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}

	cases := make(map[string]TestPluginCase)
	for _, fi := range files {
		// Allow the directory to be structured.
		if fi.IsDir() {
			continue
		}

		// Only use files ending with '.sentinel'
		if !strings.HasSuffix(fi.Name(), ".sentinel") {
			continue
		}

		// Load the sentinel file and parse it.
		fp := filepath.Join(path, fi.Name())
		tc := LoadTestPluginCase(t, fp)

		// If a customization function was provided, execute it.
		if customize != nil {
			customize(&tc)
		}

		// Add the test to the set.
		cases[fi.Name()] = tc
	}

	// Run all of the tests.
	for file, tc := range cases {
		// The testing interface (mitchellh/go-testing-interface) doesn't
		// support a t.Run(), and adding context about which policy is failing
		// to the error is obtuse otherwise, so we'll just log the policy file
		// name here to give that context to the developer.
		t.Logf("Checking %s ...", file)
		TestPlugin(t, tc)
	}
}

// Clean cleans any temporary files created. This should always be called
// at the end of any set of plugin tests.
func Clean() {
	// Delete our build directory
	if pluginBuildDir != "" {
		os.RemoveAll(pluginBuildDir)
	}

	// Reset all globals
	pluginMap = map[string]string{}
	pluginErr = map[string]error{}
}

// TestPlugin tests that a sdk.Plugin implementation works as expected.
func TestPlugin(t testing.T, c TestPluginCase) {
	// Infer the path
	path, err := PluginPath(c.PluginPath)
	if err != nil {
		t.Fatalf("error inferring GOPATH: %s", err)
	}

	// If we already errored building this, report it
	if err, ok := pluginErr[path]; ok {
		t.Fatalf("error building plugin: %s", err)
	}

	// Get the path to the built plugin, or build it
	binaryPath, ok := pluginMap[path]
	if !ok {
		binaryPath = buildPlugin(t, path)
	}

	// Build the full source which requires importing the subject
	src := `import "subject"`
	if c.PluginName != "" {
		src += " as " + c.PluginName
	}
	src += "\n\n" + c.Source

	// Make the test directory where we'll run the test.
	td, err := ioutil.TempDir("", "sentinel-sdk")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(td)

	// Write the policy
	policyPath := filepath.Join(td, "policy.sentinel")
	if err := ioutil.WriteFile(policyPath, []byte(src), 0644); err != nil {
		t.Fatalf("error writing policy: %s", err)
	}

	// Write the configuration to execute
	configPath := filepath.Join(td, "config.json")
	config, err := json.MarshalIndent(map[string]interface{}{
		"imports": map[string]interface{}{
			"subject": map[string]interface{}{
				"path":   binaryPath,
				"config": c.Config,
			},
		},
		"global": c.Global,
		"mock":   c.Mock,
	}, "", "\t")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if err := ioutil.WriteFile(configPath, config, 0644); err != nil {
		t.Fatalf("error writing config: %s", err)
	}

	// Execute Sentinel
	cmd := exec.Command("sentinel", "apply", "-config", configPath, policyPath)
	cmd.Dir = td
	output, err := cmd.CombinedOutput()
	if err != nil {
		if c.Error != "" {
			if c.Error[:1]+c.Error[len(c.Error)-1:] == "//" {
				pattern := c.Error[1 : len(c.Error)-1]
				exp, err := regexp.Compile(pattern)
				if err != nil {
					t.Fatalf("error compiling expected error pattern: %s", err)
				}
				if !exp.Match(output) {
					t.Fatalf("the resulting error does not match the expected pattern: %s\n\nError output:\n\n%s",
						c.Error, string(output))
				}
			} else {
				if !strings.Contains(string(output), c.Error) {
					t.Fatalf("resulting error does not contain %q\n\nError output:\n\n%s",
						c.Error, string(output))
				}
			}
		} else {
			t.Fatalf("error executing test. output:\n\n%s", string(output))
		}
	} else if c.Error != "" {
		t.Fatalf("expected error %q but policy passed", c.Error)
	}
}

// pluginPathModule determines the plugin path when modules are
// enabled, through the use of "go list".
//
// The working directory is set to dir, if supplied.
func pluginPathModule(dir string) (string, error) {
	cmd := exec.Command("go", "list")
	if dir != "" {
		wd, err := filepath.Abs(dir)
		if err != nil {
			return "", err
		}

		cmd.Dir = wd
	}

	out, err := cmd.Output()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			log.Println(string(e.Stderr))
		}

		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// isUsingModules checks to see if modules are enabled on the working
// repository.
func isUsingModules() bool {
	if err := exec.Command("go", "list", "-m").Run(); err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			// Log stderr if we have it
			log.Println(strings.TrimSpace(string(e.Stderr)))
		}

		return false
	}

	return true
}

// PluginPath attempts to infer the plugin path based on the GOPATH
// environment variable and the directory.
func PluginPath(dir string) (string, error) {
	if isUsingModules() {
		return pluginPathModule(dir)
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return "", errors.New("no GOPATH set")
	}

	// Append src to the GOPATH since we're looking for a source path
	gopath = filepath.Join(gopath, "src")

	// Create the absolute path for the directory
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("error expanding %q: %s", dir, err)
	}

	// The directory should have the gopath as a prefix if its within the GOPATH
	if !strings.HasPrefix(dir, gopath) {
		return "", fmt.Errorf("Directory %q doesn't appear in GOPATH %q", dir, gopath)
	}

	// Trim the gopath from the front. If we have a slash remaining, trim that
	path := strings.TrimPrefix(dir, gopath)
	if path[0] == '/' {
		path = path[1:]
	}

	return path, nil
}

// buildPlugin compiles the plugin binary with the given Go import path.
// The path to the completed binary is inserted into the global pluginMap.
func buildPlugin(t testing.T, path string) string {
	log.Printf("Building binary: %s", path)

	tpl, err := content.ReadFile("data/main.go.tpl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	// Create the main.go
	main := bytes.Replace(
		tpl,
		[]byte("PATH"), []byte(path), -1)

	// If we don't have a build dir, make one
	if pluginBuildDir == "" {
		// Create the directory to compile this
		wd, err := os.Getwd()
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		td, err := ioutil.TempDir(wd, "sentinel-sdk")
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		pluginBuildDir = td
	}

	// Create the build dir for this plugin
	td, err := ioutil.TempDir(pluginBuildDir, "sentinel-sdk")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Write the file
	if err := ioutil.WriteFile(filepath.Join(td, "main.go"), main, 0644); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Build.  Note that when running on Windows systems the
	// plugin will need an .EXE extension
	buildOutput := "plugin-test"
	if isWindows() {
		buildOutput += ".exe"
	}

	cmd := exec.Command("go", "build", "-o", buildOutput)
	cmd.Dir = td
	output, err := cmd.CombinedOutput()
	if err != nil {
		pluginErr[path] = err
		t.Fatalf("err building the test binary. output:\n\n%s", string(output))
	}

	// Record it
	pluginMap[path] = filepath.Join(td, buildOutput)
	log.Printf("Plugin binary built at: %s", pluginMap[path])
	return pluginMap[path]
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}
