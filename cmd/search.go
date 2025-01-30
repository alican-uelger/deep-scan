package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log/slog"
)

func NewSearchCmd(flagStartingPoint string, scanner Scanner) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "search",
		RunE: search(flagStartingPoint, scanner),
	}
	addSearchFlags(cmd.PersistentFlags())
	bindFlags(cmd)
	return cmd
}

func search(flagStartingPoint string, scanner Scanner) RunE {
	return func(_ *cobra.Command, _ []string) error {
		options := searchOptions()
		slog.Debug(fmt.Sprintf("running search with options: %v", options))
		files, err := scanner.Search(viper.GetString(flagStartingPoint), options)
		if err != nil {
			slog.Error(fmt.Sprintf("Error searching for files: %v", err))
			return err
		}
		slog.Debug(fmt.Sprintf("found %d files", len(files)))
		return nil
	}
}

func addSearchFlags(flagSet *pflag.FlagSet) {
	// filename flags
	flagSet.StringSliceP(flagName, "n", []string{}, "Search for files with specific names (exact match)")
	flagSet.StringSlice(flagNameContains, []string{}, "Search for files with names containing this string")
	flagSet.StringSlice(flagNameRegex, []string{}, "Search for files with names matching this regex")

	// file path flags
	flagSet.StringSliceP(flagPath, "p", []string{}, "Search in specific directories (exact match)")
	flagSet.StringSlice(flagPathContains, []string{}, "Search in directories containing this string")
	flagSet.StringSlice(flagPathRegex, []string{}, "Search in directories matching this regex")

	// content flags
	flagSet.StringSliceP(flagContent, "c", []string{}, "Search for files containing specific content")
	flagSet.StringSlice(flagContentRegex, []string{}, "Search for files containing content matching this regex")

	// sops flags
	flagSet.BoolP(flagSops, "s", false, "Search for SOPS-encrypted files")
	flagSet.StringSlice(flagSopsKey, []string{}, "Specify SOPS keys for decryption")

	// exclude filename flags
	flagSet.StringSlice(flagExcludeName, []string{}, "Exclude files with specific names (exact match)")
	flagSet.StringSlice(flagExcludeNameContains, []string{}, "Exclude files with names containing this string")

	// exclude file path flags
	flagSet.StringSlice(flagExcludePath, []string{}, "Exclude specific directories (exact match)")
	flagSet.StringSlice(flagExcludePathContains, []string{}, "Exclude directories containing this string")
	flagSet.StringSlice(flagExcludeContent, []string{}, "Exclude files containing specific content")
}
