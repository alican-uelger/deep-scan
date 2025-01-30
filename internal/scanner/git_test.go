//go:build unit

package scanner

import (
	"errors"
	"testing"

	"github.com/alican-uelger/deep-scan/internal/git"
	"github.com/alican-uelger/deep-scan/internal/sops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGitSearchSuccessfulSearch(t *testing.T) {
	mockClient := NewGitClientMock(t)
	mockClient.
		On("ListGroupProjects", "org").
		Return([]git.Project{{ID: 1, PathWithNamespace: "org/repo"}}, nil)
	mockClient.
		On("ListRepositoryTree", 1).
		Return([]git.TreeNode{{Path: "file.txt", IsTree: false}}, nil)
	mockClient.
		On("GetRawFile", mock.Anything, "file.txt").
		Return([]byte("file content"), nil)

	mockStorage := NewStorageMock(t)

	g := &Git{
		Client:  mockClient,
		Storage: mockStorage,
		Sops:    sops.New(mockStorage),
	}

	result, err := g.Search("org", SearchOptions{})
	expected := []File{
		{Name: "file.txt", Path: "org/repo", Type: FILE},
	}

	assert.Equal(t, expected, result)
	assert.Nil(t, err)
}

func TestGitSearchClientError(t *testing.T) {
	mockClient := NewGitClientMock(t)
	mockClient.
		On("ListGroupProjects", "org").
		Return(nil, errors.New("client error"))

	mockStorage := NewStorageMock(t)

	g := &Git{
		Client:  mockClient,
		Storage: mockStorage,
		Sops:    sops.New(mockStorage),
	}

	result, err := g.Search("org", SearchOptions{})

	mockClient.AssertNumberOfCalls(t, "ListGroupProjects", 1)

	expected := []File{}

	assert.Equal(t, expected, result)
	assert.Equal(t, errors.New("client error"), err)
}

func TestGitSearchEmptyGroupProjects(t *testing.T) {
	mockClient := NewGitClientMock(t)
	mockClient.
		On("ListGroupProjects", "org").
		Return([]git.Project{}, nil)

	mockStorage := NewStorageMock(t)

	g := &Git{
		Client:  mockClient,
		Storage: mockStorage,
		Sops:    sops.New(mockStorage),
	}

	result, err := g.Search("org", SearchOptions{})

	mockClient.AssertNumberOfCalls(t, "ListGroupProjects", 1)
	mockClient.AssertNumberOfCalls(t, "ListRepositoryTree", 0)

	expected := []File{}

	assert.Equal(t, expected, result)
	assert.Nil(t, err)
}

func TestGitSearchEmptyRepositoryTree(t *testing.T) {
	mockClient := NewGitClientMock(t)
	mockClient.
		On("ListGroupProjects", "org").
		Return([]git.Project{{ID: 1, PathWithNamespace: "org/repo"}}, nil)
	mockClient.
		On("ListRepositoryTree", 1).
		Return([]git.TreeNode{}, nil)

	mockStorage := NewStorageMock(t)

	g := &Git{
		Client:  mockClient,
		Storage: mockStorage,
		Sops:    sops.New(mockStorage),
	}

	result, err := g.Search("org", SearchOptions{})
	expected := []File{}

	assert.Equal(t, expected, result)
	assert.Nil(t, err)
}

func TestGitSearchWithOptions(t *testing.T) {
	mockClient := NewGitClientMock(t)
	mockClient.
		On("ListGroupProjects", "org").
		Return([]git.Project{{ID: 1, PathWithNamespace: "org/repo"}}, nil)
	mockClient.
		On("ListRepositoryTree", 1).
		Return([]git.TreeNode{{Path: "file.txt", IsTree: false}}, nil)
	mockClient.
		On("GetRawFile", mock.Anything, "file.txt").
		Return([]byte("file content"), nil)

	mockStorage := NewStorageMock(t)

	g := &Git{
		Client:  mockClient,
		Storage: mockStorage,
		Sops:    sops.New(mockStorage),
	}

	tests := []struct {
		name     string
		options  SearchOptions
		expected []File
	}{
		{
			name: "ExcludeName",
			options: SearchOptions{
				ExcludeName: []string{"file.txt"},
			},
			expected: []File{},
		},
		{
			name: "ExcludeNameContains",
			options: SearchOptions{
				ExcludeNameContains: []string{"file"},
			},
			expected: []File{},
		},
		{
			name: "ExcludeDir",
			options: SearchOptions{
				ExcludeDir: []string{"org/repo"},
			},
			expected: []File{},
		},
		{
			name: "ExcludeDirContains",
			options: SearchOptions{
				ExcludeDirContains: []string{"repo"},
			},
			expected: []File{},
		},
		{
			name: "ExcludeContent",
			options: SearchOptions{
				ExcludeContent: []string{"file content"},
			},
			expected: []File{},
		},
		{
			name: "Name",
			options: SearchOptions{
				Name: []string{"file.txt"},
			},
			expected: []File{
				{Name: "file.txt", Path: "org/repo", Type: FILE},
			},
		},
		{
			name: "NameContains",
			options: SearchOptions{
				NameContains: []string{"file"},
			},
			expected: []File{
				{Name: "file.txt", Path: "org/repo", Type: FILE},
			},
		},
		{
			name: "NameRegex",
			options: SearchOptions{
				NameRegex: []string{`file\.txt`},
			},
			expected: []File{
				{Name: "file.txt", Path: "org/repo", Type: FILE},
			},
		},
		{
			name: "Path",
			options: SearchOptions{
				Path: []string{"org/repo"},
			},
			expected: []File{
				{Name: "file.txt", Path: "org/repo", Type: FILE},
			},
		},
		{
			name: "PathContains",
			options: SearchOptions{
				PathContains: []string{"repo"},
			},
			expected: []File{
				{Name: "file.txt", Path: "org/repo", Type: FILE},
			},
		},
		{
			name: "PathRegex",
			options: SearchOptions{
				PathRegex: []string{`org/repo`},
			},
			expected: []File{
				{Name: "file.txt", Path: "org/repo", Type: FILE},
			},
		},
		{
			name: "Content",
			options: SearchOptions{
				Content: []string{"file content"},
			},
			expected: []File{
				{Name: "file.txt", Path: "org/repo", Type: FILE},
			},
		},
		{
			name: "ContentRegex",
			options: SearchOptions{
				ContentRegex: []string{`file content`},
			},
			expected: []File{
				{Name: "file.txt", Path: "org/repo", Type: FILE},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := g.Search("org", tt.options)
			assert.Equal(t, tt.expected, result)
			assert.Nil(t, err)
		})
	}
}
