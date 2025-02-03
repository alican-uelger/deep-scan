// Code generated by mockery v2.51.0. DO NOT EDIT.

package git

import (
	mock "github.com/stretchr/testify/mock"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// GitLabAPIMock is an autogenerated mock type for the GitLabAPI type
type GitLabAPIMock struct {
	mock.Mock
}

type GitLabAPIMock_Expecter struct {
	mock *mock.Mock
}

func (_m *GitLabAPIMock) EXPECT() *GitLabAPIMock_Expecter {
	return &GitLabAPIMock_Expecter{mock: &_m.Mock}
}

// GetRawFile provides a mock function with given fields: projectName, path, opts
func (_m *GitLabAPIMock) GetRawFile(projectName string, path string, opts *gitlab.GetRawFileOptions) ([]byte, *gitlab.Response, error) {
	ret := _m.Called(projectName, path, opts)

	if len(ret) == 0 {
		panic("no return value specified for GetRawFile")
	}

	var r0 []byte
	var r1 *gitlab.Response
	var r2 error
	if rf, ok := ret.Get(0).(func(string, string, *gitlab.GetRawFileOptions) ([]byte, *gitlab.Response, error)); ok {
		return rf(projectName, path, opts)
	}
	if rf, ok := ret.Get(0).(func(string, string, *gitlab.GetRawFileOptions) []byte); ok {
		r0 = rf(projectName, path, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, *gitlab.GetRawFileOptions) *gitlab.Response); ok {
		r1 = rf(projectName, path, opts)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*gitlab.Response)
		}
	}

	if rf, ok := ret.Get(2).(func(string, string, *gitlab.GetRawFileOptions) error); ok {
		r2 = rf(projectName, path, opts)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GitLabAPIMock_GetRawFile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRawFile'
type GitLabAPIMock_GetRawFile_Call struct {
	*mock.Call
}

// GetRawFile is a helper method to define mock.On call
//   - projectName string
//   - path string
//   - opts *gitlab.GetRawFileOptions
func (_e *GitLabAPIMock_Expecter) GetRawFile(projectName interface{}, path interface{}, opts interface{}) *GitLabAPIMock_GetRawFile_Call {
	return &GitLabAPIMock_GetRawFile_Call{Call: _e.mock.On("GetRawFile", projectName, path, opts)}
}

func (_c *GitLabAPIMock_GetRawFile_Call) Run(run func(projectName string, path string, opts *gitlab.GetRawFileOptions)) *GitLabAPIMock_GetRawFile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(*gitlab.GetRawFileOptions))
	})
	return _c
}

func (_c *GitLabAPIMock_GetRawFile_Call) Return(_a0 []byte, _a1 *gitlab.Response, _a2 error) *GitLabAPIMock_GetRawFile_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *GitLabAPIMock_GetRawFile_Call) RunAndReturn(run func(string, string, *gitlab.GetRawFileOptions) ([]byte, *gitlab.Response, error)) *GitLabAPIMock_GetRawFile_Call {
	_c.Call.Return(run)
	return _c
}

// ListGroupProjects provides a mock function with given fields: group, opts
func (_m *GitLabAPIMock) ListGroupProjects(group string, opts *gitlab.ListGroupProjectsOptions) ([]*gitlab.Project, *gitlab.Response, error) {
	ret := _m.Called(group, opts)

	if len(ret) == 0 {
		panic("no return value specified for ListGroupProjects")
	}

	var r0 []*gitlab.Project
	var r1 *gitlab.Response
	var r2 error
	if rf, ok := ret.Get(0).(func(string, *gitlab.ListGroupProjectsOptions) ([]*gitlab.Project, *gitlab.Response, error)); ok {
		return rf(group, opts)
	}
	if rf, ok := ret.Get(0).(func(string, *gitlab.ListGroupProjectsOptions) []*gitlab.Project); ok {
		r0 = rf(group, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*gitlab.Project)
		}
	}

	if rf, ok := ret.Get(1).(func(string, *gitlab.ListGroupProjectsOptions) *gitlab.Response); ok {
		r1 = rf(group, opts)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*gitlab.Response)
		}
	}

	if rf, ok := ret.Get(2).(func(string, *gitlab.ListGroupProjectsOptions) error); ok {
		r2 = rf(group, opts)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GitLabAPIMock_ListGroupProjects_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListGroupProjects'
type GitLabAPIMock_ListGroupProjects_Call struct {
	*mock.Call
}

// ListGroupProjects is a helper method to define mock.On call
//   - group string
//   - opts *gitlab.ListGroupProjectsOptions
func (_e *GitLabAPIMock_Expecter) ListGroupProjects(group interface{}, opts interface{}) *GitLabAPIMock_ListGroupProjects_Call {
	return &GitLabAPIMock_ListGroupProjects_Call{Call: _e.mock.On("ListGroupProjects", group, opts)}
}

func (_c *GitLabAPIMock_ListGroupProjects_Call) Run(run func(group string, opts *gitlab.ListGroupProjectsOptions)) *GitLabAPIMock_ListGroupProjects_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*gitlab.ListGroupProjectsOptions))
	})
	return _c
}

func (_c *GitLabAPIMock_ListGroupProjects_Call) Return(_a0 []*gitlab.Project, _a1 *gitlab.Response, _a2 error) *GitLabAPIMock_ListGroupProjects_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *GitLabAPIMock_ListGroupProjects_Call) RunAndReturn(run func(string, *gitlab.ListGroupProjectsOptions) ([]*gitlab.Project, *gitlab.Response, error)) *GitLabAPIMock_ListGroupProjects_Call {
	_c.Call.Return(run)
	return _c
}

// ListRepositoryTree provides a mock function with given fields: projectPath, opts
func (_m *GitLabAPIMock) ListRepositoryTree(projectPath string, opts *gitlab.ListTreeOptions) ([]*gitlab.TreeNode, *gitlab.Response, error) {
	ret := _m.Called(projectPath, opts)

	if len(ret) == 0 {
		panic("no return value specified for ListRepositoryTree")
	}

	var r0 []*gitlab.TreeNode
	var r1 *gitlab.Response
	var r2 error
	if rf, ok := ret.Get(0).(func(string, *gitlab.ListTreeOptions) ([]*gitlab.TreeNode, *gitlab.Response, error)); ok {
		return rf(projectPath, opts)
	}
	if rf, ok := ret.Get(0).(func(string, *gitlab.ListTreeOptions) []*gitlab.TreeNode); ok {
		r0 = rf(projectPath, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*gitlab.TreeNode)
		}
	}

	if rf, ok := ret.Get(1).(func(string, *gitlab.ListTreeOptions) *gitlab.Response); ok {
		r1 = rf(projectPath, opts)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*gitlab.Response)
		}
	}

	if rf, ok := ret.Get(2).(func(string, *gitlab.ListTreeOptions) error); ok {
		r2 = rf(projectPath, opts)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GitLabAPIMock_ListRepositoryTree_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListRepositoryTree'
type GitLabAPIMock_ListRepositoryTree_Call struct {
	*mock.Call
}

// ListRepositoryTree is a helper method to define mock.On call
//   - projectPath string
//   - opts *gitlab.ListTreeOptions
func (_e *GitLabAPIMock_Expecter) ListRepositoryTree(projectPath interface{}, opts interface{}) *GitLabAPIMock_ListRepositoryTree_Call {
	return &GitLabAPIMock_ListRepositoryTree_Call{Call: _e.mock.On("ListRepositoryTree", projectPath, opts)}
}

func (_c *GitLabAPIMock_ListRepositoryTree_Call) Run(run func(projectPath string, opts *gitlab.ListTreeOptions)) *GitLabAPIMock_ListRepositoryTree_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*gitlab.ListTreeOptions))
	})
	return _c
}

func (_c *GitLabAPIMock_ListRepositoryTree_Call) Return(_a0 []*gitlab.TreeNode, _a1 *gitlab.Response, _a2 error) *GitLabAPIMock_ListRepositoryTree_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *GitLabAPIMock_ListRepositoryTree_Call) RunAndReturn(run func(string, *gitlab.ListTreeOptions) ([]*gitlab.TreeNode, *gitlab.Response, error)) *GitLabAPIMock_ListRepositoryTree_Call {
	_c.Call.Return(run)
	return _c
}

// SearchProjects provides a mock function with given fields: projectName, opts
func (_m *GitLabAPIMock) SearchProjects(projectName string, opts *gitlab.SearchOptions) ([]*gitlab.Project, *gitlab.Response, error) {
	ret := _m.Called(projectName, opts)

	if len(ret) == 0 {
		panic("no return value specified for SearchProjects")
	}

	var r0 []*gitlab.Project
	var r1 *gitlab.Response
	var r2 error
	if rf, ok := ret.Get(0).(func(string, *gitlab.SearchOptions) ([]*gitlab.Project, *gitlab.Response, error)); ok {
		return rf(projectName, opts)
	}
	if rf, ok := ret.Get(0).(func(string, *gitlab.SearchOptions) []*gitlab.Project); ok {
		r0 = rf(projectName, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*gitlab.Project)
		}
	}

	if rf, ok := ret.Get(1).(func(string, *gitlab.SearchOptions) *gitlab.Response); ok {
		r1 = rf(projectName, opts)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*gitlab.Response)
		}
	}

	if rf, ok := ret.Get(2).(func(string, *gitlab.SearchOptions) error); ok {
		r2 = rf(projectName, opts)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GitLabAPIMock_SearchProjects_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SearchProjects'
type GitLabAPIMock_SearchProjects_Call struct {
	*mock.Call
}

// SearchProjects is a helper method to define mock.On call
//   - projectName string
//   - opts *gitlab.SearchOptions
func (_e *GitLabAPIMock_Expecter) SearchProjects(projectName interface{}, opts interface{}) *GitLabAPIMock_SearchProjects_Call {
	return &GitLabAPIMock_SearchProjects_Call{Call: _e.mock.On("SearchProjects", projectName, opts)}
}

func (_c *GitLabAPIMock_SearchProjects_Call) Run(run func(projectName string, opts *gitlab.SearchOptions)) *GitLabAPIMock_SearchProjects_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*gitlab.SearchOptions))
	})
	return _c
}

func (_c *GitLabAPIMock_SearchProjects_Call) Return(_a0 []*gitlab.Project, _a1 *gitlab.Response, _a2 error) *GitLabAPIMock_SearchProjects_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *GitLabAPIMock_SearchProjects_Call) RunAndReturn(run func(string, *gitlab.SearchOptions) ([]*gitlab.Project, *gitlab.Response, error)) *GitLabAPIMock_SearchProjects_Call {
	_c.Call.Return(run)
	return _c
}

// NewGitLabAPIMock creates a new instance of GitLabAPIMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGitLabAPIMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *GitLabAPIMock {
	mock := &GitLabAPIMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
