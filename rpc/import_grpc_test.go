package rpc

import (
	"reflect"
	"testing"

	sdk "github.com/hashicorp/sentinel-sdk"
	"github.com/stretchr/testify/mock"
)

func TestImport_gRPC_configure(t *testing.T) {
	// Create a mock object
	importMock := new(sdk.MockImport)
	importMock.On("Configure",
		map[string]int64{"key": int64(42)}).
		Return(nil)

	obj, closer := testImportServeGRPC(t, importMock)
	defer closer()

	// Get
	err := obj.Configure(map[string]interface{}{"key": 42})
	importMock.AssertExpectations(t)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestImport_gRPC_get(t *testing.T) {
	cases := []struct {
		Name     string
		Requests []*sdk.GetReq
		Results  []*sdk.GetResult
	}{
		{
			Name: "basic",
			Requests: []*sdk.GetReq{
				&sdk.GetReq{
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
				&sdk.GetResult{
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
				&sdk.GetReq{
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
				&sdk.GetResult{
					KeyId: 42,
					Keys:  []string{"key"},
					Value: "bar",
				},
			},
		},
		{
			Name: "not a function",
			Requests: []*sdk.GetReq{
				&sdk.GetReq{
					KeyId: 42,
					Keys: []sdk.GetKey{
						{
							Key: "foo",
						},
					},
				},
			},
			Results: []*sdk.GetResult{
				&sdk.GetResult{
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
			importMock := new(sdk.MockImport)
			importMock.On("Configure", sdk.Null).Return(nil)
			importMock.On("Get",
				mock.MatchedBy(func(reqs []*sdk.GetReq) bool {
					if len(tc.Requests) != len(reqs) {
						return false
					}

					for i := range tc.Requests {
						tc.Requests[i].ExecDeadline = reqs[i].ExecDeadline
					}

					return reflect.DeepEqual(tc.Requests, reqs)
				})).Return(tc.Results, nil)

			obj, closer := testImportServeGRPC(t, importMock)
			defer closer()

			// We need to configure first
			if err := obj.Configure(nil); err != nil {
				t.Fatalf("err: %s", err)
			}

			// Get
			results, err := obj.Get(tc.Requests)
			importMock.AssertExpectations(t)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if !reflect.DeepEqual(results, tc.Results) {
				t.Fatalf("expected %#v, got %#v", tc.Results, results)
			}
		})
	}
}
