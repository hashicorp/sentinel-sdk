package rpc

import (
	"testing"

	goplugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/sentinel-sdk"
)

func TestImport_impl(t *testing.T) {
	var _ goplugin.Plugin = new(ImportPlugin)
	var _ goplugin.GRPCPlugin = new(ImportPlugin)
}

func testImportServeGRPC(t *testing.T, o sdk.Import) (sdk.Import, func()) {
	client, _ := goplugin.TestPluginGRPCConn(t, pluginMap(&ServeOpts{
		ImportFunc: testImportFixed(o),
	}))

	// Request the Import
	raw, err := client.Dispense(ImportPluginName)
	if err != nil {
		client.Close()
		t.Fatalf("err: %s", err)
	}

	return raw.(sdk.Import), func() {
		client.Close()
	}
}

func testImportFixed(p sdk.Import) ImportFunc {
	return func() sdk.Import {
		return p
	}
}
