// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rpc

import (
	"testing"

	goplugin "github.com/hashicorp/go-plugin"
	sdk "github.com/hashicorp/sentinel-sdk"
)

func TestPlugin_impl(t *testing.T) {
	var _ goplugin.Plugin = new(Plugin)
	var _ goplugin.GRPCPlugin = new(Plugin)
}

func testPluginServeGRPC(t *testing.T, o sdk.Plugin) (sdk.Plugin, func()) {
	client, _ := goplugin.TestPluginGRPCConn(t, pluginMap(&ServeOpts{
		PluginFunc: testPluginFixed(o),
	}))

	// Request the Plugin
	raw, err := client.Dispense(PluginName)
	if err != nil {
		client.Close()
		t.Fatalf("err: %s", err)
	}

	return raw.(sdk.Plugin), func() {
		client.Close()
	}
}

func testPluginFixed(p sdk.Plugin) PluginFunc {
	return func() sdk.Plugin {
		return p
	}
}
