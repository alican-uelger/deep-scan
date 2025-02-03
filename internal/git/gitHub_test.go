//go:build unit

package git

import (
	"errors"
	"testing"

	"github.com/google/go-github/v50/github"
	"github.com/stretchr/testify/assert"
)

func TestGitHubGetProjectByName(t *testing.T) {
	tests := []struct {
		name            string
		projectName     string
		mockResponse    *github.RepositoriesSearchResult
		mockError       error
		expectedError   error
		expectedProject Project
	}{
		{
			name:        "Valid project",
			projectName: "test",
			mockResponse: &github.RepositoriesSearchResult{
				Repositories: []*github.Repository{
					{ID: github.Int64(1), FullName: github.String("test/project")},
				},
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
			projectName:     "not-found-test",
			mockResponse:    &github.RepositoriesSearchResult{Repositories: []*github.Repository{}},
			mockError:       nil,
			expectedError:   errors.New("project not found: not-found-test"),
			expectedProject: Project{},
		},
		{
			name:            "Error from API",
			projectName:     "test",
			mockResponse:    nil,
			mockError:       errors.New("API error"),
			expectedError:   errors.New("API error"),
			expectedProject: Project{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := NewGitHubAPIMock(t)
			clientMock.On("SearchProjects", tt.projectName, &github.SearchOptions{
				ListOptions: github.ListOptions{
					Page:    1,
					PerPage: 1,
				},
			}).Return(tt.mockResponse, nil, tt.mockError)

			g := GitHub{
				client: clientMock,
			}

			project, err := g.GetProjectByName(tt.projectName)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedProject, project)
		})
	}
}

func TestGitHubGetRawFile(t *testing.T) {
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
				Name:              "project",
				PathWithNamespace: "test/project",
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
				Name:              "project",
				PathWithNamespace: "test/project",
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
			clientMock := NewGitHubAPIMock(t)
			clientMock.On("GetRawFile", "test", "project", tt.path, &github.RepositoryContentGetOptions{}).Return(tt.mockResponse, nil, tt.mockError)

			g := GitHub{
				client: clientMock,
			}

			result, err := g.GetRawFile(tt.project, tt.path)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestGitHubListRepositoryTree(t *testing.T) {
	tests := []struct {
		name           string
		mockProject    Project
		mockResponse   []*github.TreeEntry
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
			mockResponse: []*github.TreeEntry{
				{Path: github.String("file1"), Type: github.String("blob")},
				{Path: github.String("dir1"), Type: github.String("tree")},
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
			clientMock := NewGitHubAPIMock(t)
			clientMock.On("ListRepositoryTree", tt.mockProject.Owner(), tt.mockProject.Name).Return(tt.mockResponse, nil, tt.mockError)

			g := GitHub{
				client: clientMock,
			}

			result, err := g.ListRepositoryTree(tt.mockProject)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestGitHubListGroupProjects(t *testing.T) {
	tests := []struct {
		name           string
		group          string
		mockResponse   []*github.Repository
		mockError      error
		expectedError  error
		expectedResult []Project
	}{
		{
			name:  "Valid group projects",
			group: "test-group",
			mockResponse: []*github.Repository{
				{ID: github.Int64(1), FullName: github.String("test-group/project1")},
				{ID: github.Int64(2), FullName: github.String("test-group/project2")},
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
			clientMock := NewGitHubAPIMock(t)
			clientMock.On("ListGroupProjects", tt.group, &github.RepositoryListByOrgOptions{
				ListOptions: github.ListOptions{
					PerPage: 100,
					Page:    1,
				},
			}).Return(tt.mockResponse, nil, tt.mockError)

			g := GitHub{
				client: clientMock,
			}

			result, err := g.ListGroupProjects(tt.group)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
