package testing

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/scanner"

	"github.com/mitchellh/go-testing-interface"
)

//go:generate go-bindata -nomemcopy -pkg=testing ./data/...

// importMap is the list of built import binaries keyed by import path.
// This import path should be canonicalized via ImportPath.
var importBuildDir string
var importMap = map[string]string{}
var importErr = map[string]error{}

// TestImportCase is a single test case for configuring TestImport.
type TestImportCase struct {
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

	// Mock is mocked import data
	Mock map[string]map[string]interface{}

	// ImportPath is the path to a Go package on your GOPATH containing
	// the import to test. If this is blank, the test case uses heuristics
	// to extract the GOPATH and use the current package for testing.
	// This package is expected to expose a "New" function which adheres to
	// the sdk/rpc.ImportFunc signature.
	//
	// This should usually be blank. This maximizes portability of the
	// import if it were to be forked or moved.
	//
	// For a given import path, the test binary will be built exactly once
	// per test run.
	ImportPath string

	// ImportName allows passing a custom name for the import to be used in
	// test cases. By default, the import is simply named "subject". The
	// import name is what is used within this policy's source to access
	// functionality provided by the import.
	ImportName string

	// A string containing any expected runtime error during evaluation. If
	// this field is non-empty, a runtime error is expected to occur, and the
	// Sentinel output is searched for the string given here. If a match is
	// found, the test passes. If it is not found the test will fail.
	Error string
}

// LoadTestImportCase is used to load a TestImportCase from a Sentinel policy
// file. The policy file is expected to begin with comments in the first
// line. Certain configuration options are supported in the comment body. The
// following is a completely valid example:
//
//     // config: {"option1": "value1"}
//     // error: failed to do the thing
//     main = rule { true }
//
// The above would load a TestImport case using the specified options. The
// config is loaded as a JSON string and unmarshaled into the Config field.
// The error field is loaded as a string into the Error field.
//
// This makes boilerplate very simple for a large number of Sentinel tests,
// and allows an entire test to be captured neatly into a single file which
// also happens to be the policy being tested.
func LoadTestImportCase(t testing.T, path string) TestImportCase {
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
		content := strings.TrimPrefix(raw, "// ")

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
			t.Fatalf("unsupported magic comment: %q", parts[0])
		}
	}

	if _, err := fh.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	policyBytes, err := ioutil.ReadAll(fh)
	if err != nil {
		t.Fatal(err)
	}

	tc := TestImportCase{
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

func TestDirectory(t testing.T, path string, customize func(*TestImportCase)) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}

	cases := make(map[string]TestImportCase)
	for _, fi := range files {
		// Allow the directory to be structured.
		if fi.IsDir() {
			continue
		}

		// Load the sentinel file and parse it.
		fp := filepath.Join(path, fi.Name())
		tc := LoadTestImportCase(t, fp)

		// If a customization function was provided, execute it.
		if customize != nil {
			customize(&tc)
		}

		// Add the test to the set.
		cases[fi.Name()] = tc
	}

	for name, tc := range cases {
		// This would be nice but the interface doesn't have Run() yet.
		//t.Run(name, func(t testing.T) {
		//	TestImport(t, tc)
		//})
		t.Log(name)
		TestImport(t, tc)
	}
}

// Clean cleans any temporary files created. This should always be called
// at the end of any set of import tests.
func Clean() {
	// Delete our build directory
	if importBuildDir != "" {
		os.RemoveAll(importBuildDir)
	}

	// Reset all globals
	importMap = map[string]string{}
	importErr = map[string]error{}
}

// TestImport tests that a sdk.Import implementation works as expected.
func TestImport(t testing.T, c TestImportCase) {
	// Infer the path
	path, err := ImportPath(c.ImportPath)
	if err != nil {
		t.Fatalf("error inferring GOPATH: %s", err)
	}

	// If we already errored building this, report it
	if err, ok := importErr[path]; ok {
		t.Fatalf("error building import: %s", err)
	}

	// Get the path to the built import, or build it
	binaryPath, ok := importMap[path]
	if !ok {
		binaryPath = buildImport(t, path)
	}

	// Build the full source which requires importing the subject
	src := `import "subject"`
	if c.ImportName != "" {
		src += " as " + c.ImportName
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
			if !strings.Contains(string(output), c.Error) {
				t.Fatalf("expected error %q not found:\n\n%s",
					c.Error, string(output))
			}
		} else {
			t.Fatalf("error executing test. output:\n\n%s", string(output))
		}
	} else if c.Error != "" {
		t.Fatalf("expected error %q but policy passed", c.Error)
	}
}

// importPathModule determines the import path when modules are
// enabled, through the use of "go list".
//
// The working directory is set to dir, if supplied.
func importPathModule(dir string) (string, error) {
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

// ImportPath attempts to infer the import path based on the GOPATH
// environment variable and the directory.
func ImportPath(dir string) (string, error) {
	if isUsingModules() {
		return importPathModule(dir)
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

// buildImport compiles the import binary with the given Go import path.
// The path to the completed binary is inserted into the global importMap.
func buildImport(t testing.T, path string) string {
	log.Printf("Building binary: %s", path)

	// Create the main.go
	main := bytes.Replace(
		MustAsset("data/main.go.tpl"),
		[]byte("PATH"), []byte(path), -1)

	// If we don't have a build dir, make one
	if importBuildDir == "" {
		// Create the directory to compile this
		wd, err := os.Getwd()
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		td, err := ioutil.TempDir(wd, "sentinel-sdk")
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		importBuildDir = td
	}

	// Create the build dir for this import
	td, err := ioutil.TempDir(importBuildDir, "sentinel-sdk")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Write the file
	if err := ioutil.WriteFile(filepath.Join(td, "main.go"), main, 0644); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Build
	cmd := exec.Command("go", "build", "-o", "import-test")
	cmd.Dir = td
	output, err := cmd.CombinedOutput()
	if err != nil {
		importErr[path] = err
		t.Fatalf("err building the test binary. output:\n\n%s", string(output))
	}

	// Record it
	importMap[path] = filepath.Join(td, "import-test")
	log.Printf("Import binary built at: %s", importMap[path])
	return importMap[path]
}
