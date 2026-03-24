// Copyright IBM Corp. 2017, 2025
// SPDX-License-Identifier: MPL-2.0

// Do not edit mock_Plugin_Closer.go directly as your changes will be
// overwritten. Instead, edit mock_Plugin_Closer.go.src and re-run
// "go generate ./" in the root SDK package.
package sdk

import "io"

// MockPluginCloser augments MockPlugin to also implement io.Closer
type MockPluginCloser struct {
	MockPlugin
}

// Close mocks Close for MockPluginCloser
func (_m *MockPluginCloser) Close() error {
	ret := _m.Called()
	return ret.Error(0)
}

var _ io.Closer = (*MockPluginCloser)(nil)
