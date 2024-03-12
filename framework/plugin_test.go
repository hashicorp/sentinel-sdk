// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package framework

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kr/pretty"

	sdk "github.com/hashicorp/sentinel-sdk"
)

func TestPlugin_impl(t *testing.T) {
	var _ sdk.Plugin = new(Plugin)
}

//-------------------------------------------------------------------
// Configure

func TestPluginConfigure(t *testing.T) {
	mockRoot := new(MockNamespaceCreator)
	mockRoot.On("Configure",
		map[string]interface{}{"key": 42}).Return(nil)

	impt := &Plugin{Root: mockRoot}
	err := impt.Configure(map[string]interface{}{"key": 42})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	mockRoot.AssertExpectations(t)
}

func TestPluginConfigure_noNamespace(t *testing.T) {
	cases := []struct {
		Name string
		Root Root
		Err  bool
	}{
		{
			"root with no other implementations",
			&rootNoImpl{},
			true,
		},

		{
			"root with Namespace",
			&rootNamespace{},
			false,
		},

		{
			"root with NamespaceCreator",
			&rootNamespaceCreator{},
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			impt := &Plugin{Root: tc.Root}
			err := impt.Configure(map[string]interface{}{})
			if (err != nil) != tc.Err {
				t.Fatalf("err: %s", err)
			}
		})
	}
}

type rootNoImpl struct{}

func (r *rootNoImpl) Configure(map[string]interface{}) error { return nil }

type rootNamespace struct{}

func (r *rootNamespace) Configure(map[string]interface{}) error { return nil }
func (r *rootNamespace) Get(string) (interface{}, error)        { return nil, nil }

type rootNamespaceCreator struct{}

func (r *rootNamespaceCreator) Configure(map[string]interface{}) error { return nil }
func (r *rootNamespaceCreator) Namespace() Namespace                   { return nil }

//-------------------------------------------------------------------
// Get

var getCases = []struct {
	Name        string
	Root        Root
	Req         []*sdk.GetReq
	Resp        []*sdk.GetResult
	ExpectedErr string
}{
	{
		"key get",
		&rootEmbedNamespace{&nsKeyValue{Key: "foo", Value: "bar"}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo"},
				KeyId: 42,
				Value: "bar",
			},
		},
		"",
	},

	{
		"key get nil",
		&rootEmbedNamespace{&nsKeyValue{Key: "foo", Value: nil}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo"},
				KeyId: 42,
				Value: sdk.Undefined,
			},
		},
		"",
	},

	{
		"key get map",
		&rootEmbedNamespace{&nsKeyValue{
			Key: "foo",
			Value: map[string]interface{}{
				"bar": 42,
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo"},
				KeyId: 42,
				Value: map[string]interface{}{
					"bar": 42,
				},
			},
		},
		"",
	},

	{
		"key get map with int key",
		&rootEmbedNamespace{&nsKeyValue{
			Key: "foo",
			Value: map[int]interface{}{
				42: "bar",
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo"},
				KeyId: 42,
				Value: map[int]interface{}{
					42: "bar",
				},
			},
		},
		"",
	},

	{
		"key get map with nil value",
		&rootEmbedNamespace{&nsKeyValue{
			Key: "foo",
			Value: map[string]interface{}{
				"bar": nil,
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo"},
				KeyId: 42,
				Value: map[string]interface{}{
					"bar": nil,
				},
			},
		},
		"",
	},

	{
		"key get slice with nil value",
		&rootEmbedNamespace{&nsKeyValue{
			Key:   "foo",
			Value: []interface{}{nil},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo"},
				KeyId: 42,
				Value: []interface{}{nil},
			},
		},
		"",
	},

	{
		"key get map with nil value, full key in get request",
		&rootEmbedNamespace{&nsKeyValue{
			Key: "foo",
			Value: map[string]interface{}{
				"bar": nil,
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
					{Key: "bar"},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo", "bar"},
				KeyId: 42,
				Value: sdk.Null,
			},
		},
		"",
	},

	{
		"key get map with unknown key, full key in get request",
		&rootEmbedNamespace{&nsKeyValue{
			Key:   "foo",
			Value: map[string]interface{}{},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
					{Key: "bar"},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo", "bar"},
				KeyId: 42,
				Value: sdk.Undefined,
			},
		},
		"",
	},

	{
		"key get invalid",
		&rootEmbedNamespace{&nsKeyValue{Key: "foo", Value: "bar"}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "unknown"},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"unknown"},
				KeyId: 42,
				Value: sdk.Undefined,
			},
		},
		"",
	},

	{
		"key get nested",
		&rootEmbedNamespace{&nsKeyValue{
			Key: "foo",
			Value: &nsKeyValue{
				Key:   "child",
				Value: "bar",
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
					{Key: "child"},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo", "child"},
				KeyId: 42,
				Value: "bar",
			},
		},
		"",
	},

	{
		"key get nested list",
		&rootEmbedNamespace{&nsKeyValue{
			Key: "foo",
			Value: &nsList{
				Value: []interface{}{"bar", "baz"},
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo"},
				KeyId: 42,
				Value: []interface{}{"bar", "baz"},
			},
		},
		"",
	},

	{
		"key get map value",
		&rootEmbedNamespace{&nsKeyValue{
			Key: "foo",
			Value: map[string]interface{}{
				"child": "bar",
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
					{Key: "child"},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo", "child"},
				KeyId: 42,
				Value: "bar",
			},
		},
		"",
	},

	{
		"key get map value with specific type",
		&rootEmbedNamespace{&nsKeyValue{
			Key: "foo",
			Value: map[string]int64{
				"child": 84,
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
					{Key: "child"},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo", "child"},
				KeyId: 42,
				Value: int64(84),
			},
		},
		"",
	},

	{
		"key get missing map value with specific type",
		&rootEmbedNamespace{&nsKeyValue{
			Key: "foo",
			Value: map[string]int64{
				"child": 84,
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
					{Key: "nope"},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo", "nope"},
				KeyId: 42,
				Value: sdk.Undefined,
			},
		},
		"",
	},

	{
		"key get map value that is a namespace",
		&rootEmbedNamespace{&nsKeyValue{
			Key: "foo",
			Value: map[string]interface{}{
				"child": &nsKeyValueMap{
					Value: map[string]interface{}{
						"foo": "bar",
					},
				},
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo"},
				KeyId: 42,
				Value: map[string]interface{}{
					"child": map[string]interface{}{
						"foo": "bar",
					},
				},
			},
		},
		"",
	},

	{
		"key get map value that is a namespace (two levels)",
		&rootEmbedNamespace{&nsKeyValue{
			Key: "foo",
			Value: map[string]interface{}{
				"child": &nsKeyValueMap{
					Value: map[string]interface{}{
						"foo": &nsKeyValueMap{
							Value: map[string]interface{}{
								"bar": "baz",
							},
						},
					},
				},
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo"},
				KeyId: 42,
				Value: map[string]interface{}{
					"child": map[string]interface{}{
						"foo": map[string]interface{}{
							"bar": "baz",
						},
					},
				},
			},
		},
		"",
	},

	{
		"key get slice value that is a namespace",
		&rootEmbedNamespace{&nsKeyValue{
			Key: "foo",
			Value: []interface{}{
				&nsKeyValueMap{
					Value: map[string]interface{}{
						"foo": "bar",
					},
				},
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo"},
				KeyId: 42,
				Value: []interface{}{
					map[string]interface{}{
						"foo": "bar",
					},
				},
			},
		},
		"",
	},

	{
		"key get nested invalid",
		&rootEmbedNamespace{&nsKeyValue{
			Key: "foo",
			Value: &nsKeyValue{
				Key:   "child",
				Value: "bar",
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
					{Key: "unknown"},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo", "unknown"},
				KeyId: 42,
				Value: sdk.Undefined,
			},
		},
		"",
	},

	{
		"key get multiple",
		&rootEmbedNamespace{&nsKeyValueMap{
			Value: map[string]interface{}{
				"foo": "foovalue",
				"bar": "barvalue",
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId: 1,
			},
			{
				Keys: []sdk.GetKey{
					{Key: "bar"},
				},
				KeyId: 3,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo"},
				KeyId: 1,
				Value: "foovalue",
			},
			{
				Keys:  []string{"bar"},
				KeyId: 3,
				Value: "barvalue",
			},
		},
		"",
	},

	{
		"key get map",
		&rootEmbedNamespace{&nsKeyValue{
			Key: "foo",
			Value: &nsKeyValueMap{
				Value: map[string]interface{}{
					"key":     "value",
					"another": "value",
				},
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo"},
				KeyId: 42,
				Value: map[string]interface{}{
					"key":     "value",
					"another": "value",
				},
			},
		},
		"",
	},

	{
		"key get result is a namespace that does not implement map",
		&rootEmbedNamespace{
			&nsKeyValue{
				Key: "foo",
				Value: &nsKeyValue{
					Key:   "one",
					Value: "two",
				},
			},
		},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo"},
				KeyId: 42,
				Value: &nsKeyValue{
					Key:   "one",
					Value: "two",
				},
			},
		},
		"",
	},

	{
		"key call",
		&rootEmbedCall{&nsCall{
			F: func(v string) (interface{}, error) {
				return v, nil
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo", Args: []interface{}{"asdf"}},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo"},
				KeyId: 42,
				Value: "asdf",
			},
		},
		"",
	},

	{
		"key call with invalid but convertable type",
		&rootEmbedCall{&nsCall{
			F: func(v string) (interface{}, error) {
				return v, nil
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo", Args: []interface{}{42}},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo"},
				KeyId: 42,
				Value: "42",
			},
		},
		"",
	},

	{
		"key call with namespace return",
		&rootEmbedCall{&nsCall{
			F: func(v string) (interface{}, error) {
				return &nsKeyValueMap{Value: map[string]interface{}{v: "bar"}}, nil
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo", Args: []interface{}{"asdf"}},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo"},
				KeyId: 42,
				Value: map[string]interface{}{
					"asdf": "bar",
				},
			},
		},
		"",
	},

	{
		"key call with no error result",
		&rootEmbedCall{&nsCall{
			F: func(v string) interface{} {
				return v
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo", Args: []interface{}{"asdf"}},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo"},
				KeyId: 42,
				Value: "asdf",
			},
		},
		"",
	},

	{
		"multiple levels, multiple calls",
		&rootEmbedCall{&nsCall{
			F: func(v string) (interface{}, error) {
				if v != "one" {
					return nil, fmt.Errorf("expected \"one\", got %q", v)
				}

				return &nsCall{
					F: func(a, b int) (interface{}, error) {
						if a != 2 && b != 3 {
							return nil, fmt.Errorf("expected: 2, 3; got: %d, %d", a, b)
						}

						return "baz", nil
					},
				}, nil
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo", Args: []interface{}{"one"}},
					{Key: "bar", Args: []interface{}{2, 3}},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo", "bar"},
				KeyId: 42,
				Value: "baz",
			},
		},
		"",
	},

	{
		"call, get, call",
		&rootEmbedCall{&nsCall{
			F: func(v string) (interface{}, error) {
				if v != "one" {
					return nil, fmt.Errorf("expected \"one\", got %q", v)
				}

				return &nsKeyValue{
					Key: "bar",
					Value: &nsCall{
						F: func(a, b int) (interface{}, error) {
							if a != 2 && b != 3 {
								return nil, fmt.Errorf("expected: 2, 3; got: %d, %d", a, b)
							}

							return "qux", nil
						},
					},
				}, nil
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo", Args: []interface{}{"one"}},
					{Key: "bar"},
					{Key: "baz", Args: []interface{}{2, 3}},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo", "bar", "baz"},
				KeyId: 42,
				Value: "qux",
			},
		},
		"",
	},

	{
		"get call with receiver",
		&rootNew{},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId:   42,
				Context: map[string]interface{}{"a": "b"},
			},
		},
		[]*sdk.GetResult{
			{
				Keys:     []string{"foo"},
				KeyId:    42,
				Value:    map[string]interface{}{"result": "New called"},
				Context:  map[string]interface{}{"foo": map[string]interface{}{"result": "New called"}},
				Callable: true,
			},
		},
		"",
	},

	{
		"get call with receiver (assert input)",
		&rootNew{
			F: func(data map[string]interface{}) (Namespace, error) {
				if data["a"] == "b" {
					return &nsKeyValueMap{map[string]interface{}{
						"foo": map[string]interface{}{
							"result": "OK",
						},
					}}, nil
				}

				return nil, nil
			},
		},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId:   42,
				Context: map[string]interface{}{"a": "b"},
			},
		},
		[]*sdk.GetResult{
			{
				Keys:     []string{"foo"},
				KeyId:    42,
				Value:    map[string]interface{}{"result": "OK"},
				Context:  map[string]interface{}{"foo": map[string]interface{}{"result": "OK"}},
				Callable: true,
			},
		},
		"",
	},

	{
		"func call with receiver (non-callable result, mutate receiver)",
		&rootNew{
			F: func(data map[string]interface{}) (Namespace, error) {
				return &nsMutable{Value: data["value"].(string)}, nil
			},
		},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo", Args: []interface{}{"two"}},
				},
				KeyId:   42,
				Context: map[string]interface{}{"value": "one"},
			},
		},
		[]*sdk.GetResult{
			{
				Keys:    []string{"foo"},
				KeyId:   42,
				Value:   "OK",
				Context: map[string]interface{}{"value": "two"},
			},
		},
		"",
	},

	{
		"get without receiver on New implementation",
		&rootNew{},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:     []string{"foo"},
				KeyId:    42,
				Value:    map[string]interface{}{"result": "New not called (Get)"},
				Callable: true,
			},
		},
		"",
	},

	{
		"func call without receiver on New implementation",
		&rootNew{},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo", Args: []interface{}{}},
				},
				KeyId: 42,
			},
		},
		[]*sdk.GetResult{
			{
				Keys:     []string{"foo"},
				KeyId:    42,
				Value:    map[string]interface{}{"result": "New not called (Func)"},
				Callable: true,
			},
		},
		"",
	},

	{
		"unknown receiver data from instantiation",
		&rootNew{
			F: func(map[string]interface{}) (Namespace, error) {
				return nil, nil
			},
		},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId:   42,
				Context: map[string]interface{}{"a": "b"},
			},
		},
		[]*sdk.GetResult{
			{
				Keys:  []string{"foo"},
				KeyId: 42,
				Value: sdk.Undefined,
			},
		},
		"",
	},

	{
		"key call unsupported",
		&rootEmbedNamespace{&nsKeyValue{Key: "foo", Value: "bar"}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo", Args: []interface{}{"asdf"}},
				},
				KeyId: 42,
			},
		},
		nil,
		`key "foo" doesn't support function calls`,
	},

	{
		"key call with too few arguments",
		&rootEmbedCall{&nsCall{
			F: func(v string) (interface{}, error) {
				return v, nil
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo", Args: []interface{}{}},
				},
				KeyId: 42,
			},
		},
		nil,
		`error calling function "foo": expected 1 arguments, got 0`,
	},

	{
		"key call with too many arguments",
		&rootEmbedCall{&nsCall{
			F: func(v string) (interface{}, error) {
				return v, nil
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo", Args: []interface{}{1, 2}},
				},
				KeyId: 42,
			},
		},
		nil,
		`error calling function "foo": expected 1 arguments, got 2`,
	},

	{
		"multi-level key call error message",
		&rootEmbedNamespace{&nsKeyValue{
			Key: "foo",
			Value: &nsCall{
				F: func() (interface{}, error) {
					return "", fmt.Errorf("foo")
				},
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
					{Key: "bar", Args: []interface{}{}},
				},
			},
		},
		nil,
		`error calling function "bar": foo`,
	},

	{
		"bad get",
		&rootEmbedNamespace{&nsGetErr{}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId: 42,
			},
		},
		nil,
		`error retrieving key "foo": get error`,
	},

	{
		"bad map",
		&rootEmbedNamespace{&nsKeyValue{Key: "foo", Value: &nsMapErr{}}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId: 42,
			},
		},
		nil,
		`error retrieving key "foo": map error`,
	},

	{
		"multi-call, error in outer",
		&rootEmbedCall{&nsCall{
			F: func(v string) (interface{}, error) {
				if v != "one" {
					return nil, fmt.Errorf("expected \"one\", got %q", v)
				}

				return &nsCall{
					F: func(a, b int) (interface{}, error) {
						if a != 2 && b != 3 {
							return nil, fmt.Errorf("expected: 2, 3; got: %d, %d", a, b)
						}

						return "baz", nil
					},
				}, nil
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo", Args: []interface{}{"bad"}},
					{Key: "bar", Args: []interface{}{2, 3}},
				},
				KeyId: 42,
			},
		},
		nil,
		`error calling function "foo": expected "one", got "bad"`,
	},

	{
		"multi-call, error in inner",
		&rootEmbedCall{&nsCall{
			F: func(v string) (interface{}, error) {
				if v != "one" {
					return nil, fmt.Errorf("expected \"one\", got %q", v)
				}

				return &nsCall{
					F: func(a, b int) (interface{}, error) {
						if a != 2 && b != 3 {
							return nil, fmt.Errorf("expected: 2, 3; got: %d, %d", a, b)
						}

						return "baz", nil
					},
				}, nil
			},
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo", Args: []interface{}{"one"}},
					{Key: "bar", Args: []interface{}{42, 43}},
				},
				KeyId: 42,
			},
		},
		nil,
		`error calling function "bar": expected: 2, 3; got: 42, 43`,
	},

	{
		"error from receiver constructor",
		&rootNew{
			F: func(map[string]interface{}) (Namespace, error) {
				return nil, fmt.Errorf("OK")
			},
		},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId:   42,
				Context: map[string]interface{}{"a": "b"},
			},
		},
		nil,
		"error instantiating namespace: OK",
	},

	{
		"error from receiver constructor, function call",
		&rootNew{
			F: func(map[string]interface{}) (Namespace, error) {
				return nil, nil
			},
		},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
					{Key: "bar", Args: []interface{}{"one"}},
				},
				KeyId:   42,
				Context: map[string]interface{}{"a": "b"},
			},
		},
		nil,
		`attempting to call function "foo.bar" on undefined receiver`,
	},

	{
		"Context supplied but New not implemented",
		&rootEmbedNamespace{&nsKeyValue{
			Key:   "foo",
			Value: "bar",
		}},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId:   42,
				Context: map[string]interface{}{"a": "b"},
			},
		},
		nil,
		"sdk.GetReq.Context present but plugin does not support framework.New",
	},

	{
		"receiver marshal error",
		&rootNew{
			F: func(data map[string]interface{}) (Namespace, error) {
				return &nsKeyValueMap{map[string]interface{}{
					"foo": map[string]interface{}{
						"result": "Not OK",
					},
					"bar": &nsMapErr{},
				}}, nil
			},
		},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId:   42,
				Context: map[string]interface{}{"a": "b"},
			},
		},
		nil,
		`error marshaling receiver after retrieving key "foo": map error`,
	},

	{
		"receiver non-object",
		&rootNew{
			F: func(data map[string]interface{}) (Namespace, error) {
				return &nsKeyValue{
					Key:   "foo",
					Value: "Not OK",
				}, nil
			},
		},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId:   42,
				Context: map[string]interface{}{"a": "b"},
			},
		},
		nil,
		`error marshaling receiver after retrieving key "foo": receiver is no longer an object`,
	},

	{
		"receiver nil object",
		&rootNew{
			F: func(data map[string]interface{}) (Namespace, error) {
				return &nsNilable{}, nil
			},
		},
		[]*sdk.GetReq{
			{
				Keys: []sdk.GetKey{
					{Key: "foo"},
				},
				KeyId:   42,
				Context: map[string]interface{}{"a": "b"},
			},
		},
		nil,
		`error marshaling receiver after retrieving key "foo": receiver is now nil`,
	},
}

func TestPluginGet(t *testing.T) {
	for _, tc := range getCases {
		t.Run(tc.Name, func(t *testing.T) {
			impt := &Plugin{
				Root: tc.Root,
			}

			// Configure
			err := impt.Configure(map[string]interface{}{})
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			// Perform the req
			actual, err := impt.Get(tc.Req)
			if err != nil {
				if tc.ExpectedErr != "" {
					if err.Error() != tc.ExpectedErr {
						t.Fatalf("expected error to be %q, got %q", tc.ExpectedErr, err.Error())
					}
				} else {
					t.Fatalf("err: %s", err)
				}
			}

			// Compare the response
			if !reflect.DeepEqual(actual, tc.Resp) {
				t.Fatalf("bad: %s", pretty.Sprint(actual))
			}
		})
	}
}

func TestPluginGetConcurrent(t *testing.T) {
	for _, tc := range getCases {
		t.Run(tc.Name, func(t *testing.T) {
			impt := &Plugin{
				Root: tc.Root,
			}

			// Configure
			err := impt.Configure(map[string]interface{}{})
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			// Perform the reqs, in parallel, 1000x
			var wg sync.WaitGroup
			runErr := make(chan error, 1000)
			actuals := make(chan []*sdk.GetResult, 1000)
			for i := 0; i < 1000; i++ {
				wg.Add(1)
				go func() {
					actual, err := impt.Get(tc.Req)
					if err != nil {
						if tc.ExpectedErr != "" {
							if err.Error() != tc.ExpectedErr {
								runErr <- fmt.Errorf("expected error to be %q, got %q", tc.ExpectedErr, err.Error())
							}
						} else {
							runErr <- fmt.Errorf("err: %s", err)
						}

						return
					}

					actuals <- actual
				}()

				wg.Done()
			}

			wg.Wait()

			nErrs := len(runErr)
			if nErrs > 0 {
				t.Fatalf("%d errors, first error: %s", nErrs, <-runErr)
			}

			for len(actuals) > 0 {
				actual := <-actuals
				if !reflect.DeepEqual(actual, tc.Resp) {
					t.Fatalf("bad response encountered: %s", pretty.Sprint(actual))
				}
			}
		})
	}
}

// TestPluginGetImmutable checks to ensure that any deep response
// reflection we do does not alter the structure of the original
// underlying plugin data.
func TestPluginGetImmutable(t *testing.T) {
	imptF := func() *Plugin {
		return &Plugin{
			Root: &rootEmbedNamespace{&nsKeyValue{
				Key: "foo",
				Value: &nsKeyValueMap{
					Value: map[string]interface{}{
						"key":     "value",
						"another": "value",
						"embedded_map": &nsKeyValueMap{
							Value: map[string]interface{}{
								"key": "value",
							},
						},
						"embedded_slice": []*nsKeyValueMap{
							{
								Value: map[string]interface{}{
									"key": "value",
								},
							},
						},
					},
				},
			}},
		}
	}

	actual := imptF()
	expected := imptF()

	err := actual.Configure(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	_, err = actual.Get([]*sdk.GetReq{
		{
			Keys: []sdk.GetKey{
				{Key: "foo"},
			},
			KeyId: 42,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatal("plugin data should not have been altered")
	}
}

// rootEmbedNamespace embeds a Namespace for easy testing.
type rootEmbedNamespace struct{ Namespace }

func (r *rootEmbedNamespace) Configure(map[string]interface{}) error { return nil }

// rootEmbedCall embeds a Call for easy testing.
type rootEmbedCall struct{ C Call }

func (r *rootEmbedCall) Configure(map[string]interface{}) error { return nil }

func (r *rootEmbedCall) Get(k string) (interface{}, error) {
	return r.C.Get(k)
}

func (r *rootEmbedCall) Func(k string) interface{} {
	return r.C.Func(k)
}

// nsKeyValue implements Namespace and returns a value for a specific key.
type nsKeyValue struct {
	Key   string
	Value interface{}
}

func (v *nsKeyValue) Get(key string) (interface{}, error) {
	if v.Key != key {
		return nil, nil
	}

	return v.Value, nil
}

// nsKeyValueMap implements Namespace and returns a value by looking up
// the key in a static map.
type nsKeyValueMap struct{ Value map[string]interface{} }

func (v *nsKeyValueMap) Get(key string) (interface{}, error) {
	result, ok := v.Value[key]
	if !ok {
		return nil, nil
	}

	return result, nil
}

func (v *nsKeyValueMap) Map() (map[string]interface{}, error) {
	return v.Value, nil
}

// nsList implements Namespace and returns a value by looking up
// the index in a slice
type nsList struct{ Value []interface{} }

func (v *nsList) Get(key string) (interface{}, error) {
	return nil, nil
}

func (v *nsList) List() ([]interface{}, error) {
	return v.Value, nil
}

// nsCall implements Call that you can implement with a function.
type nsCall struct {
	F interface{}
}

func (v *nsCall) Func(key string) interface{} {
	return v.F
}

func (v *nsCall) Get(key string) (interface{}, error) {
	return nil, fmt.Errorf("can't get")
}

// nsGetErr implements Namespace and just stubs an error response.
type nsGetErr struct{}

func (v *nsGetErr) Get(string) (interface{}, error) { return nil, errors.New("get error") }

// nsvMapErr implements a Map Namespace and just stubs an error response.
type nsMapErr struct{}

func (v *nsMapErr) Get(string) (interface{}, error)      { return map[string]interface{}{}, nil }
func (v *nsMapErr) Map() (map[string]interface{}, error) { return nil, errors.New("map error") }

// rootNew implements a mock namespace for testing framework.New.
type rootNew struct {
	F func(map[string]interface{}) (Namespace, error)
}

func (r *rootNew) Configure(map[string]interface{}) error { return nil }

func (r *rootNew) Get(k string) (interface{}, error) {
	return map[string]interface{}{"result": "New not called (Get)"}, nil
}

func (r *rootNew) Func(k string) interface{} {
	return func() (interface{}, error) {
		return map[string]interface{}{"result": "New not called (Func)"}, nil
	}
}

func (r *rootNew) New(data map[string]interface{}) (Namespace, error) {
	if r.F != nil {
		return r.F(data)
	}

	return &nsKeyValueMap{map[string]interface{}{
		"foo": map[string]interface{}{
			"result": "New called",
		},
	}}, nil
}

// nsMutable represents a mutable namespace.
type nsMutable struct {
	Value string
}

func (v *nsMutable) Get(key string) (interface{}, error) {
	return v.Value, nil
}

func (v *nsMutable) Map() (map[string]interface{}, error) {
	return map[string]interface{}{"value": v.Value}, nil
}

func (v *nsMutable) Func(key string) interface{} {
	return func(s string) (interface{}, error) {
		v.Value = s
		return "OK", nil
	}
}

// nsNilable is a fake namespace that returns nil for everything.
type nsNilable struct{}

func (v *nsNilable) Get(key string) (interface{}, error) {
	return nil, nil
}

func (v *nsNilable) Map() (map[string]interface{}, error) {
	return nil, nil
}

// Test Get with a Root that implements NamespaceCreator.
func TestPluginGet_namespaceCreator(t *testing.T) {
	impt := &Plugin{
		Root: &rootCounter{},
	}

	// Configure
	err := impt.Configure(map[string]interface{}{})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	{
		// Make a Get request, assert the response
		actual, err := impt.Get([]*sdk.GetReq{
			{
				ExecId: 1,
				Keys:   []sdk.GetKey{{Key: "foo"}},
				KeyId:  1,
			},
			{
				ExecId: 2,
				Keys:   []sdk.GetKey{{Key: "bar"}},
				KeyId:  3,
			},
			{
				ExecId: 1,
				Keys:   []sdk.GetKey{{Key: "baz"}},
				KeyId:  5,
			},
		})
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		expected := []*sdk.GetResult{
			{
				Keys:  []string{"foo"},
				KeyId: 1,
				Value: uint64(1),
			},
			{
				Keys:  []string{"bar"},
				KeyId: 3,
				Value: uint64(1),
			},
			{
				Keys:  []string{"baz"},
				KeyId: 5,
				Value: uint64(2),
			},
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("expected %#v, got %#v", expected, actual)
		}
	}
}

// Test Get with a Root that implements NamespaceCreator expires the
// created namespaces properly.
func TestPluginGet_namespaceCreatorExpire(t *testing.T) {
	impt := &Plugin{
		Root: &rootCounter{},
	}

	// Configure
	err := impt.Configure(map[string]interface{}{})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Create the deadlines
	deadline1 := time.Now().Add(10 * time.Millisecond)
	deadline2 := time.Now().Add(100 * time.Millisecond)

	{
		// Make a Get request, assert the response
		_, err := impt.Get([]*sdk.GetReq{
			{
				ExecId:       1,
				ExecDeadline: deadline1,
				Keys:         []sdk.GetKey{{Key: "foo"}},
				KeyId:        1,
			},
			{
				ExecId:       2,
				ExecDeadline: deadline2,
				Keys:         []sdk.GetKey{{Key: "bar"}},
				KeyId:        3,
			},
		})
		if err != nil {
			t.Fatalf("err: %s", err)
		}
	}

	// Sleep for one
	time.Sleep(time.Until(deadline1) + 5*time.Millisecond)

	// Verify we have only one
	impt.namespaceLock.RLock()
	if len(impt.namespaceMap) != 1 {
		t.Fatal("should be one")
	}
	impt.namespaceLock.RUnlock()

	// Sleep for two
	time.Sleep(time.Until(deadline2) + 5*time.Millisecond)

	// Verify we have only one
	impt.namespaceLock.RLock()
	if len(impt.namespaceMap) != 0 {
		t.Fatal("should be empty")
	}
	impt.namespaceLock.RUnlock()
}

type rootCounter struct{}

func (r *rootCounter) Configure(map[string]interface{}) error { return nil }
func (r *rootCounter) Namespace() Namespace                   { return &nsCounter{} }

// nsCounter is a stateful Namespace that increments a counter every Get.
type nsCounter struct {
	Count uint64
}

func (v *nsCounter) Get(string) (interface{}, error) {
	return atomic.AddUint64(&v.Count, 1), nil
}
