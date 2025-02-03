package git

import (
	"context"
	"github.com/google/go-github/v50/github"
)

type githubClientWrapper struct {
	client *github.Client
}

func (w *githubClientWrapper) ListGroupProjects(org string, opts *github.RepositoryListByOrgOptions) ([]*github.Repository, *github.Response, error) {
	return w.client.Repositories.ListByOrg(context.Background(), org, opts)
}

func (w *githubClientWrapper) SearchProjects(name string, opts *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error) {
	return w.client.Search.Repositories(context.Background(), name, opts)
}

func (w *githubClientWrapper) GetRawFile(owner, repo, path string, opts *github.RepositoryContentGetOptions) ([]byte, *github.Response, error) {
	fileContent, _, _, err := w.client.Repositories.GetContents(context.Background(), owner, repo, path, opts)
	if err != nil {
		return nil, nil, err
	}
	content, err := fileContent.GetContent()
	return []byte(content), nil, err
}

func (w *githubClientWrapper) ListRepositoryTree(owner, repo string) ([]*github.TreeEntry, *github.Response, error) {
	tree, _, err := w.client.Git.GetTree(context.Background(), owner, repo, "HEAD", true)
	if tree == nil {
		return nil, nil, err
	}
	return tree.Entries, nil, err
}
