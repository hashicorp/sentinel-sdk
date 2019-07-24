package framework

import (
	"errors"
	"fmt"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kr/pretty"

	"github.com/hashicorp/sentinel-sdk"
)

func TestImport_impl(t *testing.T) {
	var _ sdk.Import = new(Import)
}

//-------------------------------------------------------------------
// Configure

func TestImportConfigure(t *testing.T) {
	mockRoot := new(MockNamespaceCreator)
	mockRoot.On("Configure",
		map[string]interface{}{"key": 42}).Return(nil)

	impt := &Import{Root: mockRoot}
	err := impt.Configure(map[string]interface{}{"key": 42})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	mockRoot.AssertExpectations(t)
}

func TestImportConfigure_noNamespace(t *testing.T) {
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
			impt := &Import{Root: tc.Root}
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

func TestImportGet(t *testing.T) {
	// Used a lot
	undefined := sdk.Undefined

	cases := []struct {
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
					Keys: []*sdk.GetKey{
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
					Keys: []*sdk.GetKey{
						{Key: "foo"},
					},
					KeyId: 42,
				},
			},
			[]*sdk.GetResult{
				{
					Keys:  []string{"foo"},
					KeyId: 42,
					Value: undefined,
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
					Keys: []*sdk.GetKey{
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
			"key get map with nil value",
			&rootEmbedNamespace{&nsKeyValue{
				Key: "foo",
				Value: map[string]interface{}{
					"bar": nil,
				},
			}},
			[]*sdk.GetReq{
				{
					Keys: []*sdk.GetKey{
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
			"key get invalid",
			&rootEmbedNamespace{&nsKeyValue{Key: "foo", Value: "bar"}},
			[]*sdk.GetReq{
				{
					Keys: []*sdk.GetKey{
						{Key: "unknown"},
					},
					KeyId: 42,
				},
			},
			[]*sdk.GetResult{
				{
					Keys:  []string{"unknown"},
					KeyId: 42,
					Value: undefined,
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
					Keys: []*sdk.GetKey{
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
			"key get map value",
			&rootEmbedNamespace{&nsKeyValue{
				Key: "foo",
				Value: map[string]interface{}{
					"child": "bar",
				},
			}},
			[]*sdk.GetReq{
				{
					Keys: []*sdk.GetKey{
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
					Keys: []*sdk.GetKey{
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
					Keys: []*sdk.GetKey{
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
					Value: undefined,
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
					Keys: []*sdk.GetKey{
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
					Keys: []*sdk.GetKey{
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
					Keys: []*sdk.GetKey{
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
					Keys: []*sdk.GetKey{
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
					Value: undefined,
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
					Keys: []*sdk.GetKey{
						{Key: "foo"},
					},
					KeyId: 1,
				},
				{
					Keys: []*sdk.GetKey{
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
					Keys: []*sdk.GetKey{
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
			"key call",
			&rootEmbedCall{&nsCall{
				F: func(v string) (interface{}, error) {
					return v, nil
				},
			}},
			[]*sdk.GetReq{
				{
					Keys: []*sdk.GetKey{
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
					Keys: []*sdk.GetKey{
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
					Keys: []*sdk.GetKey{
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
			"key call unsupported",
			&rootEmbedNamespace{&nsKeyValue{Key: "foo", Value: "bar"}},
			[]*sdk.GetReq{
				{
					Keys: []*sdk.GetKey{
						{Key: "foo", Args: []interface{}{"asdf"}},
					},
					KeyId: 42,
				},
			},
			nil,
			`key "foo" doesn't support function calls`,
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
					Keys: []*sdk.GetKey{
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
					Keys: []*sdk.GetKey{
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
					Keys: []*sdk.GetKey{
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
			"key call with too few arguments",
			&rootEmbedCall{&nsCall{
				F: func(v string) (interface{}, error) {
					return v, nil
				},
			}},
			[]*sdk.GetReq{
				{
					Keys: []*sdk.GetKey{
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
					Keys: []*sdk.GetKey{
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
					Keys: []*sdk.GetKey{
						{Key: "foo"},
						{Key: "bar", Args: []interface{}{}},
					},
				},
			},
			nil,
			`error calling function "foo.bar": foo`,
		},

		{
			"bad get",
			&rootEmbedNamespace{&nsGetErr{}},
			[]*sdk.GetReq{
				{
					Keys: []*sdk.GetKey{
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
					Keys: []*sdk.GetKey{
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
					Keys: []*sdk.GetKey{
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
					Keys: []*sdk.GetKey{
						{Key: "foo", Args: []interface{}{"one"}},
						{Key: "bar", Args: []interface{}{42, 43}},
					},
					KeyId: 42,
				},
			},
			nil,
			`error calling function "foo.bar": expected: 2, 3; got: 42, 43`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			impt := &Import{
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

// Test Get with a Root that implements NamespaceCreator.
func TestImportGet_namespaceCreator(t *testing.T) {
	impt := &Import{
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
				Keys:   []*sdk.GetKey{{Key: "foo"}},
				KeyId:  1,
			},
			{
				ExecId: 2,
				Keys:   []*sdk.GetKey{{Key: "bar"}},
				KeyId:  3,
			},
			{
				ExecId: 1,
				Keys:   []*sdk.GetKey{{Key: "baz"}},
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
func TestImportGet_namespaceCreatorExpire(t *testing.T) {
	impt := &Import{
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
				Keys:         []*sdk.GetKey{{Key: "foo"}},
				KeyId:        1,
			},
			{
				ExecId:       2,
				ExecDeadline: deadline2,
				Keys:         []*sdk.GetKey{{Key: "bar"}},
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
