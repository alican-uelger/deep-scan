package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/alican-uelger/deep-scan/internal/git"
	"github.com/alican-uelger/deep-scan/internal/scanner"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const flagGitOrg = "org"

const (
	envGitlabHost  = "GITLAB_HOST"
	envGitlabToken = "GITLAB_TOKEN"
)

func NewGitlabScannerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "gitlab",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return requireEnvs(envGitlabHost, envGitlabToken)
		},
	}
	addGitlabScannerFlags(cmd.PersistentFlags())
	err := cmd.MarkPersistentFlagRequired(flagGitOrg)
	if err != nil {
		slog.Error(fmt.Sprintf("Error marking flag as required: %v", err))
	}
	bindFlags(cmd)
	gitlabClient, err := git.NewGitlab(os.Getenv(envGitlabToken), os.Getenv(envGitlabHost))
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating gitlab client: %v", err))
	}
	gitlabScanner := scanner.NewGitlab(gitlabClient)
	cmd.AddCommand(NewSearchCmd(flagGitOrg, gitlabScanner))
	return cmd
}

func addGitlabScannerFlags(flagSet *pflag.FlagSet) {
	flagSet.StringP(flagGitOrg, "o", "", "The git org")
}

func requireEnvs(envs ...string) error {
	for _, env := range envs {
		if os.Getenv(env) == "" {
			return fmt.Errorf("set required ENV 'export %s=<YOUR_%s>'", env, env)
		}
	}
	return nil
}
