package scanner

import (
	"fmt"
	"github.com/alican-uelger/deep-scan/internal/matcher"
	"path/filepath"
	"sync"
)

const contextLength = 15

type TextMatcher interface {
	Match(text, search string, searchType matcher.SearchType, context int) (bool, bool, []matcher.MatchResult)
}

type Base struct {
	TextMatcher TextMatcher
	Storage     Storage
	Sops        Sops
}

func (s *Base) filterFile(file File, options SearchOptions) (bool, []matcher.MatchResult) {
	var results []matcher.MatchResult

	// excludes
	if len(options.ExcludeName) > 0 {
		for _, name := range options.ExcludeName {
			_, exactMatch, matches := s.TextMatcher.Match(file.Name, name, matcher.TextSearch, contextLength)
			if exactMatch {
				results = append(results, matches...)
				return false, results
			}
		}
	}
	if len(options.ExcludeNameContains) > 0 {
		for _, nameContains := range options.ExcludeNameContains {
			matched, _, matches := s.TextMatcher.Match(file.Name, nameContains, matcher.TextSearch, contextLength)
			if matched {
				results = append(results, matches...)
				return false, results
			}
		}
	}
	if len(options.ExcludePath) > 0 {
		for _, p := range options.ExcludePath {
			_, exactMatch, matches := s.TextMatcher.Match(file.Path, p, matcher.TextSearch, contextLength)
			if exactMatch {
				results = append(results, matches...)
				return false, results
			}
		}
	}
	if len(options.ExcludePathContains) > 0 {
		for _, pContains := range options.ExcludePathContains {
			matched, _, matches := s.TextMatcher.Match(file.Path, pContains, matcher.TextSearch, contextLength)
			if matched {
				results = append(results, matches...)
				return false, results
			}
		}
	}

	// name
	if len(options.Name) > 0 {
		for _, name := range options.Name {
			_, exactMatch, matches := s.TextMatcher.Match(file.Name, name, matcher.TextSearch, contextLength)
			results = append(results, matches...)
			if !exactMatch {
				return false, results
			}
		}
	}
	if len(options.NameContains) > 0 {
		for _, nameContains := range options.NameContains {
			matched, _, matches := s.TextMatcher.Match(file.Name, nameContains, matcher.TextSearch, contextLength)
			results = append(results, matches...)
			if !matched {
				return false, results
			}
		}
	}

	if len(options.NameRegex) > 0 {
		for _, nameRegex := range options.NameRegex {
			matched, _, matches := s.TextMatcher.Match(file.Name, nameRegex, matcher.RegexSearch, contextLength)
			results = append(results, matches...)
			if !matched {
				return false, results
			}
		}
	}

	// path
	if len(options.Path) > 0 {
		for _, p := range options.Path {
			_, exactMatch, matches := s.TextMatcher.Match(file.Path, p, matcher.TextSearch, contextLength)
			results = append(results, matches...)
			if !exactMatch {
				return false, results
			}
		}
	}
	if len(options.PathContains) > 0 {
		for _, pContains := range options.PathContains {
			matched, _, matches := s.TextMatcher.Match(file.Path, pContains, matcher.TextSearch, contextLength)
			results = append(results, matches...)
			if !matched {
				return false, results
			}
		}
	}

	if len(options.PathRegex) > 0 {
		for _, pRegex := range options.PathRegex {
			matched, _, matches := s.TextMatcher.Match(file.Path, pRegex, matcher.RegexSearch, contextLength)
			results = append(results, matches...)
			if !matched {
				return false, results
			}
		}
	}

	return true, results
}

func (s *Base) filterSopsKey(content string, options SearchOptions) (bool, []matcher.MatchResult) {
	var results []matcher.MatchResult
	if len(options.SopsKey) > 0 {
		for _, c := range options.SopsKey {
			matched, _, matches := s.TextMatcher.Match(content, c, matcher.TextSearch, contextLength)
			results = append(results, matches...)
			if !matched {
				return false, results
			}
		}
	}
	return true, results
}

func (s *Base) filterContent(content string, options SearchOptions) (bool, []matcher.MatchResult) {
	var results []matcher.MatchResult

	// excludes
	if len(options.ExcludeContent) > 0 {
		for _, c := range options.ExcludeContent {
			matched, _, matches := s.TextMatcher.Match(content, c, matcher.TextSearch, contextLength)
			if matched {
				results = append(results, matches...)
				return false, results
			}
		}
	}

	// content
	if len(options.Content) > 0 {
		for _, c := range options.Content {
			matched, _, matches := s.TextMatcher.Match(content, c, matcher.TextSearch, contextLength)
			results = append(results, matches...)
			if !matched {
				return false, results
			}
		}
	}
	if len(options.ContentRegex) > 0 {
		for _, cRegex := range options.ContentRegex {
			matched, _, matches := s.TextMatcher.Match(content, cRegex, matcher.RegexSearch, contextLength)
			results = append(results, matches...)
			if !matched {
				return false, results
			}
		}
	}

	return true, results
}

func (s *Base) decryptContent(file File, encryptedContent []byte) (string, error) {
	decryptedContent := string(encryptedContent)
	decryptedContent, err := s.Sops.DecryptFile(filepath.Join(file.Path, file.Name))
	if err != nil {
		return decryptedContent, fmt.Errorf("decrypt error: %s", err)
	}
	return decryptedContent, nil
}

var printMu sync.Mutex

func printFileMatches(fileMatches []FileMatch) {
	go func() {
		printMu.Lock()
		defer printMu.Unlock()
		result := ""
		for _, fileMatch := range fileMatches {
			result += buildFileMatchOutput(fileMatch)
		}
		fmt.Print(result)
	}()
}

func printFileMatch(fileMatch FileMatch) {
	go func() {
		printMu.Lock()
		defer printMu.Unlock()
		fmt.Print(buildFileMatchOutput(fileMatch))
	}()

}

func buildFileMatchOutput(fileMatch FileMatch) string {
	result := "+----------------------------------------+\n"
	result += "Match:\t" + filepath.Join(fileMatch.Path, fileMatch.Name) + "\n"
	for i, m := range fileMatch.Matches {
		result += fmt.Sprintf("\tLine:%d, ColStart:%d, ColEnd:%d\n", m.Line, m.StartCol, m.EndCol)
		result += fmt.Sprintf("\t'%s'\n", m.CompressedFormattedSnippet)
		if i < len(fileMatch.Matches)-1 {
			result += "\n"
		}
	}
	result += "+----------------------------------------+\n\n"
	return result
}

// This prevents the unnecessary downloading/reading files by checking if the content is needed
func isFileContentNeeded(options SearchOptions) bool {
	if len(options.Content) > 0 {
		return true
	}
	if len(options.ContentRegex) > 0 {
		return true
	}
	if options.Sops {
		return true
	}
	if len(options.ExcludeContent) > 0 {
		return true
	}
	return false
}
