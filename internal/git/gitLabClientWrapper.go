package git

import gitlab "gitlab.com/gitlab-org/api/client-go"

/*
gitlabClientWrapper is a wrapper around gitlab.Client to make it easier to mock the gitlab.Client in tests
*/

type gitlabClientWrapper struct {
	client *gitlab.Client
}

func (w *gitlabClientWrapper) ListGroupProjects(group string, opts *gitlab.ListGroupProjectsOptions) ([]*gitlab.Project, *gitlab.Response, error) {
	return w.client.Groups.ListGroupProjects(group, opts)
}

/*
project string: The name or ID of the project
*/

func (w *gitlabClientWrapper) SearchProjects(project string, opts *gitlab.SearchOptions) ([]*gitlab.Project, *gitlab.Response, error) {
	return w.client.Search.Projects(project, opts)
}

func (w *gitlabClientWrapper) GetRawFile(project string, path string, opts *gitlab.GetRawFileOptions) ([]byte, *gitlab.Response, error) {
	return w.client.RepositoryFiles.GetRawFile(project, path, opts)
}

func (w *gitlabClientWrapper) ListRepositoryTree(project string, opts *gitlab.ListTreeOptions) ([]*gitlab.TreeNode, *gitlab.Response, error) {
	return w.client.Repositories.ListTree(project, opts)
}
