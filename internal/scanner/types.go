package scanner

import "github.com/alican-uelger/deep-scan/internal/git"

type GitClient interface {
	ListGroupProjects(group string) ([]git.Project, error)
	ListRepositoryTree(projectID any) ([]git.TreeNode, error)
	GetRawFile(project git.Project, path string) ([]byte, error)
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
	Name string
	Path string
	Type FileType
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
	SopsKey             []string
	ExcludeName         []string
	ExcludeNameContains []string
	ExcludeDir          []string
	ExcludeDirContains  []string
	ExcludeContent      []string
}
