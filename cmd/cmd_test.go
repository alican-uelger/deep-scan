//go:build unit

package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCmd(t *testing.T) {
	cmd := NewRootCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "", cmd.Use)
	assert.True(t, cmd.TraverseChildren)
}

func TestBindFlags(t *testing.T) {
	cmd := &cobra.Command{}
	addRootFlags(cmd.PersistentFlags())
	bindFlags(cmd)
	assert.Equal(t, "INFO", viper.GetString(flagLogLevel))
}
