package testing

import (
	"os"
	"path/filepath"
	"testing"

	testingiface "github.com/mitchellh/go-testing-interface"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	Clean()
	os.Exit(exitCode)
}

func TestTestImport(t *testing.T) {
	path, err := filepath.Abs("testimport")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Test good
	t.Run("success", func(t *testing.T) {
		TestImport(t, TestImportCase{
			ImportPath: path,
			Source:     `main = subject.foo == "foo!!"`,
		})
	})

	// TestBad
	t.Run("failure", func(t *testing.T) {
		// Use a defer to catch a panic that RuntimeT will throw. We can
		// detect the failure this way.
		defer func() {
			if e := recover(); e == nil {
				t.Fatal("should fail")
			}
		}()

		TestImport(&testingiface.RuntimeT{}, TestImportCase{
			ImportPath: path,
			Source:     `main = subject.foo == "foo!"`,
		})
	})

	// Test runtime error
	t.Run("error", func(t *testing.T) {
		TestImport(&testingiface.RuntimeT{}, TestImportCase{
			ImportPath: path,
			Source:     `main = rule { error("nope") }`,
			Error:      "nope",
		})
	})

	// Test configuration
	t.Run("config", func(t *testing.T) {
		TestImport(t, TestImportCase{
			ImportPath: path,
			Config:     map[string]interface{}{"suffix": "??"},
			Source:     `main = subject.foo == "foo??"`,
		})
	})

	// Test globals
	t.Run("global", func(t *testing.T) {
		TestImport(t, TestImportCase{
			ImportPath: path,
			Global:     map[string]interface{}{"value": "foo??"},
			Source:     `main = value == "foo??"`,
		})
	})

	// Test mocks
	t.Run("mock", func(t *testing.T) {
		TestImport(t, TestImportCase{
			ImportPath: path,
			Mock: map[string]map[string]interface{}{
				"data": map[string]interface{}{
					"value": "foo??",
				},
			},
			Source: `import "data"; main = data.value == "foo??"`,
		})
	})

}
