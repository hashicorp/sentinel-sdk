// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rpc

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"

	sdk "github.com/hashicorp/sentinel-sdk"
)

func TestPlugin_gRPC_configure(t *testing.T) {
	// Create a mock object
	pluginMock := new(sdk.MockPlugin)
	pluginMock.On("Configure",
		map[string]interface{}{"key": int64(42)}).
		Return(nil)

	obj, closer := testPluginServeGRPC(t, pluginMock)
	defer closer()

	// Get
	err := obj.Configure(map[string]interface{}{"key": 42})
	pluginMock.AssertExpectations(t)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPlugin_gRPC_get(t *testing.T) {
	cases := []struct {
		Name     string
		Requests []*sdk.GetReq
		Results  []*sdk.GetResult
	}{
		{
			Name: "basic",
			Requests: []*sdk.GetReq{
				{
					KeyId: 42,
					Keys: []sdk.GetKey{
						{
							Key:  "foo",
							Args: []interface{}{"foo", int64(42)},
						},
					},
					Context: map[string]interface{}{
						"_type": "SomeNamespace",
						"data": map[string]interface{}{
							"string": "foo",
							"number": int64(0),
						},
					},
				},
			},
			Results: []*sdk.GetResult{
				{
					KeyId: 42,
					Keys:  []string{"key"},
					Value: "bar",
					Context: map[string]interface{}{
						"_type": "SomeNamespace",
						"data": map[string]interface{}{
							"string": "bar",
							"number": int64(1),
						},
					},
					Callable: true,
				},
			},
		},
		{
			Name: "niladic",
			Requests: []*sdk.GetReq{
				{
					KeyId: 42,
					Keys: []sdk.GetKey{
						{
							Key:  "foo",
							Args: []interface{}{},
						},
					},
				},
			},
			Results: []*sdk.GetResult{
				{
					KeyId: 42,
					Keys:  []string{"key"},
					Value: "bar",
				},
			},
		},
		{
			Name: "not a function",
			Requests: []*sdk.GetReq{
				{
					KeyId: 42,
					Keys: []sdk.GetKey{
						{
							Key: "foo",
						},
					},
				},
			},
			Results: []*sdk.GetResult{
				{
					KeyId: 42,
					Keys:  []string{"key"},
					Value: "bar",
				},
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			pluginMock := new(sdk.MockPlugin)
			pluginMock.On("Configure", map[string]interface{}{}).Return(nil)
			pluginMock.On("Get",
				mock.MatchedBy(func(reqs []*sdk.GetReq) bool {
					if len(tc.Requests) != len(reqs) {
						return false
					}

					for i := range tc.Requests {
						tc.Requests[i].ExecDeadline = reqs[i].ExecDeadline
					}

					return reflect.DeepEqual(tc.Requests, reqs)
				})).Return(tc.Results, nil)

			obj, closer := testPluginServeGRPC(t, pluginMock)
			defer closer()

			// We need to configure first
			if err := obj.Configure(nil); err != nil {
				t.Fatalf("err: %s", err)
			}

			// Get
			results, err := obj.Get(tc.Requests)
			pluginMock.AssertExpectations(t)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if !reflect.DeepEqual(results, tc.Results) {
				t.Fatalf("expected %#v, got %#v", tc.Results, results)
			}
		})
	}
}
