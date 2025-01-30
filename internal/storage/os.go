package storage

import (
	"os"
	"path/filepath"
)

type Os struct{}

func NewOs() *Os {
	return &Os{}
}

func (s *Os) MkdirAll(path string) error {
	return os.MkdirAll(path, 0777)
}

func (s *Os) WriteFile(path string, conent []byte) error {
	return os.WriteFile(path, conent, 0600)
}

func (s *Os) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (s *Os) ReadDir(path string) ([]string, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(dirEntries))
	for _, entry := range dirEntries {
		result = append(result, filepath.Join(path, entry.Name()))
	}
	return result, nil
}

func (s *Os) IsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), err
}
