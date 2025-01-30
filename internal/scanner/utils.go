package scanner

import (
	"fmt"
	"log/slog"
	"regexp"
	"slices"
	"strings"
)

func strContains(ss []string, s string) bool {
	for _, substr := range ss {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

func regMatchContains(regs []string, s string) bool {
	for _, reg := range regs {
		matched, err := regexp.MatchString(reg, s)
		if err != nil {
			slog.Error(fmt.Sprintf("regex error: %s", err))
		}
		if matched {
			return true
		}
	}
	return false
}

// nolint:gocognit
func filterFile(file File, content string, options SearchOptions) bool {
	for _, excludeName := range options.ExcludeName {
		if file.Name == excludeName {
			return false
		}
	}
	for _, excludeNameContains := range options.ExcludeNameContains {
		if strings.Contains(file.Name, excludeNameContains) {
			return false
		}
	}
	for _, excludeDir := range options.ExcludeDir {
		if file.Path == excludeDir {
			return false
		}
	}
	for _, excludeDirContains := range options.ExcludeDirContains {
		if strings.Contains(file.Path, excludeDirContains) {
			return false
		}
	}
	for _, excludeContent := range options.ExcludeContent {
		if strings.Contains(content, excludeContent) {
			return false
		}
	}

	if len(options.Name) > 0 && !slices.Contains(options.Name, file.Name) {
		return false
	}
	if len(options.NameContains) > 0 && !strContains(options.NameContains, file.Name) {
		return false
	}
	if len(options.NameRegex) > 0 && !regMatchContains(options.NameRegex, file.Name) {
		return false
	}

	if len(options.Path) > 0 && !slices.Contains(options.Path, file.Path) {
		return false
	}
	if len(options.PathContains) > 0 && !strContains(options.PathContains, file.Path) {
		return false
	}
	if len(options.PathRegex) > 0 && !regMatchContains(options.PathRegex, file.Path) {
		return false
	}

	if len(options.Content) > 0 && !strContains(options.Content, content) {
		return false
	}

	if len(options.ContentRegex) > 0 && !regMatchContains(options.ContentRegex, content) {
		return false
	}

	return true
}

func decryptContent(file File, rawContent []byte, store Storage, s Sops) (string, error) {
	content := string(rawContent)
	err := store.MkdirAll(file.Path)
	if err != nil {
		return content, fmt.Errorf("decrypt error: %w", err)
	}
	err = store.WriteFile(file.Path, rawContent)
	if err != nil {
		return content, fmt.Errorf("write file error: %w", err)
	}
	content, err = s.DecryptFile(file.Path)
	if err != nil {
		return content, fmt.Errorf("decrypt error: %w", err)
	}
	return content, nil
}
