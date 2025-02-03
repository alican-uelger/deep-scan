package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/alican-uelger/deep-scan/internal/scanner"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"log/slog"
	"os"
	"strings"
)

func NewSearchCmd(flagStartingPoint string, scanner Scanner) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "search",
		RunE: search(flagStartingPoint, scanner),
		PreRun: func(cmd *cobra.Command, args []string) {
			bindFlags(cmd)
		},
	}
	addGitScannerFlags(cmd.PersistentFlags())
	addSearchFlags(cmd.PersistentFlags())
	addOutputFLags(cmd.PersistentFlags())
	bindFlags(cmd)
	return cmd
}

const (
	JSON string = "json"
	YAML string = "yaml"
)

func search(flagStartingPoint string, scanner Scanner) RunE {
	return func(_ *cobra.Command, _ []string) error {
		options := searchOptions()
		startingPoint := viper.GetString(flagStartingPoint)
		slog.Debug(fmt.Sprintf("running search with %s=%s, options: %v", flagStartingPoint, startingPoint, options))
		files, err := scanner.Search(startingPoint, options)
		if err != nil {
			slog.Error(fmt.Sprintf("Error searching for files: %v", err))
			return err
		}
		slog.Debug(fmt.Sprintf("found %d files", len(files)))
		o := viper.GetString(flagOutput)
		if o != "" {
			err = output(o, files)
			if err != nil {
				slog.Error(fmt.Sprintf("Error outputting files: %v", err))
				return err
			}
		}
		return nil
	}
}

func output(outputType string, files []scanner.FileMatch) error {
	switch strings.ToLower(outputType) {
	case JSON:
		return jsonOutput(files)
	case YAML:
		return yamlOutput(files)
	default:
		return fmt.Errorf("unsupported output type: %s", outputType)
	}
}

func jsonOutput(files []scanner.FileMatch) error {
	filesJson, err := json.MarshalIndent(files, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile("output.json", filesJson, 0644)
	if err != nil {
		return err
	}
	slog.Info("Output written to output.json")
	return nil
}

func yamlOutput(files []scanner.FileMatch) error {
	filesYaml, err := yaml.Marshal(files)
	if err != nil {
		return err
	}
	err = os.WriteFile("output.yaml", filesYaml, 0644)
	if err != nil {
		return err
	}
	slog.Info("Output written to output.yaml")
	return nil
}

func addOutputFLags(flagSet *pflag.FlagSet) {
	flagSet.String(flagOutput, "", "Output to file (JSON, YAML)")
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
	flagSet.Bool(flagSopsOnly, false, "Search for files that are only SOPS-encrypted")
	flagSet.StringSlice(flagSopsKey, []string{}, "Search for files encrypted with a specific key")

	// exclude filename flags
	flagSet.StringSlice(flagExcludeName, []string{}, "Exclude files with specific names (exact match)")
	flagSet.StringSlice(flagExcludeNameContains, []string{}, "Exclude files with names containing this string")

	// exclude file path flags
	flagSet.StringSlice(flagExcludePath, []string{}, "Exclude specific directories (exact match)")
	flagSet.StringSlice(flagExcludePathContains, []string{}, "Exclude directories containing this string")
	flagSet.StringSlice(flagExcludeContent, []string{}, "Exclude files containing specific content")
}
