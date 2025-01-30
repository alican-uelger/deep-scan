//go:build unit

package scanner

import (
	"testing"

	"github.com/alican-uelger/deep-scan/internal/matcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBaseFilter(t *testing.T) {
	mockTextMatcher := NewTextMatcherMock(t)
	mockStorage := NewStorageMock(t)
	mockSops := NewSopsMock(t)

	base := &Base{
		TextMatcher: mockTextMatcher,
		Storage:     mockStorage,
		Sops:        mockSops,
	}

	file := File{Name: "file.txt", Path: "org/repo", Type: FILE}
	content := "file content"

	tests := []struct {
		name     string
		options  SearchOptions
		expected bool
	}{
		{
			name: "ExcludeName",
			options: SearchOptions{
				ExcludeName: []string{"file.txt"},
			},
			expected: false,
		},
		{
			name: "ExcludeNameContains",
			options: SearchOptions{
				ExcludeNameContains: []string{"file"},
			},
			expected: false,
		},
		{
			name: "ExcludePath",
			options: SearchOptions{
				ExcludePath: []string{"org/repo"},
			},
			expected: false,
		},
		{
			name: "ExcludePathContains",
			options: SearchOptions{
				ExcludePathContains: []string{"repo"},
			},
			expected: false,
		},
		{
			name: "ExcludeContent",
			options: SearchOptions{
				ExcludeContent: []string{"file content"},
			},
			expected: false,
		},
		{
			name: "Name",
			options: SearchOptions{
				Name: []string{"file.txt"},
			},
			expected: true,
		},
		{
			name: "NameContains",
			options: SearchOptions{
				NameContains: []string{"file"},
			},
			expected: true,
		},
		{
			name: "NameRegex",
			options: SearchOptions{
				NameRegex: []string{`file\.txt`},
			},
			expected: true,
		},
		{
			name: "Path",
			options: SearchOptions{
				Path: []string{"org/repo"},
			},
			expected: true,
		},
		{
			name: "PathContains",
			options: SearchOptions{
				PathContains: []string{"repo"},
			},
			expected: true,
		},
		{
			name: "PathRegex",
			options: SearchOptions{
				PathRegex: []string{`org/repo`},
			},
			expected: true,
		},
		{
			name: "Content",
			options: SearchOptions{
				Content: []string{"file content"},
			},
			expected: true,
		},
		{
			name: "ContentRegex",
			options: SearchOptions{
				ContentRegex: []string{`file content`},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTextMatcher.On("Match", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(true, true, []matcher.MatchResult{}).Maybe()
			result, _ := base.filter(file, content, tt.options)
			assert.Equal(t, tt.expected, result)
		})
	}
}
