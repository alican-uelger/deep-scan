//go:build unit

package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewSearchCmd(t *testing.T) {
	scanner := NewScannerMock(t)
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
		flagSops, flagSopsContentBeforeDecryption,
		flagExcludeName, flagExcludeNameContains,
		flagExcludePath, flagExcludePathContains, flagExcludeContent,
	}

	for _, flag := range flags {
		assert.NotNil(t, cmd.PersistentFlags().Lookup(flag), "Flag %s should be present", flag)
	}
}

func TestSearch_NeitherOrgNorProject(t *testing.T) {
	t.Cleanup(viper.Reset)

	scanner := NewScannerMock(t)
	cmd := NewSearchCmd(flagGitOrg, scanner)
	cmd.SetArgs([]string{})
	err := cmd.Execute()

	assert.ErrorContains(t, err, "provide at least one of")
}

func TestSearch_BothOrgAndProject(t *testing.T) {
	t.Cleanup(viper.Reset)

	scanner := NewScannerMock(t)
	cmd := NewSearchCmd(flagGitOrg, scanner)
	cmd.SetArgs([]string{"--org", "myorg", "--project", "owner/repo"})
	err := cmd.Execute()

	assert.ErrorContains(t, err, "mutually exclusive")
}
