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

	"github.com/mitchellh/go-testing-interface"
)

//go:generate go-bindata -nomemcopy -pkg=testing ./data/...

// importMap is the list of built import binaries keyed by import path.
// This import path should be canonicalized via ImportPath.
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
	src := "import \"subject\"\n\n" + c.Source

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
		t.Fatalf("error executing test. output:\n\n%s", string(output))
	}
}

// ImportPath attempts to infer the import path based on the GOPATH
// environment variable and the directory.
func ImportPath(dir string) (string, error) {
	// Get the GOPATH
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

	// Create the directory to compile this
	td, err := ioutil.TempDir("", "sentinel-sdk")
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
