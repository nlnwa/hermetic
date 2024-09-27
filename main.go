package main

import (
	"hermetic/cmd"
	"log/slog"
	"os"
)

func main() {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	if err := cmd.NewRootCommand().Execute(); err != nil {
		slog.Error("failed to execute command, got error:", err)
		os.Exit(1)
	}
}
