package encoding

import (
	"reflect"
	"testing"

	"github.com/hashicorp/sentinel-sdk"
)

func TestEncoding(t *testing.T) {
	for _, tc := range encodingTests {
		t.Run(tc.Name, func(t *testing.T) {
			// Go => Value
			value, err := GoToValue(tc.Source)
			if err != nil {
				t.Fatalf("GoToValue err: %s", err)
			}

			// Determine the expected type and value. This is just the
			// value set to tc.Expected unless we wrap it in targetType. We
			// use that wrapper to signal that we want to set the typ to nil
			// to force automatic type inference and test that.
			typ := reflect.ValueOf(tc.Expected).Type()
			expected := tc.Expected
			if n, ok := expected.(targetType); ok {
				expected = n.Expected
				typ = n.Type
			}

			// Value => Go
			actual, err := ValueToGo(value, typ)
			if err != nil {
				if tc.Err {
					return
				}

				t.Fatalf("ValueToGo: %s", err)
			}

			// It should be what we expect
			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("bad: %#v", actual)
			}
		})
	}
}

// targetType is a wrapper struct that can be put around Expected values
// in the table below to signal that the test should run with the "typ"
// parameter set to nil. This will force the Go conversion to enforce the type.
type targetType struct {
	Type     reflect.Type // can be nil to automatically infer
	Expected interface{}
}

// encodingTests are the test cases for all encodings
var encodingTests = []struct {
	Name     string
	Source   interface{}
	Expected interface{}
	Err      bool
}{
	//-----------------------------------------------------------
	// Map

	{
		"map to matching map type",
		map[string]interface{}{
			"foo": 42,
			"bar": 21,
		},
		map[string]int8{
			"foo": 42,
			"bar": 21,
		},
		false,
	},

	{
		"map to interface type",
		map[string]interface{}{
			"foo": 42,
			"bar": 21,
		},
		map[string]interface{}{
			"foo": int64(42),
			"bar": int64(21),
		},
		false,
	},

	{
		"map with nil type",
		map[string]interface{}{
			"foo": 42,
			"bar": 21,
		},
		targetType{Expected: map[string]int64{
			"foo": 42,
			"bar": 21,
		}},
		false,
	},

	{
		"map with interface type",
		map[string]interface{}{
			"foo": 42,
			"bar": 21,
		},
		targetType{
			Type: interfaceTyp,
			Expected: map[string]int64{
				"foo": 42,
				"bar": 21,
			},
		},
		false,
	},

	//-----------------------------------------------------------
	// Slice

	{
		"slice to matching slice type",
		[]int32{1, 2, 3, 4},
		[]int{1, 2, 3, 4},
		false,
	},

	{
		"slice to interface{} slice type",
		[]interface{}{1, "foo"},
		[]interface{}{int64(1), "foo"},
		false,
	},

	{
		"slice to interface{} type",
		[]interface{}{1, "foo"},
		[]interface{}{int64(1), "foo"},
		false,
	},

	{
		"slice to nil type",
		[]int32{1, 2, 3, 4},
		targetType{Expected: []int64{1, 2, 3, 4}},
		false,
	},

	{
		"slice with interface type",
		[]int32{1, 2, 3, 4},
		targetType{
			Type:     interfaceTyp,
			Expected: []int64{1, 2, 3, 4},
		},
		false,
	},

	{
		"slice to incompatible slice type",
		[]int32{1, 2, 3, 4},
		[]bool{},
		true,
	},

	//-----------------------------------------------------------
	// Bool

	{
		"bool to int",
		true,
		int(0),
		true,
	},

	{
		"bool to uint",
		true,
		uint(0),
		true,
	},

	{
		"bool to string",
		true,
		`42`,
		true,
	},

	{
		"bool to bool",
		true,
		true,
		false,
	},

	//-----------------------------------------------------------
	// Int

	{
		"int to int",
		42,
		int(42),
		false,
	},

	{
		"int to int8",
		42,
		int8(42),
		false,
	},

	{
		"int to int16",
		42,
		int16(42),
		false,
	},

	{
		"int to int32",
		42,
		int32(42),
		false,
	},

	{
		"int to int64",
		42,
		int64(42),
		false,
	},

	{
		"int to uint",
		42,
		uint(42),
		false,
	},

	{
		"int to uint8",
		42,
		uint8(42),
		false,
	},

	{
		"int to uint16",
		42,
		uint16(42),
		false,
	},

	{
		"int to uint32",
		42,
		uint32(42),
		false,
	},

	{
		"int to uint64",
		42,
		uint64(42),
		false,
	},

	{
		"int to float32",
		42,
		float32(42),
		false,
	},

	{
		"int to float64",
		42,
		float64(42),
		false,
	},

	{
		"int to uint negative",
		-42,
		uint(0),
		true,
	},

	{
		"int to string",
		42,
		`42`,
		false,
	},

	//-----------------------------------------------------------
	// String

	{
		"string to int",
		`42`,
		int(42),
		false,
	},

	{
		"string to int8",
		`42`,
		int8(42),
		false,
	},

	{
		"string to int16",
		`42`,
		int16(42),
		false,
	},

	{
		"string to int32",
		`42`,
		int32(42),
		false,
	},

	{
		"string to int64",
		`42`,
		int64(42),
		false,
	},

	{
		"string to uint",
		`42`,
		uint(42),
		false,
	},

	{
		"string to uint8",
		`42`,
		uint8(42),
		false,
	},

	{
		"string to uint16",
		`42`,
		uint16(42),
		false,
	},

	{
		"string to uint32",
		`42`,
		uint32(42),
		false,
	},

	{
		"string to uint64",
		`42`,
		uint64(42),
		false,
	},

	{
		"string to float32",
		`42`,
		float32(42),
		false,
	},

	{
		"string to float64",
		`42`,
		float64(42),
		false,
	},

	{
		"string to string",
		`42`,
		`42`,
		false,
	},

	//-----------------------------------------------------------
	// Bool

	{
		"null to null",
		sdk.Null,
		sdk.Null,
		false,
	},

	{
		"null to nil type",
		sdk.Null,
		targetType{Expected: sdk.Null},
		false,
	},

	//-----------------------------------------------------------
	// Undefined

	{
		"undefined to undefined",
		sdk.Undefined,
		sdk.Undefined,
		false,
	},

	{
		"undefined to nil type",
		sdk.Undefined,
		targetType{Expected: sdk.Undefined},
		false,
	},
}
