package cmd

import (
	"fmt"
	"github.com/alican-uelger/deep-scan/internal/git"
	"github.com/alican-uelger/deep-scan/internal/scanner"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

const (
	envGitHubToken = "GITHUB_TOKEN"
	envGitHubHost  = "GITHUB_HOST"
)

func NewGitHubScannerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "github",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return requireEnvs(envGitHubToken)
		},
	}
	githubClient, err := git.NewGitHub(os.Getenv(envGitHubToken), os.Getenv(envGitHubHost))
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating github client: %v", err))
	}
	githubScanner := scanner.NewGitlab(githubClient)
	cmd.AddCommand(NewSearchCmd(flagGitOrg, githubScanner))
	return cmd
}
