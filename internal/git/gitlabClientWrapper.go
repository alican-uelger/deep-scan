package git

import gitlab "gitlab.com/gitlab-org/api/client-go"

/**
gitlabClientWrapper is a wrapper around gitlab.Client to make it easier to mock the gitlab.Client in tests
*/

type gitlabClientWrapper struct {
	client *gitlab.Client
}

func (w *gitlabClientWrapper) SearchProjects(name string, opts *gitlab.SearchOptions) ([]*gitlab.Project, *gitlab.Response, error) {
	return w.client.Search.Projects(name, opts)
}

func (w *gitlabClientWrapper) GetRawFile(projectID int, path string, opts *gitlab.GetRawFileOptions) ([]byte, *gitlab.Response, error) {
	return w.client.RepositoryFiles.GetRawFile(projectID, path, opts)
}

func (w *gitlabClientWrapper) ListRepositoryTree(projectID any, opts *gitlab.ListTreeOptions) ([]*gitlab.TreeNode, *gitlab.Response, error) {
	return w.client.Repositories.ListTree(projectID, opts)
}

func (w *gitlabClientWrapper) ListGroupProjects(group string, opts *gitlab.ListGroupProjectsOptions) ([]*gitlab.Project, *gitlab.Response, error) {
	return w.client.Groups.ListGroupProjects(group, opts)
}
