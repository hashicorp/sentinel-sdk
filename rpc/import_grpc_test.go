package rpc

import (
	"reflect"
	"testing"

	"github.com/hashicorp/sentinel-sdk"
	"github.com/stretchr/testify/mock"
)

func TestImport_gRPC_configure(t *testing.T) {
	// Create a mock object
	importMock := new(sdk.MockImport)
	importMock.On("Configure",
		map[string]interface{}{"key": int64(42)}).
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
	// Create a mock object
	importMock := new(sdk.MockImport)
	importMock.On("Configure", map[string]interface{}{}).Return(nil)
	importMock.On("Get",
		mock.MatchedBy(func(reqs []*sdk.GetReq) bool {
			return len(reqs) == 1 &&
				len(reqs[0].Keys) == 1 &&
				len(reqs[0].Keys[0].Args) == 2 &&
				reqs[0].Keys[0].Key == "foo" &&
				reqs[0].Keys[0].Args[0] == "foo" &&
				reqs[0].Keys[0].Args[1] == int64(42) &&
				reqs[0].Context["_type"] == "SomeNamespace" &&
				reqs[0].Context["data"].(map[string]interface{})["string"] == "foo" &&
				reqs[0].Context["data"].(map[string]interface{})["number"] == int64(0)
		})).
		Return([]*sdk.GetResult{
			&sdk.GetResult{
				KeyId: 42,
				Keys:  []string{"key"},
				Value: "bar",
				Context: map[string]interface{}{
					"_type": "SomeNamespace",
					"data": map[string]interface{}{
						"string": "bar",
						"number": 1,
					},
				},
				Callable: true,
			},
		}, nil)

	obj, closer := testImportServeGRPC(t, importMock)
	defer closer()

	// We need to configure first
	if err := obj.Configure(nil); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Get
	results, err := obj.Get([]*sdk.GetReq{
		&sdk.GetReq{
			KeyId: 42,
			Keys: []sdk.GetKey{
				{
					Key:  "foo",
					Args: []interface{}{"foo", 42},
				},
			},
			Context: map[string]interface{}{
				"_type": "SomeNamespace",
				"data": map[string]interface{}{
					"string": "foo",
					"number": 0,
				},
			},
		},
	})
	importMock.AssertExpectations(t)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	{
		result := sdk.GetResultList(results).KeyId(42)
		if !reflect.DeepEqual(result, &sdk.GetResult{
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
		}) {
			t.Fatalf("bad: %#v", result)
		}
	}
}
