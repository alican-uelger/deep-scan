// Code generated by mockery v2.53.4. DO NOT EDIT.

package scanner

import mock "github.com/stretchr/testify/mock"

// SopsMock is an autogenerated mock type for the Sops type
type SopsMock struct {
	mock.Mock
}

type SopsMock_Expecter struct {
	mock *mock.Mock
}

func (_m *SopsMock) EXPECT() *SopsMock_Expecter {
	return &SopsMock_Expecter{mock: &_m.Mock}
}

// DecryptFile provides a mock function with given fields: path
func (_m *SopsMock) DecryptFile(path string) (string, error) {
	ret := _m.Called(path)

	if len(ret) == 0 {
		panic("no return value specified for DecryptFile")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(path)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SopsMock_DecryptFile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DecryptFile'
type SopsMock_DecryptFile_Call struct {
	*mock.Call
}

// DecryptFile is a helper method to define mock.On call
//   - path string
func (_e *SopsMock_Expecter) DecryptFile(path interface{}) *SopsMock_DecryptFile_Call {
	return &SopsMock_DecryptFile_Call{Call: _e.mock.On("DecryptFile", path)}
}

func (_c *SopsMock_DecryptFile_Call) Run(run func(path string)) *SopsMock_DecryptFile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *SopsMock_DecryptFile_Call) Return(_a0 string, _a1 error) *SopsMock_DecryptFile_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SopsMock_DecryptFile_Call) RunAndReturn(run func(string) (string, error)) *SopsMock_DecryptFile_Call {
	_c.Call.Return(run)
	return _c
}

// NewSopsMock creates a new instance of SopsMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSopsMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *SopsMock {
	mock := &SopsMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
