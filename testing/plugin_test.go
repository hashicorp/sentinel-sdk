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

func TestTestPlugin(t *testing.T) {
	path, err := filepath.Abs("testplugin")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Test good
	t.Run("success", func(t *testing.T) {
		TestPlugin(t, TestPluginCase{
			PluginPath: path,
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

		TestPlugin(&testingiface.RuntimeT{}, TestPluginCase{
			PluginPath: path,
			Source:     `main = subject.foo == "foo!"`,
		})
	})

	// Test runtime error
	t.Run("error", func(t *testing.T) {
		TestPlugin(&testingiface.RuntimeT{}, TestPluginCase{
			PluginPath: path,
			Source:     `main = rule { error("nope") }`,
			Error:      "nope",
		})
	})

	// Test runtime error w/ regular expression
	t.Run("error with regex", func(t *testing.T) {
		TestPlugin(&testingiface.RuntimeT{}, TestPluginCase{
			PluginPath: path,
			Source:     `main = rule { error("super 1337 error") }`,
			Error:      `/super \d+ error/`,
		})
	})

	// Test runtime error w/ errored regular expression
	t.Run("error with errored regex", func(t *testing.T) {
		// Use a defer to catch a panic that RuntimeT will throw. We can
		// detect the failure this way.
		defer func() {
			if e := recover(); e == nil {
				t.Fatal("should fail")
			}
		}()

		TestPlugin(&testingiface.RuntimeT{}, TestPluginCase{
			PluginPath: path,
			Source:     `main = rule { error("super 1337 error") }`,
			Error:      `/(super \d+ error/`,
		})
	})

	// Test configuration
	t.Run("config", func(t *testing.T) {
		TestPlugin(t, TestPluginCase{
			PluginPath: path,
			Config:     map[string]interface{}{"suffix": "??"},
			Source:     `main = subject.foo == "foo??"`,
		})
	})

	// Test globals
	t.Run("global", func(t *testing.T) {
		TestPlugin(t, TestPluginCase{
			PluginPath: path,
			Global:     map[string]interface{}{"value": "foo??"},
			Source:     `main = value == "foo??"`,
		})
	})

	// Test mocks
	t.Run("mock", func(t *testing.T) {
		TestPlugin(t, TestPluginCase{
			PluginPath: path,
			Mock: map[string]map[string]interface{}{
				"data": map[string]interface{}{
					"value": "foo??",
				},
			},
			Source: `import "data"; main = data.value == "foo??"`,
		})
	})

	t.Run("custom plugin name", func(t *testing.T) {
		TestPlugin(t, TestPluginCase{
			PluginPath: path,
			PluginName: "foo",
			Source:     `main = foo.bar is "bar!!"`,
		})
	})

	// TestDirectory helper
	t.Run("directory", func(t *testing.T) {
		TestPluginDir(t, "testdata/plugin-test-dir", func(tc *TestPluginCase) {
			tc.PluginPath = path
			tc.Global = map[string]interface{}{"exclamation": "!"}
		})
	})
}
