// Do not edit mock_Import_Closer.go directly as your changes will be
// overwritten. Instead, edit mock_Import_Closer.go.src and re-run
// "go generate ./" in the root SDK package.
package sdk

import "io"

// MockImportCloser augments MockImport to also implement io.Closer
type MockImportCloser struct {
	MockImport
}

// Close mocks Close for MockImportCloser
func (_m *MockImportCloser) Close() error {
	ret := _m.Called()
	return ret.Error(0)
}

var _ io.Closer = (*MockImportCloser)(nil)
