package main

import (
	"fmt"
	"log/slog"

	"github.com/alican-uelger/deep-scan/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	err := rootCmd.Execute()
	if err != nil {
		slog.Error(fmt.Sprintf("Error executing root cmd: %v", err))
	}
}
