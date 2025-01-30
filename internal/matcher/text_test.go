//go:build unit

package matcher

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchContentText(t *testing.T) {
	m := NewText()

	tests := []struct {
		name               string
		content            string
		search             string
		context            int
		expectedLen        int
		expectedExactMatch bool
	}{
		{
			name:               "Two match",
			content:            "This is a sample text with several words. This text is for testing.",
			search:             "text",
			context:            10,
			expectedLen:        2,
			expectedExactMatch: false,
		},
		{
			name:               "No match",
			content:            "This is a sample text with several words. This text is for testing.",
			search:             "nonexistent",
			context:            10,
			expectedLen:        0,
			expectedExactMatch: false,
		},
		{
			name:               "One match",
			content:            "Another example with different text.",
			search:             "example",
			context:            5,
			expectedLen:        1,
			expectedExactMatch: false,
		},
		{
			name:               "Multiple matches",
			content:            "Multiple text occurrences in this text.",
			search:             "text",
			context:            5,
			expectedLen:        2,
			expectedExactMatch: false,
		},
		{
			name:               "Exact match",
			content:            "This Content will match 100%.",
			search:             "This Content will match 100%.",
			context:            5,
			expectedLen:        1,
			expectedExactMatch: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, results := m.Match(tt.content, tt.search, TextSearch, tt.context)
			assert.Len(t, results, tt.expectedLen)
			for _, result := range results {
				assert.Equal(t, tt.expectedExactMatch, result.ExactMatch)
				fmt.Println(result)
			}
		})
	}
}

func TestMatchContentRegex(t *testing.T) {
	m := NewText()

	tests := []struct {
		name               string
		content            string
		search             string
		context            int
		expectedLen        int
		expectedExactMatch bool
	}{
		{
			name:               "Two match",
			content:            "This is a sample text with several words. This text is for testing.",
			search:             `\btext\b`,
			context:            10,
			expectedLen:        2,
			expectedExactMatch: false,
		},
		{
			name:               "No match",
			content:            "This is a sample text with several words. This text is for testing.",
			search:             `\bnonexistent\b`,
			context:            10,
			expectedLen:        0,
			expectedExactMatch: false,
		},
		{
			name:               "One match",
			content:            "Another example with different text.",
			search:             `\bexample\b`,
			context:            5,
			expectedLen:        1,
			expectedExactMatch: false,
		},
		{
			name:               "Multiple matches",
			content:            "Multiple text occurrences in this text.",
			search:             `\btext\b`,
			context:            5,
			expectedLen:        2,
			expectedExactMatch: false,
		},
		{
			name:               "Exact match",
			content:            "This Content will match 100%.",
			search:             `.*`,
			context:            5,
			expectedLen:        1,
			expectedExactMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, results := m.Match(tt.content, tt.search, RegexSearch, tt.context)
			assert.Len(t, results, tt.expectedLen)
			for _, result := range results {
				assert.Equal(t, tt.expectedExactMatch, result.ExactMatch)
				fmt.Println(result)
			}
		})
	}
}
