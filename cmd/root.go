package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	flagLogLevel = "log-level"
	flagGitOrg   = "org"
)

// NewRootCmd creates the root command for the CLI application.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			bindFlags(cmd)
			setupLogger()
		},
	}
	cmd.TraverseChildren = true // activates the PersistentPreRun for all subcommands
	addRootFlags(cmd.PersistentFlags())
	bindFlags(cmd)
	cmd.AddCommand(NewOsScannerCmd())
	cmd.AddCommand(NewGitLabScannerCmd())
	cmd.AddCommand(NewGitHubScannerCmd())
	return cmd
}

func addGitScannerFlags(flagSet *pflag.FlagSet) {
	flagSet.StringP(flagGitOrg, "o", "", "The git org to scan")
}

func addRootFlags(flagSet *pflag.FlagSet) {
	flagSet.StringP(flagLogLevel, "l", "INFO", "Set the log level (DEBUG, INFO, WARN, ERROR)")
}

func bindFlags(cmd *cobra.Command) {
	err := viper.BindPFlags(cmd.PersistentFlags())
	if err != nil {
		slog.Error(fmt.Sprintf("failed to bind flags: %v", err))
	}
}

func setupLogger() {
	l := logLevel()
	options := &slog.HandlerOptions{Level: l}
	handler := slog.NewJSONHandler(os.Stdout, options)
	slog.SetDefault(slog.New(handler))
	slog.Info("using loglevel " + l.String())
}
