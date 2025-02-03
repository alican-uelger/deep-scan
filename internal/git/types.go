package git

import (
	"strings"
)

type TreeNode struct {
	IsTree bool // is folder
	Path   string
	Type   string
}

type Project struct {
	Name              string
	ID                int
	PathWithNamespace string // Github html_url
}

func (p *Project) Owner() string {
	parts := strings.Split(p.PathWithNamespace, "/")
	if len(parts) < 1 {
		return ""
	}
	return parts[0]
}
