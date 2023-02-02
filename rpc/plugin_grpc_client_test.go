// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rpc

import (
	"io"
	"testing"

	sdk "github.com/hashicorp/sentinel-sdk"
)

func TestPluginGRPCClient_impl(t *testing.T) {
	var _ sdk.Plugin = new(PluginGRPCClient)
	var _ io.Closer = new(PluginGRPCClient)
}
