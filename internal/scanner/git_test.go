//go:build unit

package scanner

import (
	"errors"
	"testing"

	"github.com/alican-uelger/deep-scan/internal/git"
	"github.com/alican-uelger/deep-scan/internal/sops"
	"github.com/stretchr/testify/assert"
)

func TestGitSearchSuccessfulSearch(t *testing.T) {
	mockProject := git.Project{ID: 1, PathWithNamespace: "org/repo"}
	mockClient := NewGitClientMock(t)
	mockClient.
		On("ListGroupProjects", "org").
		Return([]git.Project{mockProject}, nil)
	mockClient.
		On("ListRepositoryTree", mockProject).
		Return([]git.TreeNode{{Path: "file.txt", IsTree: false}}, nil)

	mockTextMatcher := NewTextMatcherMock(t)
	mockStorage := NewStorageMock(t)

	g := &Git{
		Client: mockClient,
		Base: Base{
			Storage:     mockStorage,
			Sops:        sops.New(mockStorage),
			TextMatcher: mockTextMatcher,
		},
	}

	result, err := g.Search("org", SearchOptions{})
	expected := []FileMatch{
		{File: File{Name: "file.txt", Path: "org/repo", Type: FILE}},
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
	mockTextMatcher := NewTextMatcherMock(t)

	g := &Git{
		Client: mockClient,
		Base: Base{
			Storage:     mockStorage,
			Sops:        sops.New(mockStorage),
			TextMatcher: mockTextMatcher,
		},
	}

	result, err := g.Search("org", SearchOptions{})

	mockClient.AssertNumberOfCalls(t, "ListGroupProjects", 1)

	var expected []FileMatch

	assert.Equal(t, expected, result)
	assert.Equal(t, errors.New("client error"), err)
}

func TestGitSearchEmptyGroupProjects(t *testing.T) {
	mockClient := NewGitClientMock(t)
	mockClient.
		On("ListGroupProjects", "org").
		Return([]git.Project{}, nil)

	mockStorage := NewStorageMock(t)
	mockTextMatcher := NewTextMatcherMock(t)

	g := &Git{
		Client: mockClient,
		Base: Base{
			Storage:     mockStorage,
			Sops:        sops.New(mockStorage),
			TextMatcher: mockTextMatcher,
		},
	}

	result, err := g.Search("org", SearchOptions{})

	mockClient.AssertNumberOfCalls(t, "ListGroupProjects", 1)
	mockClient.AssertNumberOfCalls(t, "ListRepositoryTree", 0)

	var expected []FileMatch

	assert.Equal(t, expected, result)
	assert.Nil(t, err)
}

func TestGitSearchEmptyRepositoryTree(t *testing.T) {
	mockProject := git.Project{ID: 1, PathWithNamespace: "org/repo"}
	mockClient := NewGitClientMock(t)
	mockClient.
		On("ListGroupProjects", "org").
		Return([]git.Project{mockProject}, nil)
	mockClient.
		On("ListRepositoryTree", mockProject).
		Return([]git.TreeNode{}, nil)

	mockStorage := NewStorageMock(t)
	mockTextMatcher := NewTextMatcherMock(t)

	g := &Git{
		Client: mockClient,
		Base: Base{
			Storage:     mockStorage,
			Sops:        sops.New(mockStorage),
			TextMatcher: mockTextMatcher,
		},
	}

	result, err := g.Search("org", SearchOptions{})
	var expected []FileMatch

	assert.Equal(t, expected, result)
	assert.Nil(t, err)
}
