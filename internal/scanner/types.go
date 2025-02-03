package scanner

import (
	"github.com/alican-uelger/deep-scan/internal/git"
	"github.com/alican-uelger/deep-scan/internal/matcher"
)

type GitClient interface {
	ListGroupProjects(group string) ([]git.Project, error)
	ListRepositoryTree(project git.Project) ([]git.TreeNode, error)
	GetRawFile(project git.Project, path string) ([]byte, error)
}

type Storage interface {
	ReadFile(string) ([]byte, error)
	ReadDir(string) ([]string, error)
	IsDir(string) (bool, error)
	MkdirAll(string) error
	WriteFile(string, []byte) error
}

type Sops interface {
	DecryptFile(path string) (string, error)
}

type FileType string

// nolint
const (
	FILE        FileType = "FILE"
	SOPS_SECRET FileType = "SOPS_SECRET"
	SOPS_CONFIG FileType = "SOPS_CONFIG"
)

type File struct {
	Name string   `json:"name" yaml:"name"`
	Path string   `json:"path" yaml:"path"`
	Type FileType `json:"type" yaml:"type"`
}

type FileMatch struct {
	File
	Matches []matcher.MatchResult `json:"matches" yaml:"matches"`
}

type SearchOptions struct {
	Name                []string
	NameContains        []string
	NameRegex           []string
	Path                []string
	PathContains        []string
	PathRegex           []string
	Content             []string
	ContentRegex        []string
	Sops                bool
	SopsOnly            bool
	SopsKey             []string
	ExcludeName         []string
	ExcludeNameContains []string
	ExcludePath         []string
	ExcludePathContains []string
	ExcludeContent      []string
}
