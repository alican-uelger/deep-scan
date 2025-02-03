package cmd

import (
	"github.com/alican-uelger/deep-scan/internal/scanner"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const flagDir = "dir"

type Scanner interface {
	Search(string, scanner.SearchOptions) ([]scanner.FileMatch, error)
}

func NewOsScannerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "os",
	}
	addOsScannerFlags(cmd.PersistentFlags())
	bindFlags(cmd)
	osScanner := scanner.NewOs()
	cmd.AddCommand(NewSearchCmd(flagDir, osScanner))
	return cmd
}

func addOsScannerFlags(flagSet *pflag.FlagSet) {
	flagSet.StringP(flagDir, "d", ".", "The directory to scan")
}
