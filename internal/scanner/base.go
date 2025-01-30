package scanner

import (
	"fmt"
	"github.com/alican-uelger/deep-scan/internal/matcher"
	"path/filepath"
)

type TextMatcher interface {
	Match(text, search string, searchType matcher.SearchType, context int) (bool, bool, []matcher.MatchResult)
}

type Base struct {
	TextMatcher TextMatcher
	Storage     Storage
	Sops        Sops
}

func (s *Base) filter(file File, content string, options SearchOptions) (bool, []matcher.MatchResult) {
	contextLength := 15

	var results []matcher.MatchResult

	// excludes
	if len(options.ExcludeName) > 0 {
		for _, name := range options.ExcludeName {
			_, exactMatch, matches := s.TextMatcher.Match(file.Name, name, matcher.TextSearch, contextLength)
			results = append(results, matches...)
			if exactMatch {
				return false, results
			}
		}
	}
	if len(options.ExcludeNameContains) > 0 {
		for _, nameContains := range options.ExcludeNameContains {
			matched, _, matches := s.TextMatcher.Match(file.Name, nameContains, matcher.TextSearch, contextLength)
			results = append(results, matches...)
			if matched {
				return false, results
			}
		}
	}
	if len(options.ExcludePath) > 0 {
		for _, p := range options.ExcludePath {
			_, exactMatch, matches := s.TextMatcher.Match(file.Path, p, matcher.TextSearch, contextLength)
			results = append(results, matches...)
			if exactMatch {
				return false, results
			}
		}
	}
	if len(options.ExcludePathContains) > 0 {
		for _, pContains := range options.ExcludePathContains {
			matched, _, matches := s.TextMatcher.Match(file.Name, pContains, matcher.TextSearch, contextLength)
			results = append(results, matches...)
			if matched {
				return false, results
			}
		}
	}
	if len(options.ExcludeContent) > 0 {
		for _, c := range options.ExcludeContent {
			matched, _, matches := s.TextMatcher.Match(content, c, matcher.TextSearch, contextLength)
			results = append(results, matches...)
			if matched {
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

func printFileMatch(fileMatch FileMatch) {
	fmt.Println("+----------------------------------------+")
	fmt.Println("Match:\t" + filepath.Join(fileMatch.Path, fileMatch.Name))
	for i, m := range fileMatch.Matches {
		fmt.Printf("\tLine:%d, ColStart:%d, ColEnd:%d\n", m.Line, m.StartCol, m.EndCol)
		fmt.Printf("\t'%s'\n", m.CompressedFormattedSnippet)
		if i < len(fileMatch.Matches)-1 {
			fmt.Println()
		}
	}
	fmt.Println("+----------------------------------------+")
}
