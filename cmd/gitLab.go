package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/alican-uelger/deep-scan/internal/git"
	"github.com/alican-uelger/deep-scan/internal/scanner"
	"github.com/spf13/cobra"
)

const (
	envGitlabHost  = "GITLAB_HOST"
	envGitlabToken = "GITLAB_TOKEN"
)

func NewGitLabScannerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "gitlab",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return requireEnvs(envGitlabToken)
		},
	}
	gitlabClient, err := git.NewGitLab(os.Getenv(envGitlabToken), os.Getenv(envGitlabHost))
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating gitlab client: %v", err))
	}
	gitlabScanner := scanner.NewGitlab(gitlabClient)
	cmd.AddCommand(NewSearchCmd(flagGitOrg, gitlabScanner))
	return cmd
}

func requireEnvs(envs ...string) error {
	for _, env := range envs {
		if os.Getenv(env) == "" {
			return fmt.Errorf("set required ENV 'export %s=<YOUR_%s>'", env, env)
		}
	}
	return nil
}
