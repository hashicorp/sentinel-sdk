package rpc

import (
	"fmt"
	"math"

	"github.com/hashicorp/sentinel-sdk"
	"github.com/hashicorp/sentinel-sdk/encoding"
	"github.com/hashicorp/sentinel-sdk/proto/go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// ImportGRPCClient is a gRPC server for Imports.
type ImportGRPCClient struct {
	Client proto.ImportClient

	instanceId uint64
}

func (m *ImportGRPCClient) Close() error {
	if m.instanceId > 0 {
		_, err := m.Client.Close(context.Background(), &proto.Close_Request{
			InstanceId: m.instanceId,
		})
		return err
	}

	return nil
}

func (m *ImportGRPCClient) Configure(config map[string]interface{}) error {
	v, err := encoding.GoToValue(config)
	if err != nil {
		return fmt.Errorf("config couldn't be encoded to plugin: %s", err)
	}

	resp, err := m.Client.Configure(context.Background(), &proto.Configure_Request{
		Config: v,
	})
	if err != nil {
		return err
	}

	m.instanceId = resp.InstanceId
	return nil
}

func (m *ImportGRPCClient) Get(rawReqs []*sdk.GetReq) ([]*sdk.GetResult, error) {
	reqs := make([]*proto.Get_Request, 0, len(rawReqs))
	for _, req := range rawReqs {
		// Request keys
		keys := make([]*proto.Get_Request_Key, len(req.Keys))
		for i, reqKey := range req.Keys {
			keys[i] = &proto.Get_Request_Key{Key: reqKey.Key}
			if reqKey.Args != nil {
				keys[i].Call = true
				keys[i].Args = make([]*proto.Value, len(reqKey.Args))
				for j, raw := range reqKey.Args {
					v, err := encoding.GoToValue(raw)
					if err != nil {
						return nil, err
					}

					keys[i].Args[j] = v
				}
			}
		}

		// Request context
		var reqCtx map[string]*proto.Value
		if req.Context != nil {
			reqCtx = make(map[string]*proto.Value)
			for k, raw := range req.Context {
				v, err := encoding.GoToValue(raw)
				if err != nil {
					return nil, err
				}

				reqCtx[k] = v
			}
		}

		reqs = append(reqs, &proto.Get_Request{
			InstanceId:   m.instanceId,
			ExecId:       req.ExecId,
			ExecDeadline: uint64(req.ExecDeadline.Unix()),
			Keys:         keys,
			KeyId:        req.KeyId,
			Context:      reqCtx,
		})
	}

	resp, err := m.Client.Get(
		context.Background(),
		&proto.Get_MultiRequest{
			Requests: reqs,
		},
		grpc.MaxRecvMsgSizeCallOption{MaxRecvMsgSize: math.MaxInt32},
		grpc.MaxSendMsgSizeCallOption{MaxSendMsgSize: math.MaxInt32},
	)
	if err != nil {
		return nil, err
	}

	results := make([]*sdk.GetResult, 0, len(resp.Responses))
	for _, resp := range resp.Responses {
		v, err := encoding.ValueToGo(resp.Value, nil)
		if err != nil {
			return nil, err
		}

		// Response context
		var resCtx map[string]interface{}
		if resp.Context != nil {
			resCtx = make(map[string]interface{})
			for k, raw := range resp.Context {
				v, err := encoding.ValueToGo(raw, nil)
				if err != nil {
					return nil, fmt.Errorf("error converting context value for key %q: %s", k, err)
				}

				resCtx[k] = v
			}
		}

		results = append(results, &sdk.GetResult{
			KeyId:    resp.KeyId,
			Keys:     resp.Keys,
			Value:    v,
			Context:  resCtx,
			Callable: resp.Callable,
		})
	}

	return results, nil
}
