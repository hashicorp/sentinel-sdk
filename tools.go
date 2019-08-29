// +build tools

package sdk

import (
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/vektra/mockery/cmd/mockery"
	_ "gotest.tools/gotestsum"
)
