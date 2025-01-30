//go:build unit

package scanner

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOsSearchSuccessfulSearch(t *testing.T) {
	mockStorage := NewStorageMock(t)
	mockStorage.
		On("ReadDir", "dir").
		Return([]string{"dir/file.txt"}, nil)
	mockStorage.
		On("IsDir", "dir/file.txt").
		Return(false, nil)
	mockStorage.
		On("ReadFile", "dir/file.txt").
		Return([]byte("file content"), nil)

	mockSops := NewSopsMock(t)
	mockTestMatcher := NewTextMatcherMock(t)

	o := &Os{
		Base: Base{
			Storage:     mockStorage,
			Sops:        mockSops,
			TextMatcher: mockTestMatcher,
		},
	}

	result, err := o.Search("dir", SearchOptions{})
	expected := []FileMatch{
		{File: File{Name: "file.txt", Path: "dir", Type: FILE}},
	}

	assert.Equal(t, expected, result)
	assert.Nil(t, err)
}

func TestOsSearchStorageError(t *testing.T) {
	mockStorage := NewStorageMock(t)
	mockStorage.
		On("ReadDir", "dir").
		Return(nil, errors.New("storage error"))

	mockSops := NewSopsMock(t)
	mockTestMatcher := NewTextMatcherMock(t)

	o := &Os{
		Base: Base{
			Storage:     mockStorage,
			Sops:        mockSops,
			TextMatcher: mockTestMatcher,
		},
	}

	result, err := o.Search("dir", SearchOptions{})

	var expected []FileMatch

	assert.Equal(t, expected, result)
	assert.Equal(t, errors.New("storage error"), err)
}

func TestOsSearchEmptyDirectory(t *testing.T) {
	mockStorage := NewStorageMock(t)
	mockStorage.
		On("ReadDir", "dir").
		Return([]string{}, nil)

	mockSops := NewSopsMock(t)
	mockTestMatcher := NewTextMatcherMock(t)

	o := &Os{
		Base: Base{
			Storage:     mockStorage,
			Sops:        mockSops,
			TextMatcher: mockTestMatcher,
		},
	}

	result, err := o.Search("dir", SearchOptions{})

	var expected []FileMatch

	assert.Equal(t, expected, result)
	assert.Nil(t, err)
}
