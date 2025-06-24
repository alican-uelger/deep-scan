//go:build unit

package git

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/api/client-go"
)

func TestNewGitlab(t *testing.T) {
	token := "test-token"
	host := "http://example.com"
	g, err := NewGitLab(token, host)
	assert.NoError(t, err)
	assert.NotNil(t, g)
}

func TestGetProjectByName(t *testing.T) {
	tests := []struct {
		name            string
		projctName      string
		mockResponse    []*gitlab.Project
		mockError       error
		expectedError   error
		expectedProject Project
	}{
		{
			name:       "Valid project",
			projctName: "test",
			mockResponse: []*gitlab.Project{
				{ID: 1, PathWithNamespace: "test/project"},
			},
			mockError:     nil,
			expectedError: nil,
			expectedProject: Project{
				ID:                1,
				PathWithNamespace: "test/project",
			},
		},
		{
			name:            "No projects found",
			projctName:      "not-found-test",
			mockResponse:    []*gitlab.Project{},
			mockError:       nil,
			expectedError:   errors.New("project not found: not-found-test"),
			expectedProject: Project{},
		},
		{
			name:            "Error from API",
			projctName:      "test",
			mockResponse:    nil,
			mockError:       errors.New("API error"),
			expectedError:   errors.New("API error"),
			expectedProject: Project{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := NewGitLabAPIMock(t)
			clientMock.On("SearchProjects", tt.projctName, &gitlab.SearchOptions{
				ListOptions: gitlab.ListOptions{
					Page:    1,
					PerPage: 1,
				},
			}).Return(tt.mockResponse, nil, tt.mockError)

			g := GitLab{
				client: clientMock,
			}

			project, err := g.GetProjectByName(tt.projctName)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedProject, project)
		})
	}
}

func TestGetRawFile(t *testing.T) {
	tests := []struct {
		name           string
		project        Project
		path           string
		mockResponse   []byte
		mockError      error
		expectedError  error
		expectedResult []byte
	}{
		{
			name: "Valid file",
			project: Project{
				ID:                1,
				Name:              "p1",
				PathWithNamespace: "org/p1",
			},
			path:           "README.md",
			mockResponse:   []byte("file content"),
			mockError:      nil,
			expectedError:  nil,
			expectedResult: []byte("file content"),
		},
		{
			name: "Error from API",
			project: Project{
				ID:                1,
				Name:              "p1",
				PathWithNamespace: "org/p1",
			},
			path:           "README.md",
			mockResponse:   nil,
			mockError:      errors.New("API error"),
			expectedError:  errors.New("API error"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := NewGitLabAPIMock(t)
			clientMock.On("GetRawFile", strconv.Itoa(tt.project.ID), tt.path, mock.Anything).Return(tt.mockResponse, nil, tt.mockError)

			g := GitLab{
				client: clientMock,
			}

			result, err := g.GetRawFile(tt.project, tt.path)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestListRepositoryTree(t *testing.T) {
	tests := []struct {
		name           string
		mockProject    Project
		mockResponse   []*gitlab.TreeNode
		mockError      error
		expectedError  error
		expectedResult []TreeNode
	}{
		{
			name: "Valid tree",
			mockProject: Project{
				Name:              "p-1",
				ID:                1,
				PathWithNamespace: "org/p-1",
			},
			mockResponse: []*gitlab.TreeNode{
				{Path: "file1", Type: "blob"},
				{Path: "dir1", Type: "tree"},
			},
			mockError:     nil,
			expectedError: nil,
			expectedResult: []TreeNode{
				{Path: "file1", Type: "blob", IsTree: false},
				{Path: "dir1", Type: "tree", IsTree: true},
			},
		},
		{
			name: "Error from API",
			mockProject: Project{
				Name:              "p-2",
				ID:                2,
				PathWithNamespace: "org/p-2",
			},
			mockResponse:   nil,
			mockError:      errors.New("API error"),
			expectedError:  errors.New("API error"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := NewGitLabAPIMock(t)
			clientMock.On("ListRepositoryTree", tt.mockProject.PathWithNamespace, &gitlab.ListTreeOptions{
				Recursive: gitlab.Bool(true),
				ListOptions: gitlab.ListOptions{
					PerPage: 100,
					Page:    1,
				},
			}).Return(tt.mockResponse, &gitlab.Response{NextPage: 0}, tt.mockError)

			g := GitLab{
				client: clientMock,
			}

			result, err := g.ListRepositoryTree(tt.mockProject)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestListGroupProjects(t *testing.T) {
	tests := []struct {
		name           string
		group          string
		mockResponse   []*gitlab.Project
		mockError      error
		expectedError  error
		expectedResult []Project
	}{
		{
			name:  "Valid group projects",
			group: "test-group",
			mockResponse: []*gitlab.Project{
				{ID: 1, PathWithNamespace: "test-group/project1"},
				{ID: 2, PathWithNamespace: "test-group/project2"},
			},
			mockError:     nil,
			expectedError: nil,
			expectedResult: []Project{
				{ID: 1, PathWithNamespace: "test-group/project1"},
				{ID: 2, PathWithNamespace: "test-group/project2"},
			},
		},
		{
			name:           "Error from API",
			group:          "test-group",
			mockResponse:   nil,
			mockError:      errors.New("API error"),
			expectedError:  errors.New("API error"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := NewGitLabAPIMock(t)
			clientMock.On("ListGroupProjects", tt.group, &gitlab.ListGroupProjectsOptions{
				IncludeSubGroups: gitlab.Bool(true),
				ListOptions: gitlab.ListOptions{
					PerPage: 100,
					Page:    1,
				},
			}).Return(tt.mockResponse, &gitlab.Response{NextPage: 0}, tt.mockError)

			g := GitLab{
				client: clientMock,
			}

			result, err := g.ListGroupProjects(tt.group)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
