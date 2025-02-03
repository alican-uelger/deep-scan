package matcher

import (
	"regexp"
	"strings"
)

type MatchResult struct {
	Line                       int    `json:"line"`
	StartCol                   int    `json:"startCol"`
	EndCol                     int    `json:"endCol"`
	Snippet                    string `json:"snippet"`
	FormattedSnippet           string `json:"formattedSnippet"`
	CompressedFormattedSnippet string `json:"compressedFormattedSnippet"`
	ExactMatch                 bool   `json:"exactMatch"`
}

type SearchType string

const (
	TextSearch  SearchType = "text"
	RegexSearch SearchType = "regex"
)

type Text struct{}

func NewText() Text {
	return Text{}
}

func (t Text) Match(text, search string, searchType SearchType, context int) (bool, bool, []MatchResult) {
	if exactMatch, results := checkExactMatch(text, search, searchType); exactMatch {
		return true, true, results
	}

	return findMatches(text, search, searchType, context)
}

func checkExactMatch(text, search string, searchType SearchType) (bool, []MatchResult) {
	if searchType == TextSearch && search == text {
		return true, []MatchResult{createExactMatchResult(text, search, searchType)}
	}

	if searchType == RegexSearch {
		re, _ := regexp.Compile(search)
		if re.MatchString(text) && re.FindString(text) == text {
			return true, []MatchResult{createExactMatchResult(text, search, searchType)}
		}
	}

	return false, nil
}

func createExactMatchResult(text, search string, searchType SearchType) MatchResult {
	formattedSnippet := formatSnippet(text, search, searchType)
	compressedSnippet := compressSnippet(formattedSnippet)
	return MatchResult{
		ExactMatch:                 true,
		Line:                       1,
		StartCol:                   1,
		EndCol:                     len(text),
		Snippet:                    text,
		FormattedSnippet:           formattedSnippet,
		CompressedFormattedSnippet: compressedSnippet,
	}
}

func findMatches(text, search string, searchType SearchType, context int) (bool, bool, []MatchResult) {
	var results []MatchResult
	start := 0

	for {
		index, matchLength := getNextMatch(text, search, searchType, start)
		if index == -1 {
			break
		}

		results = append(results, createMatchResult(text, search, searchType, index, matchLength, context))
		start = index + matchLength

		if start >= len(text) {
			break
		}
	}

	return len(results) > 0, false, results
}

func getNextMatch(text, search string, searchType SearchType, start int) (int, int) {
	if searchType == RegexSearch {
		re := regexp.MustCompile(search)
		loc := re.FindStringIndex(text[start:])
		if loc == nil {
			return -1, 0
		}
		return loc[0] + start, loc[1] - loc[0]
	}

	index := strings.Index(text[start:], search)
	if index == -1 {
		return -1, 0
	}

	return index + start, len(search)
}

func createMatchResult(text, search string, searchType SearchType, index, matchLength, context int) MatchResult {
	endIndex := index + matchLength
	startContext := max(0, index-context)
	endContext := min(len(text), endIndex+context)

	snippet := text[startContext:endContext]
	if startContext > 0 {
		snippet = "..." + snippet
	}
	if endContext < len(text) {
		snippet = snippet + "..."
	}

	line, startCol, endCol := calculatePosition(text, index, endIndex)
	formattedSnippet := formatSnippet(snippet, search, searchType)
	compressedSnippet := compressSnippet(formattedSnippet)

	return MatchResult{
		Line:                       line,
		StartCol:                   startCol,
		EndCol:                     endCol,
		Snippet:                    snippet,
		FormattedSnippet:           formattedSnippet,
		CompressedFormattedSnippet: compressedSnippet,
	}
}

func calculatePosition(text string, index, endIndex int) (int, int, int) {
	line := strings.Count(text[:index], "\n") + 1
	lastNewline := strings.LastIndex(text[:index], "\n")
	if lastNewline == -1 {
		lastNewline = 0
	} else {
		lastNewline++
	}

	startCol := index - lastNewline + 1
	endCol := endIndex - lastNewline
	return line, startCol, endCol
}

func compressSnippet(snippet string) string {
	snippet = strings.ReplaceAll(snippet, "\n", "\\n")
	snippet = strings.ReplaceAll(snippet, "\t", "\\t")
	return strings.Join(strings.Fields(snippet), " ")
}

func formatSnippet(snippet, search string, searchType SearchType) string {
	green, gray, reset := "\033[32m", "\033[90m", "\033[0m"

	if searchType == TextSearch {
		if index := strings.Index(snippet, search); index != -1 {
			return gray + snippet[:index] + green + snippet[index:index+len(search)] + gray + snippet[index+len(search):] + reset
		}
		return gray + snippet + reset
	}

	re, err := regexp.Compile(search)
	if err != nil {
		return gray + snippet + reset
	}

	var result strings.Builder
	lastIndex := 0
	matches := re.FindAllStringIndex(snippet, -1)

	for _, match := range matches {
		start, end := match[0], match[1]
		result.WriteString(gray + snippet[lastIndex:start] + green + snippet[start:end] + gray)
		lastIndex = end
	}

	result.WriteString(snippet[lastIndex:] + reset)
	return result.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
