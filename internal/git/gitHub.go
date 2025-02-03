package git

import (
	"context"
	"fmt"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
	"log/slog"
	"path/filepath"
)

type GitHubAPI interface {
	SearchProjects(name string, opts *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error)
	GetRawFile(owner, repo, path string, opts *github.RepositoryContentGetOptions) ([]byte, *github.Response, error)
	ListRepositoryTree(owner, repo string) ([]*github.TreeEntry, *github.Response, error)
	ListGroupProjects(org string, opts *github.RepositoryListByOrgOptions) ([]*github.Repository, *github.Response, error)
}

type GitHub struct {
	client GitHubAPI
}

func NewGitHub(token, hostname string) (*GitHub, error) {
	if hostname == "" {
		hostname = "github.com"
		slog.Debug(fmt.Sprintf("using default GitHub hostname: %s", hostname))
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &GitHub{
		client: &githubClientWrapper{client},
	}, nil
}

func (g *GitHub) GetProjectByName(name string) (Project, error) {
	project := Project{}
	sOpts := &github.SearchOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 1,
		},
	}
	slog.Debug(fmt.Sprintf("searching project: %s", name))
	r, _, err := g.client.SearchProjects(name, sOpts)
	if err != nil {
		return project, err
	}

	if len(r.Repositories) < 1 {
		return project, fmt.Errorf("project not found: %s", name)
	}
	project.ID = int(r.Repositories[0].GetID())
	project.PathWithNamespace = r.Repositories[0].GetFullName()
	return project, nil
}

func (g *GitHub) ListGroupProjects(group string) ([]Project, error) {
	var gitProjects []Project
	opts := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}
	slog.Debug(fmt.Sprintf("fetching projects for group: %s", group))
	repos, _, err := g.client.ListGroupProjects(group, opts)
	if err != nil {
		return nil, err
	}
	for _, repo := range repos {
		gitProjects = append(gitProjects, Project{
			Name:              repo.GetName(),
			ID:                int(repo.GetID()),
			PathWithNamespace: repo.GetFullName(),
		})
	}
	return gitProjects, nil
}

func (g *GitHub) GetRawFile(project Project, path string) ([]byte, error) {
	slog.Debug(fmt.Sprintf("fetching raw file: %s", filepath.Join(project.PathWithNamespace, path)))
	content, _, err := g.client.GetRawFile(project.Owner(), project.Name, path, &github.RepositoryContentGetOptions{})
	return content, err
}

func (g *GitHub) ListRepositoryTree(project Project) ([]TreeNode, error) {
	var repoTreeNodes []TreeNode
	slog.Debug(fmt.Sprintf("fetching tree for project: %v", project.Name))
	tree, _, err := g.client.ListRepositoryTree(project.Owner(), project.Name)
	if err != nil {
		return repoTreeNodes, err
	}
	for _, treeNode := range tree {
		repoTreeNodes = append(repoTreeNodes, TreeNode{
			IsTree: treeNode.GetType() == "tree",
			Path:   treeNode.GetPath(),
			Type:   treeNode.GetType(),
		})
	}
	return repoTreeNodes, nil
}
