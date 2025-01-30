package scanner

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/alican-uelger/deep-scan/internal/matcher"
	"github.com/alican-uelger/deep-scan/internal/sops"
	"github.com/alican-uelger/deep-scan/internal/storage"
)

type Os struct {
	Base
}

func NewOs() *Os {
	return &Os{
		Base: Base{
			Storage:     storage.NewOs(),
			Sops:        sops.New(storage.NewMem()),
			TextMatcher: matcher.NewText(),
		},
	}
}

func (s *Os) Search(dir string, options SearchOptions) ([]FileMatch, error) {
	var result []FileMatch

	dirEntries, err := s.Storage.ReadDir(dir)
	if err != nil {
		return result, err
	}
	for _, entry := range dirEntries {
		isDir, err := s.Storage.IsDir(entry)
		if err != nil {
			slog.Warn(fmt.Sprintf("is directory function failed with err %s - skipping %s and continueing", err, entry))
			continue
		}
		if isDir {
			nestedFiles, err := s.Search(entry, options)
			if err != nil {
				slog.Warn(fmt.Sprintf("nested directory search failed with err %s - skipping %s and continueing", err, entry))
				continue
			}
			result = append(result, nestedFiles...)
			continue
		}
		fileMatch := FileMatch{
			File: File{
				Name: filepath.Base(entry),
				Path: filepath.Dir(entry),
				Type: FILE, // TODO: detect if sops secret
			},
			Matches: nil,
		}

		rawContent, err := s.Storage.ReadFile(entry)
		if err != nil {
			slog.Warn(fmt.Sprintf("reading file content failed %s - skipping %s and continueing", err, entry))
			continue
		}
		content := string(rawContent)
		if options.Sops && fileMatch.Type == SOPS_SECRET {
			content, err = s.decryptContent(fileMatch.File, rawContent)
			if err != nil {
				slog.Error(fmt.Sprintf("decrypt error: %s", err))
				continue
			}
		}
		// filter files
		ok, matches := s.filter(fileMatch.File, content, options)
		if !ok {
			continue
		}
		fileMatch.Matches = matches
		slog.Debug(fmt.Sprintf("found file: %s", entry))
		printFileMatch(fileMatch)
		result = append(result, fileMatch)
	}
	return result, nil
}

func (s *Os) decryptContent(file File, rawContent []byte) (string, error) {
	content := string(rawContent)
	content, err := s.Sops.DecryptFile(file.Path)
	if err != nil {
		return content, fmt.Errorf("decrypt error: %s", err)
	}
	return content, nil
}
