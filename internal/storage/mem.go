/*
Copyright © 2024 BizHub servicebroker@deutschebahn.com
*/
package storage

import (
	"io/fs"
	"path/filepath"

	"gomodules.xyz/memfs"
)

type Mem struct {
	fs memfs.FS
}

func NewMem() *Mem {
	return &Mem{
		fs: *memfs.New(),
	}
}

func (s *Mem) MkdirAll(path string) error {
	return s.fs.MkdirAll(path, 0777)
}

func (s *Mem) WriteFile(path string, conent []byte) error {
	return s.fs.WriteFile(path, conent, 0644)
}

func (s *Mem) ReadFile(path string) ([]byte, error) {
	return fs.ReadFile(&s.fs, path)
}

func (s *Mem) ReadDir(path string) ([]string, error) {
	dirEntries, err := fs.ReadDir(&s.fs, path)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(dirEntries))
	for _, entry := range dirEntries {
		result = append(result, filepath.Join(path, entry.Name()))
	}
	return result, nil
}

func (s *Mem) IsDir(path string) (bool, error) {
	info, err := fs.Stat(&s.fs, path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), err
}
