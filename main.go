package main

import (
	"log/slog"
	"os"

	"github.com/nlnwa/hermetic/cmd"
)

func main() {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	if err := cmd.NewRootCommand().Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
