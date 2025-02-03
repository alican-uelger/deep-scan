package scanner

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/alican-uelger/deep-scan/internal/matcher"
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
		tree, err := s.Client.ListRepositoryTree(project)
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
					Type: FILE,
				},
				Matches: nil,
			}
			content := ""
			if isFileContentNeeded(options) {
				rawContent, err := s.Client.GetRawFile(project, treeEntry.Path)
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
					decryptedContent, err := s.decryptContent(fileMatch.File, rawContent)
					if err == nil {
						slog.Debug(fmt.Sprintf("found sops secret file: %s", entry))
						fileMatch.Type = SOPS_SECRET // set file type to sops secret
						content = decryptedContent   // set content to decrypted content
					}
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
	fileLocation := filepath.Join(file.Path, file.Name)
	content := string(rawContent)
	err := s.Storage.MkdirAll(fileLocation)
	if err != nil {
		return content, fmt.Errorf("decrypt: mkdirall error: %s", err)
	}
	err = s.Storage.WriteFile(fileLocation, rawContent)
	if err != nil {
		return content, fmt.Errorf("decrypt: write file error: %s", err)
	}
	content, err = s.Sops.DecryptFile(fileLocation)
	if err != nil {
		return content, fmt.Errorf("decrypt: decrypt error: %s", err)
	}
	return content, nil
}
