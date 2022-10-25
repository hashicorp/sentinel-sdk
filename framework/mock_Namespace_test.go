// Code generated by mockery v2.14.0. DO NOT EDIT.

// Generated code. DO NOT MODIFY.

package framework

import mock "github.com/stretchr/testify/mock"

// MockNamespace is an autogenerated mock type for the Namespace type
type MockNamespace struct {
	mock.Mock
}

// Get provides a mock function with given fields: _a0
func (_m *MockNamespace) Get(_a0 string) (interface{}, error) {
	ret := _m.Called(_a0)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string) interface{}); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockNamespace interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockNamespace creates a new instance of MockNamespace. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockNamespace(t mockConstructorTestingTNewMockNamespace) *MockNamespace {
	mock := &MockNamespace{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
