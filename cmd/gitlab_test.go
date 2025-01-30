//go:build unit

package cmd

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGitlabScannerCmd(t *testing.T) {
	err := os.Setenv(envGitlabHost, "test_host")
	require.NoError(t, err)
	err = os.Setenv(envGitlabToken, "test_token")
	require.NoError(t, err)
	cmd := NewGitlabScannerCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "gitlab", cmd.Use)
}

func TestAddGitlabScannerFlags(t *testing.T) {
	cmd := &cobra.Command{}
	addGitlabScannerFlags(cmd.PersistentFlags())
	assert.NotNil(t, cmd.PersistentFlags().Lookup(flagGitOrg))
}

func TestRequireEnvs(t *testing.T) {
	err := os.Setenv(envGitlabHost, "test_host")
	require.NoError(t, err)
	err = os.Setenv(envGitlabToken, "test_token")
	require.NoError(t, err)
	err = requireEnvs(envGitlabHost, envGitlabToken)
	require.NoError(t, err)
}
