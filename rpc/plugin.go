// Copyright IBM Corp. 2017, 2025
// SPDX-License-Identifier: MPL-2.0

package rpc

import (
	"context"

	"google.golang.org/grpc"

	goplugin "github.com/hashicorp/go-plugin"

	sdk "github.com/hashicorp/sentinel-sdk"
	proto "github.com/hashicorp/sentinel-sdk/proto/go"
)

// Plugin is the goplugin.Plugin implementation to serve sdk.Plugin.
type Plugin struct {
	goplugin.NetRPCUnsupportedPlugin

	F func() sdk.Plugin
}

func (p *Plugin) GRPCServer(_ *goplugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterPluginServer(s, &PluginGRPCServer{F: p.F})
	return nil
}

func (p *Plugin) GRPCClient(_ context.Context, _ *goplugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &PluginGRPCClient{Client: proto.NewPluginClient(c)}, nil
}
