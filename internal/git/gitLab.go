package git

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strconv"

	"gitlab.com/gitlab-org/api/client-go"
)

type GitLabAPI interface {
	ListGroupProjects(group string, opts *gitlab.ListGroupProjectsOptions) ([]*gitlab.Project, *gitlab.Response, error)
	SearchProjects(project string, opts *gitlab.SearchOptions) ([]*gitlab.Project, *gitlab.Response, error)
	GetRawFile(project string, path string, opts *gitlab.GetRawFileOptions) ([]byte, *gitlab.Response, error)
	ListRepositoryTree(project string, opts *gitlab.ListTreeOptions) ([]*gitlab.TreeNode, *gitlab.Response, error)
}

type GitLab struct {
	client GitLabAPI
}

func NewGitLab(token, hostname string) (*GitLab, error) {
	if hostname == "" {
		hostname = "gitlab.com"
		slog.Debug(fmt.Sprintf("using default GitLab hostname: %s", hostname))
	}
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(hostname))
	if err != nil {
		return nil, err
	}
	return &GitLab{
		client: &gitlabClientWrapper{client},
	}, nil
}

func (g *GitLab) GetProjectByName(projectName string) (Project, error) {
	project := Project{}
	sOpts := gitlab.SearchOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: 1,
		},
	}
	slog.Debug(fmt.Sprintf("searching project: %s", projectName))
	projects, _, err := g.client.SearchProjects(projectName, &sOpts)
	if err != nil {
		return project, err
	}
	if len(projects) < 1 {
		return project, fmt.Errorf("project not found: %s", projectName)
	}
	project.Name = projects[0].Name
	project.ID = projects[0].ID
	project.PathWithNamespace = projects[0].PathWithNamespace
	return project, nil
}

func (g *GitLab) ListGroupProjects(group string) ([]Project, error) {
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
				Name:              project.Name,
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

func (g *GitLab) GetRawFile(project Project, path string) ([]byte, error) {
	slog.Debug(fmt.Sprintf("fetching raw file: %s", filepath.Join(project.PathWithNamespace, path)))
	content, _, err := g.client.GetRawFile(strconv.Itoa(project.ID), path, &gitlab.GetRawFileOptions{})
	return content, err
}

func (g *GitLab) ListRepositoryTree(project Project) ([]TreeNode, error) {
	var repoTreeNodes []TreeNode
	opts := &gitlab.ListTreeOptions{
		Recursive: gitlab.Bool(true),
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}
	for {
		slog.Debug(fmt.Sprintf("fetching tree for project: %s, page: %d", project.Name, opts.Page))
		treeNodes, resp, err := g.client.ListRepositoryTree(project.PathWithNamespace, opts)
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
