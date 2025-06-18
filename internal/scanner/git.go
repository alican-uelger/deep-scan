package scanner

import (
	"fmt"
	"github.com/alican-uelger/deep-scan/internal/git"
	"log/slog"
	"path/filepath"
	"sync"

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
	var mu sync.Mutex
	var wg sync.WaitGroup

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
			if treeEntry.IsTree {
				continue
			}

			wg.Add(1)
			go func(treeEntry git.TreeNode) {
				defer wg.Done()
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
						slog.Warn(fmt.Sprintf("reading file content failed %s - skipping %s and continuing", err, entry))
						return
					}
					content = string(rawContent)
					if options.Sops {
						ok, matches := s.filterSopsKey(content, options)
						if !ok {
							return
						}
						fileMatch.Matches = append(fileMatch.Matches, matches...)
						decryptedContent, err := s.decryptContent(fileMatch.File, rawContent)
						if err == nil {
							slog.Debug(fmt.Sprintf("found sops secret file: %s", entry))
							fileMatch.Type = SOPS_SECRET
							content = decryptedContent
						}
					}
				}
				ok, matches := s.filter(fileMatch.File, content, options)
				if !ok {
					return
				}
				fileMatch.Matches = matches
				slog.Debug(fmt.Sprintf("found file: %s", entry))
				if !options.LogLate {
					printFileMatch(fileMatch)
				}
				mu.Lock()
				result = append(result, fileMatch)
				mu.Unlock()
			}(treeEntry)
		}
	}
	wg.Wait()
	if options.LogLate {
		printFileMatches(result)
	}
	return result, nil
}

func (s *Git) decryptContent(file File, rawContent []byte) (string, error) {
	fileLocation := filepath.Join(file.Path, file.Name)
	content := string(rawContent)
	err := s.Storage.MkdirAll(file.Path)
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
