package scanner

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/alican-uelger/deep-scan/internal/sops"
	"github.com/alican-uelger/deep-scan/internal/storage"
)

type Storage interface {
	ReadFile(string) ([]byte, error)
	ReadDir(string) ([]string, error)
	IsDir(string) (bool, error)
	MkdirAll(string) error
	WriteFile(string, []byte) error
}

type Os struct {
	Storage Storage
	Sops    Sops
}

func NewOs() *Os {
	return &Os{
		Storage: storage.NewOs(),
		Sops:    sops.New(storage.NewMem()),
	}
}

func (s *Os) Search(dir string, options SearchOptions) ([]File, error) {
	result := []File{}

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
		file := File{
			Name: filepath.Base(entry),
			Path: filepath.Dir(entry),
			Type: FILE, // TODO: detect if sops secret
		}
		rawContent, err := s.Storage.ReadFile(entry)
		if err != nil {
			slog.Warn(fmt.Sprintf("reading file content failed %s - skipping %s and continueing", err, entry))
			continue
		}
		content := string(rawContent)
		if options.Sops && file.Type == SOPS_SECRET {
			content, err = decryptContent(file, rawContent, s.Storage, s.Sops)
			if err != nil {
				slog.Error(fmt.Sprintf("decrypt error: %s", err))
				continue
			}
		}
		if !filterFile(file, content, options) {
			continue
		}
		result = append(result, file)
	}
	return result, nil
}
