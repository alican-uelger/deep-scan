package scanner

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"sync"

	"github.com/alican-uelger/deep-scan/internal/matcher"
	"github.com/alican-uelger/deep-scan/internal/sops"
	"github.com/alican-uelger/deep-scan/internal/storage"
)

type Os struct {
	Base
}

func NewOs() *Os {
	osStorage := storage.NewOs()
	return &Os{
		Base: Base{
			Storage:     osStorage,
			Sops:        sops.New(osStorage),
			TextMatcher: matcher.NewText(),
		},
	}
}

func (s *Os) Search(dir string, options SearchOptions) ([]FileMatch, error) {
	var result []FileMatch
	var mu sync.Mutex
	var wg sync.WaitGroup

	dirEntries, err := s.Storage.ReadDir(dir)
	if err != nil {
		return result, err
	}

	for _, entry := range dirEntries {
		wg.Add(1)
		go func(entry string) {
			defer wg.Done()
			isDir, err := s.Storage.IsDir(entry)
			if err != nil {
				slog.Warn(fmt.Sprintf("is directory function failed with err %s - skipping %s and continuing", err, entry))
				return
			}
			if isDir {
				nestedFiles, err := s.Search(entry, options)
				if err != nil {
					slog.Warn(fmt.Sprintf("nested directory search failed with err %s - skipping %s and continuing", err, entry))
					return
				}
				mu.Lock()
				result = append(result, nestedFiles...)
				mu.Unlock()
				return
			}
			fileMatch := FileMatch{
				File: File{
					Name: filepath.Base(entry),
					Path: filepath.Dir(entry),
					Type: FILE,
				},
				Matches: nil,
			}
			var matches []matcher.MatchResult
			ok, filterFileMatches := s.filterFile(fileMatch.File, options)
			if !ok {
				return
			}
			matches = append(matches, filterFileMatches...)

			content := ""
			if isFileContentNeeded(options) {
				rawContent, err := s.Storage.ReadFile(entry)
				if err != nil {
					slog.Warn(fmt.Sprintf("reading file content failed %s - skipping %s and continuing", err, entry))
					return
				}
				content = string(rawContent)
				if options.Sops {
					ok, _ := s.filterSopsKey(content, options)
					if !ok {
						return
					}
					decryptedContent, err := s.decryptContent(fileMatch.File)
					if err == nil {
						slog.Debug(fmt.Sprintf("found sops secret file: %s", entry))
						fileMatch.Type = SOPS_SECRET
						content = decryptedContent
					}
				}
			}

			// when sops-only is enabled, only search for sops files
			if options.SopsOnly && fileMatch.Type != SOPS_SECRET {
				return
			}

			ok, contentMatches := s.filterContent(content, options)
			if !ok {
				return
			}
			matches = append(matches, contentMatches...)
			fileMatch.Matches = matches
			slog.Debug(fmt.Sprintf("found file: %s", entry))
			if !options.LogLate {
				printFileMatch(fileMatch)
			}
			mu.Lock()
			result = append(result, fileMatch)
			mu.Unlock()
		}(entry)
	}
	wg.Wait()
	if options.LogLate {
		printFileMatches(result)
	}
	return result, nil
}

func (s *Os) decryptContent(file File) (string, error) {
	content, err := s.Sops.DecryptFile(filepath.Join(file.Path, file.Name))
	if err != nil {
		return content, fmt.Errorf("decrypt error: %s", err)
	}
	return content, nil
}
