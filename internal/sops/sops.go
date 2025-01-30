package sops

import (
	"fmt"
	"path/filepath"
	"strings"
)

type FileType = string

const (
	YAML   FileType = "yaml"
	DotEnv FileType = "dotenv"
)

type SopsAPI interface {
	DecryptFile([]byte, string) ([]byte, error)
}

type Storage interface {
	ReadFile(string) ([]byte, error)
}

type Sops struct {
	Storage Storage
	Client  SopsAPI
}

func New(storage Storage) *Sops {
	return &Sops{
		Storage: storage,
		Client:  &sopsCLientWrapper{},
	}
}

func (s *Sops) DecryptFile(path string) (string, error) {
	fileType := getFileType(path)
	content, err := s.Storage.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("could not read secret file: %w", err)
	}
	decryptedContent, err := s.Client.DecryptFile(content, fileType)
	if err != nil {
		return "", fmt.Errorf("could not decrypt secret file: %w", err)
	}
	return string(decryptedContent), nil
}

func getFileType(path string) FileType {
	fileType := strings.Replace(filepath.Ext(path), ".", "", 1)
	switch fileType {
	case "env":
		return DotEnv
	case "yml":
		return YAML
	}
	return fileType
}
