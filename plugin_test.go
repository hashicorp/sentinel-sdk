package sdk

import (
	"reflect"
	"testing"
)

func TestGetResultListKeyId(t *testing.T) {
	cases := []struct {
		Name     string
		KeyId    uint64
		Expected *GetResult
	}{
		{
			Name:     "found",
			KeyId:    42,
			Expected: &GetResult{KeyId: 42},
		},
		{
			Name:     "not found",
			KeyId:    43,
			Expected: nil,
		},
	}
	results := GetResultList([]*GetResult{
		&GetResult{
			KeyId: 42,
		},
	})

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			actual := results.KeyId(tc.KeyId)
			if !reflect.DeepEqual(tc.Expected, actual) {
				t.Fatalf("expected %#v, got %#v", tc.Expected, actual)
			}
		})
	}
}
