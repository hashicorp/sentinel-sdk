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
			args := reqs[0].Args
			return len(reqs) == 1 &&
				reflect.DeepEqual(reqs[0].Keys, []string{"key"}) &&
				len(args) == 2 &&
				args[0] == "foo" &&
				args[1] == 42
		})).
		Return([]*sdk.GetResult{
			&sdk.GetResult{
				KeyId: 42,
				Keys:  []string{"key"},
				Value: "bar",
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
			Keys:  []string{"key"},
			Args:  []interface{}{"foo", 42},
		},
	})
	importMock.AssertExpectations(t)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	{
		result := sdk.GetResultList(results).KeyId(42).Value
		if result != "bar" {
			t.Fatalf("bad: %#v", result)
		}
	}
}
