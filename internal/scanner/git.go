package scanner

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/alican-uelger/deep-scan/internal/sops"
	"github.com/alican-uelger/deep-scan/internal/storage"
)

type Git struct {
	Client  GitClient
	Storage Storage
	Sops    Sops
}

func NewGitlab(client GitClient) *Git {
	strg := storage.NewMem()
	return &Git{
		Storage: strg,
		Sops:    sops.New(strg),
		Client:  client,
	}
}

func (s *Git) Search(org string, options SearchOptions) ([]File, error) {
	result := []File{}
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
			file := File{
				Name: filepath.Base(entry),
				Path: filepath.Dir(entry),
				Type: FILE, // TODO: detect if sops secret
			}
			rawContent, err := s.Client.GetRawFile(project, treeEntry.Path)
			if err != nil {
				return result, err
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
	}
	return result, nil
}
