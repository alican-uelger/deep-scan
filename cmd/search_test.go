//go:build unit

package cmd

import (
	"testing"

	"github.com/alican-uelger/deep-scan/test/mocks"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNewSearchCmd(t *testing.T) {
	scanner := mocks.NewScanner(t)
	cmd := NewSearchCmd("test", scanner)
	assert.NotNil(t, cmd)
	assert.Equal(t, "search", cmd.Use)
}

func TestAddSearchFlags(t *testing.T) {
	cmd := &cobra.Command{}
	addSearchFlags(cmd.PersistentFlags())

	flags := []string{
		flagName, flagNameContains, flagNameRegex,
		flagPath, flagPathContains, flagPathRegex,
		flagContent, flagContentRegex,
		flagSops, flagSopsKey,
		flagExcludeName, flagExcludeNameContains,
		flagExcludePath, flagExcludePathContains, flagExcludeContent,
	}

	for _, flag := range flags {
		assert.NotNil(t, cmd.PersistentFlags().Lookup(flag), "Flag %s should be present", flag)
	}
}
