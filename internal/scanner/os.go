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
		content := ""
		// read file content if needed
		if isFileContentNeeded(options) {
			rawContent, err := s.Storage.ReadFile(entry)
			if err != nil {
				slog.Warn(fmt.Sprintf("reading file content failed %s - skipping %s and continueing", err, entry))
				continue
			}
			content = string(rawContent) // set content to raw content
			// check if sops is enabled
			if options.Sops {
				// check the encrypted content for sops-key
				ok, matches := s.filterSopsKey(content, options)
				if !ok {
					continue
				}
				fileMatch.Matches = append(fileMatch.Matches, matches...)
				// detect weather the file is sops-secret file which is encrypted or not
				decryptedContent, err := s.decryptContent(fileMatch.File)
				if err == nil {
					slog.Debug(fmt.Sprintf("found sops secret file: %s", entry))
					fileMatch.Type = SOPS_SECRET // set file type to sops secret
					content = decryptedContent   // set content to decrypted content
				}
			}
		}
		// filter files
		ok, matches := s.filter(fileMatch.File, content, options)
		if !ok {
			continue
		}
		fileMatch.Matches = append(fileMatch.Matches, matches...)
		slog.Debug(fmt.Sprintf("found file: %s", entry))
		if !options.LogLate {
			printFileMatch(fileMatch)
		}
		result = append(result, fileMatch)
	}
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
