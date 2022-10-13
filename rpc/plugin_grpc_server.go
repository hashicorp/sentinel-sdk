package rpc

import (
	"fmt"
	"io"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	sdk "github.com/hashicorp/sentinel-sdk"
	"github.com/hashicorp/sentinel-sdk/encoding"
	proto "github.com/hashicorp/sentinel-sdk/proto/go"
	"golang.org/x/net/context"
)

// PluginGRPCServer is a gRPC server for Plugins.
type PluginGRPCServer struct {
	F func() sdk.Plugin

	// instanceId is the current instance ID. This should be modified
	// with sync/atomic.
	instanceId    uint64
	instances     map[uint64]sdk.Plugin
	instancesLock sync.RWMutex
}

func (m *PluginGRPCServer) Close(
	ctx context.Context, v *proto.Close_Request) (*proto.Empty, error) {
	// Get the plugin and remove it immediately
	m.instancesLock.Lock()
	impt, ok := m.instances[v.InstanceId]
	delete(m.instances, v.InstanceId)
	m.instancesLock.Unlock()

	// If we have it, attempt to call Close on the plugin if it is
	// a closer.
	if ok {
		if c, ok := impt.(io.Closer); ok {
			c.Close()
		}
	}

	return &proto.Empty{}, nil
}

func (m *PluginGRPCServer) Configure(
	ctx context.Context, v *proto.Configure_Request) (*proto.Configure_Response, error) {
	// Build the configuration
	var config map[string]interface{}
	configRaw, err := encoding.ValueToGo(v.Config, reflect.TypeOf(config))
	if err != nil {
		return nil, fmt.Errorf("error converting config: %s", err)
	}
	config = configRaw.(map[string]interface{})

	// Configure is called once to configure a new plugin. Allocate the plugin.
	impt := m.F()

	// Call configure
	if err := impt.Configure(config); err != nil {
		return nil, err
	}

	// We have to allocate a new instance ID.
	id := atomic.AddUint64(&m.instanceId, 1)

	// Put the plugin into the store
	m.instancesLock.Lock()
	if m.instances == nil {
		m.instances = make(map[uint64]sdk.Plugin)
	}
	m.instances[id] = impt
	m.instancesLock.Unlock()

	// Configure the plugin
	return &proto.Configure_Response{
		InstanceId: id,
	}, nil
}

func (m *PluginGRPCServer) Get(
	ctx context.Context, v *proto.Get_MultiRequest) (*proto.Get_MultiResponse, error) {
	// Build the mapping of requests by instance ID. Then we can make the
	// calls for each proper instance easily.
	requestsById := make(map[uint64][]*sdk.GetReq)
	for _, req := range v.Requests {
		// Request keys
		keys := make([]sdk.GetKey, len(req.Keys))
		for i, reqKey := range req.Keys {
			keys[i] = sdk.GetKey{Key: reqKey.Key}
			if reqKey.Call {
				keys[i].Args = make([]interface{}, len(reqKey.Args))
				for j, arg := range reqKey.Args {
					obj, err := encoding.ValueToGo(arg, nil)
					if err != nil {
						return nil, fmt.Errorf("error converting arg %d: %s", i, err)
					}

					keys[i].Args[j] = obj
				}
			}
		}

		// Object context
		var reqCtx map[string]interface{}
		if req.Context != nil {
			reqCtx = make(map[string]interface{})
			for k, raw := range req.Context {
				v, err := encoding.ValueToGo(raw, nil)
				if err != nil {
					return nil, fmt.Errorf("error converting context value for key %q: %s", k, err)
				}

				reqCtx[k] = v
			}
		}

		getReq := &sdk.GetReq{
			ExecId:       req.ExecId,
			ExecDeadline: time.Unix(int64(req.ExecDeadline), 0),
			Keys:         keys,
			KeyId:        req.KeyId,
			Context:      reqCtx,
		}

		requestsById[req.InstanceId] = append(requestsById[req.InstanceId], getReq)
	}

	responses := make([]*proto.Get_Response, 0, len(v.Requests))
	for id, reqs := range requestsById {
		m.instancesLock.RLock()
		impt, ok := m.instances[id]
		m.instancesLock.RUnlock()
		if !ok {
			return nil, fmt.Errorf("unknown instance ID given: %d", id)
		}

		results, err := impt.Get(reqs)
		if err != nil {
			return nil, err
		}

		for _, result := range results {
			// Return value
			v, err := encoding.GoToValue(result.Value)
			if err != nil {
				return nil, err
			}

			// Return context
			var resCtx map[string]*proto.Value
			if result.Context != nil {
				resCtx = make(map[string]*proto.Value)
				for k, raw := range result.Context {
					v, err := encoding.GoToValue(raw)
					if err != nil {
						return nil, err
					}

					resCtx[k] = v
				}
			}

			responses = append(responses, &proto.Get_Response{
				InstanceId: id,
				KeyId:      result.KeyId,
				Keys:       result.Keys,
				Value:      v,
				Context:    resCtx,
				Callable:   result.Callable,
			})
		}
	}

	return &proto.Get_MultiResponse{Responses: responses}, nil
}
