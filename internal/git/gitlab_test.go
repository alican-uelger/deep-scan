//go:build unit

package git

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/api/client-go"
)

func TestNewGitlab(t *testing.T) {
	token := "test-token"
	host := "http://example.com"
	g, err := NewGitlab(token, host)
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
			clientMock := NewGitlabAPIMock(t)
			clientMock.On("SearchProjects", tt.projctName, &gitlab.SearchOptions{
				ListOptions: gitlab.ListOptions{
					Page:    1,
					PerPage: 1,
				},
			}).Return(tt.mockResponse, nil, tt.mockError)

			g := Gitlab{
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
				ID: 1,
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
				ID: 1,
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
			clientMock := NewGitlabAPIMock(t)
			clientMock.On("GetRawFile", tt.project.ID, tt.path, &gitlab.GetRawFileOptions{}).Return(tt.mockResponse, nil, tt.mockError)

			g := Gitlab{
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
		projectID      any
		mockResponse   []*gitlab.TreeNode
		mockError      error
		expectedError  error
		expectedResult []TreeNode
	}{
		{
			name:      "Valid tree",
			projectID: 1,
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
			name:           "Error from API",
			projectID:      1,
			mockResponse:   nil,
			mockError:      errors.New("API error"),
			expectedError:  errors.New("API error"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := NewGitlabAPIMock(t)
			clientMock.On("ListRepositoryTree", tt.projectID, &gitlab.ListTreeOptions{
				Recursive: gitlab.Bool(true),
				ListOptions: gitlab.ListOptions{
					PerPage: 100,
					Page:    1,
				},
			}).Return(tt.mockResponse, &gitlab.Response{NextPage: 0}, tt.mockError)

			g := Gitlab{
				client: clientMock,
			}

			result, err := g.ListRepositoryTree(tt.projectID)
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
			clientMock := NewGitlabAPIMock(t)
			clientMock.On("ListGroupProjects", tt.group, &gitlab.ListGroupProjectsOptions{
				IncludeSubGroups: gitlab.Bool(true),
				ListOptions: gitlab.ListOptions{
					PerPage: 100,
					Page:    1,
				},
			}).Return(tt.mockResponse, &gitlab.Response{NextPage: 0}, tt.mockError)

			g := Gitlab{
				client: clientMock,
			}

			result, err := g.ListGroupProjects(tt.group)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
