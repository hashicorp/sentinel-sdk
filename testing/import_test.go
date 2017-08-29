package testing

import (
	"path/filepath"
	"testing"

	testingiface "github.com/mitchellh/go-testing-interface"
)

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

	// Test configuration
	t.Run("config", func(t *testing.T) {
		TestImport(t, TestImportCase{
			ImportPath: path,
			Config:     map[string]interface{}{"suffix": "??"},
			Source:     `main = subject.foo == "foo??"`,
		})
	})
}
