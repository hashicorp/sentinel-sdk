// This is a template used to generate the main file for the test binaries
// built by the "testing" package for Sentinel plugins. This isn't expected
// to be modified manually.

package main

import (
	"github.com/hashicorp/sentinel-sdk/rpc"

	impl "PATH"
)

func main() {
	rpc.Serve(&rpc.ServeOpts{
		PluginFunc: impl.New,
	})
}
