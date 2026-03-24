// Copyright IBM Corp. 2017, 2025
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"bytes"
	"encoding/json"
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
		{
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

func Test_Null_MarshalJSON(t *testing.T) {
	res, err := json.Marshal(Null)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(res, []byte("null")) {
		t.Fatalf("unexpected response, marshal of Null should be \"null\", got %q", string(res))
	}
}
