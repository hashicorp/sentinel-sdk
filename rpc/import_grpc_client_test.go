package rpc

import (
	"io"
	"testing"

	"github.com/hashicorp/sentinel-sdk"
)

func TestImportGRPCClient_impl(t *testing.T) {
	var _ sdk.Import = new(ImportGRPCClient)
	var _ io.Closer = new(ImportGRPCClient)
}
