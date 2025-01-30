package scanner

import (
	"fmt"
	"github.com/alican-uelger/deep-scan/internal/matcher"
	"log/slog"
	"path/filepath"

	"github.com/alican-uelger/deep-scan/internal/sops"
	"github.com/alican-uelger/deep-scan/internal/storage"
)

type Git struct {
	Base
	Client GitClient
}

func NewGitlab(client GitClient) *Git {
	memStorage := storage.NewMem()
	return &Git{
		Base: Base{
			Storage:     memStorage,
			Sops:        sops.New(memStorage),
			TextMatcher: matcher.Text{},
		},
		Client: client,
	}
}

func (s *Git) Search(org string, options SearchOptions) ([]FileMatch, error) {
	var result []FileMatch
	projects, err := s.Client.ListGroupProjects(org)
	if err != nil {
		return result, err
	}
	for _, project := range projects {
		tree, err := s.Client.ListRepositoryTree(project.ID)
		if err != nil {
			return result, err
		}
		for _, treeEntry := range tree {
			// if its a tree/folder just skip
			if treeEntry.IsTree {
				continue
			}
			entry := filepath.Join(project.PathWithNamespace, treeEntry.Path)
			fileMatch := FileMatch{
				File: File{
					Name: filepath.Base(entry),
					Path: filepath.Dir(entry),
					Type: FILE, // TODO: detect if sops secret
				},
				Matches: nil,
			}
			rawContent, err := s.Client.GetRawFile(project, treeEntry.Path)
			if err != nil {
				return result, err
			}
			content := string(rawContent)
			if options.Sops && fileMatch.Type == SOPS_SECRET {
				content, err = s.decryptContent(fileMatch.File, rawContent)
				if err != nil {
					slog.Error(fmt.Sprintf("decrypt error: %s", err))
					continue
				}
			}
			ok, matches := s.filter(fileMatch.File, content, options)
			if !ok {
				continue
			}

			fileMatch.Matches = matches
			slog.Debug(fmt.Sprintf("found file: %s", entry))
			printFileMatch(fileMatch)
			result = append(result, fileMatch)
		}
	}
	return result, nil
}

func (s *Git) decryptContent(file File, rawContent []byte) (string, error) {
	content := string(rawContent)
	err := s.Storage.MkdirAll(file.Path)
	if err != nil {
		return content, fmt.Errorf("decrypt error: %s", err)
	}
	err = s.Storage.WriteFile(file.Path, rawContent)
	if err != nil {
		return content, fmt.Errorf("write file error: %s", err)
	}
	content, err = s.Sops.DecryptFile(file.Path)
	if err != nil {
		return content, fmt.Errorf("decrypt error: %s", err)
	}
	return content, nil
}
