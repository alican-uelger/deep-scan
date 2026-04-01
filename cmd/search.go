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
			name := viper.GetString(flagOutputName)
			err = output(o, name, files)
			if err != nil {
				slog.Error(fmt.Sprintf("Error outputting files: %v", err))
				return err
			}
		}
		return nil
	}
}

func output(outputType string, name string, files []scanner.FileMatch) error {
	switch strings.ToLower(outputType) {
	case JSON:
		return jsonOutput(name, files)
	case YAML:
		return yamlOutput(name, files)
	default:
		return fmt.Errorf("unsupported output type: %s", outputType)
	}
}

func jsonOutput(name string, files []scanner.FileMatch) error {
	if name == "" {
		name = "output.json"
	}
	filesJson, err := json.MarshalIndent(files, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(name, filesJson, 0644)
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("Output written to %s", name))
	return nil
}

func yamlOutput(name string, files []scanner.FileMatch) error {
	if name == "" {
		name = "output.yaml"
	}
	filesYaml, err := yaml.Marshal(files)
	if err != nil {
		return err
	}
	err = os.WriteFile(name, filesYaml, 0644)
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("Output written to %s", name))
	return nil
}

func addOutputFLags(flagSet *pflag.FlagSet) {
	flagSet.String(flagOutput, "", "Output to file (JSON, YAML)")
	flagSet.String(flagOutputName, "", "Output file name (default: output.json / output.yaml)")
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
	flagSet.StringSlice(flagSopsContentBeforeDecryption, []string{}, "Search for content in SOPS-encrypted files before decryption")

	// output flags
	flagSet.Bool(flagNoSnippets, false, "Suppress match snippets in output")

	// exclude filename flags
	flagSet.StringSlice(flagExcludeName, []string{}, "Exclude files with specific names (exact match)")
	flagSet.StringSlice(flagExcludeNameContains, []string{}, "Exclude files with names containing this string")

	// exclude file path flags
	flagSet.StringSlice(flagExcludePath, []string{}, "Exclude specific directories (exact match)")
	flagSet.StringSlice(flagExcludePathContains, []string{}, "Exclude directories containing this string")
	flagSet.StringSlice(flagExcludeContent, []string{}, "Exclude files containing specific content")

	flagSet.Bool(flagLogLate, false, "This flag will log the results after the search is complete. This is useful for large searches, when you want to be as fast as possible.")
}
