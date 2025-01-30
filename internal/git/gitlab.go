package git

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"gitlab.com/gitlab-org/api/client-go"
)

type GitlabAPI interface {
	SearchProjects(name string, opts *gitlab.SearchOptions) ([]*gitlab.Project, *gitlab.Response, error)
	GetRawFile(projectID int, path string, opts *gitlab.GetRawFileOptions) ([]byte, *gitlab.Response, error)
	ListRepositoryTree(projectID any, opts *gitlab.ListTreeOptions) ([]*gitlab.TreeNode, *gitlab.Response, error)
	ListGroupProjects(group string, opts *gitlab.ListGroupProjectsOptions) ([]*gitlab.Project, *gitlab.Response, error)
}

type Gitlab struct {
	client GitlabAPI
}

func NewGitlab(token, host string) (*Gitlab, error) {
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(host))
	if err != nil {
		return nil, err
	}
	return &Gitlab{
		client: &gitlabClientWrapper{client},
	}, nil
}

func (g *Gitlab) GetProjectByName(name string) (Project, error) {
	project := Project{}
	sOpts := gitlab.SearchOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: 1,
		},
	}
	slog.Debug(fmt.Sprintf("searching project: %s", name))
	projects, _, err := g.client.SearchProjects(name, &sOpts)
	if err != nil {
		return project, err
	}
	if len(projects) < 1 {
		return project, fmt.Errorf("project not found: %s", name)
	}
	project.ID = projects[0].ID
	project.PathWithNamespace = projects[0].PathWithNamespace
	return project, nil
}

func (g *Gitlab) GetRawFile(project Project, path string) ([]byte, error) {
	slog.Debug(fmt.Sprintf("fetching raw file: %s", filepath.Join(project.PathWithNamespace, path)))
	content, _, err := g.client.GetRawFile(project.ID, path, &gitlab.GetRawFileOptions{})
	return content, err
}

func (g *Gitlab) ListRepositoryTree(projectID any) ([]TreeNode, error) {
	var repoTreeNodes []TreeNode
	opts := &gitlab.ListTreeOptions{
		Recursive: gitlab.Bool(true),
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}
	for {
		slog.Debug(fmt.Sprintf("fetching tree for project: %v, page: %d", projectID, opts.Page))
		treeNodes, resp, err := g.client.ListRepositoryTree(projectID, opts)
		if err != nil {
			return repoTreeNodes, err
		}
		for _, treeNode := range treeNodes {
			repoTreeNodes = append(repoTreeNodes, TreeNode{
				IsTree: treeNode.Type == "tree",
				Path:   treeNode.Path,
				Type:   treeNode.Type,
			})
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return repoTreeNodes, nil
}

func (g *Gitlab) ListGroupProjects(group string) ([]Project, error) {
	var gitProjects []Project
	opts := &gitlab.ListGroupProjectsOptions{
		IncludeSubGroups: gitlab.Bool(true),
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}
	for {
		slog.Debug(fmt.Sprintf("fetching projects for group: %s, page: %d", group, opts.Page))
		projects, resp, err := g.client.ListGroupProjects(group, opts)
		if err != nil {
			return nil, err
		}
		for _, project := range projects {
			gitProjects = append(gitProjects, Project{
				ID:                project.ID,
				PathWithNamespace: project.PathWithNamespace,
			})
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return gitProjects, nil
}
