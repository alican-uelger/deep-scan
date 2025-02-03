//go:build unit

package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNewOsScannerCmd(t *testing.T) {
	cmd := NewOsScannerCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "os", cmd.Use)
}

func TestAddOsScannerFlags(t *testing.T) {
	cmd := &cobra.Command{}
	addOsScannerFlags(cmd.PersistentFlags())
	assert.NotNil(t, cmd.PersistentFlags().Lookup(flagDir))
}
