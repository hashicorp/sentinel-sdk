// Package testplugin contains a test plugin that the testing package uses
// for unit tests. This plugin should not be actually used for anything.
package testplugin

import (
	sdk "github.com/hashicorp/sentinel-sdk"
	"github.com/hashicorp/sentinel-sdk/framework"
)

// New creates a new Plugin.
func New() sdk.Plugin {
	return &framework.Plugin{
		Root: &root{},
	}
}

type root struct {
	suffix string
}

// framework.Root impl.
func (m *root) Configure(raw map[string]interface{}) error {
	if v, ok := raw["suffix"]; ok {
		m.suffix = v.(string)
	}

	return nil
}

// framework.Namespace impl.
func (m *root) Get(key string) (interface{}, error) {
	suffix := m.suffix
	if suffix == "" {
		suffix = "!!"
	}

	return key + suffix, nil
}
