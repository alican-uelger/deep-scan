package cmd

import (
	"log/slog"
	"strings"

	"github.com/alican-uelger/deep-scan/internal/scanner"
	"github.com/spf13/viper"
)

const (
	flagName                = "name"
	flagNameContains        = "name-contains"
	flagNameRegex           = "name-regex"
	flagPath                = "path"
	flagPathContains        = "path-contains"
	flagPathRegex           = "path-regex"
	flagContent             = "content"
	flagContentRegex        = "content-regex"
	flagSops                = "sops"
	flagSopsOnly            = "sops-only"
	flagSopsKey             = "sops-key"
	flagExcludeName         = "exclude-name"
	flagExcludeNameContains = "exclude-name-contains"
	flagExcludePath         = "exclude-path"
	flagExcludePathContains = "exclude-path-contains"
	flagExcludeContent      = "exclude-content"
	flagLogLate             = "log-late"
)

const (
	flagOutput = "output"
)

func searchOptions() scanner.SearchOptions {
	return scanner.SearchOptions{
		Name:                viper.GetStringSlice(flagName),
		NameContains:        viper.GetStringSlice(flagNameContains),
		NameRegex:           viper.GetStringSlice(flagNameRegex),
		Path:                viper.GetStringSlice(flagPath),
		PathContains:        viper.GetStringSlice(flagPathContains),
		PathRegex:           viper.GetStringSlice(flagPathRegex),
		Content:             viper.GetStringSlice(flagContent),
		ContentRegex:        viper.GetStringSlice(flagContentRegex),
		Sops:                viper.GetBool(flagSops),
		SopsOnly:            viper.GetBool(flagSopsOnly),
		SopsKey:             viper.GetStringSlice(flagSopsKey),
		ExcludeName:         viper.GetStringSlice(flagExcludeName),
		ExcludeNameContains: viper.GetStringSlice(flagExcludeNameContains),
		ExcludePath:         viper.GetStringSlice(flagExcludePath),
		ExcludePathContains: viper.GetStringSlice(flagExcludePathContains),
		ExcludeContent:      viper.GetStringSlice(flagExcludeContent),
		LogLate:             viper.GetBool(flagLogLate),
	}
}

func logLevel() slog.Level {
	defaultLogLevel := slog.LevelInfo
	level := viper.GetString(flagLogLevel)

	if level == "" {
		return defaultLogLevel
	}
	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return defaultLogLevel
	}
}
