package git

type TreeNode struct {
	IsTree bool // is folder
	Path   string
	Type   string
}

type Project struct {
	ID                int
	PathWithNamespace string // Github html_url
}
