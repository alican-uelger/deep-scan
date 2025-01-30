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

	o := &Os{
		Storage: mockStorage,
		Sops:    mockSops,
	}

	result, err := o.Search("dir", SearchOptions{})
	expected := []File{
		{Name: "file.txt", Path: "dir", Type: FILE},
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

	o := &Os{
		Storage: mockStorage,
		Sops:    mockSops,
	}

	result, err := o.Search("dir", SearchOptions{})

	expected := []File{}

	assert.Equal(t, expected, result)
	assert.Equal(t, errors.New("storage error"), err)
}

func TestOsSearchEmptyDirectory(t *testing.T) {
	mockStorage := NewStorageMock(t)
	mockStorage.
		On("ReadDir", "dir").
		Return([]string{}, nil)

	mockSops := NewSopsMock(t)

	o := &Os{
		Storage: mockStorage,
		Sops:    mockSops,
	}

	result, err := o.Search("dir", SearchOptions{})

	expected := []File{}

	assert.Equal(t, expected, result)
	assert.Nil(t, err)
}

func TestOsSearchWithOptions(t *testing.T) {
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

	o := &Os{
		Storage: mockStorage,
		Sops:    mockSops,
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
				ExcludeDir: []string{"dir"},
			},
			expected: []File{},
		},
		{
			name: "ExcludeDirContains",
			options: SearchOptions{
				ExcludeDirContains: []string{"dir"},
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
				{Name: "file.txt", Path: "dir", Type: FILE},
			},
		},
		{
			name: "NameContains",
			options: SearchOptions{
				NameContains: []string{"file"},
			},
			expected: []File{
				{Name: "file.txt", Path: "dir", Type: FILE},
			},
		},
		{
			name: "NameRegex",
			options: SearchOptions{
				NameRegex: []string{`file\.txt`},
			},
			expected: []File{
				{Name: "file.txt", Path: "dir", Type: FILE},
			},
		},
		{
			name: "Path",
			options: SearchOptions{
				Path: []string{"dir"},
			},
			expected: []File{
				{Name: "file.txt", Path: "dir", Type: FILE},
			},
		},
		{
			name: "PathContains",
			options: SearchOptions{
				PathContains: []string{"dir"},
			},
			expected: []File{
				{Name: "file.txt", Path: "dir", Type: FILE},
			},
		},
		{
			name: "PathRegex",
			options: SearchOptions{
				PathRegex: []string{`dir`},
			},
			expected: []File{
				{Name: "file.txt", Path: "dir", Type: FILE},
			},
		},
		{
			name: "Content",
			options: SearchOptions{
				Content: []string{"file content"},
			},
			expected: []File{
				{Name: "file.txt", Path: "dir", Type: FILE},
			},
		},
		{
			name: "ContentRegex",
			options: SearchOptions{
				ContentRegex: []string{`file content`},
			},
			expected: []File{
				{Name: "file.txt", Path: "dir", Type: FILE},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := o.Search("dir", tt.options)
			assert.Equal(t, tt.expected, result)
			assert.Nil(t, err)
		})
	}
}
