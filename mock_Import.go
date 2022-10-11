// Code generated by mockery v2.14.0. DO NOT EDIT.

package sdk

import mock "github.com/stretchr/testify/mock"

// MockImport is an autogenerated mock type for the Import type
type MockImport struct {
	mock.Mock
}

// Configure provides a mock function with given fields: _a0
func (_m *MockImport) Configure(_a0 interface{}) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: reqs
func (_m *MockImport) Get(reqs []*GetReq) ([]*GetResult, error) {
	ret := _m.Called(reqs)

	var r0 []*GetResult
	if rf, ok := ret.Get(0).(func([]*GetReq) []*GetResult); ok {
		r0 = rf(reqs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*GetResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]*GetReq) error); ok {
		r1 = rf(reqs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockImport interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockImport creates a new instance of MockImport. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockImport(t mockConstructorTestingTNewMockImport) *MockImport {
	mock := &MockImport{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
