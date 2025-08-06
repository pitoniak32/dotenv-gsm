package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	// secretmanager "cloud.google.com/go/secretmanager/apiv1"
	// "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"

	"github.com/direnv/direnv/v2/pkg/dotenv"
)

func parseLogLevel(env string) slog.Level {
	switch strings.ToLower(env) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelWarn // fallback default
	}
}

func main() {
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				Level: parseLogLevel(os.Getenv("LOG_LEVEL")),
			}),
		),
	)

	args := os.Args

	var shell Shell
	var newenv Env
	var target string

	if len(args) > 1 {
		shell = DetectShell(args[1])
	} else {
		shell = Bash
	}

	if len(args) > 2 {
		target = args[2]
	}

	var data []byte
	data, err := os.ReadFile(target)
	if err != nil {
		slog.Error("reading env file failed", "file", target, "err", err)
		os.Exit(1)
	}

	// Set PWD env var to the directory the .env file resides in. This results
	// in the least amount of surprise, as a dotenv file is most often defined
	// in the same directory it's loaded from, so referring to PWD should match
	// the directory of the .env file.
	path, err := filepath.Abs(target)
	if err != nil {
		slog.Error("finding absolute directory of env file failed", "file", target, "err", err)
		os.Exit(1)
	}
	if err := os.Setenv("PWD", filepath.Dir(path)); err != nil {
		slog.Error("finding pwd of env file failed", "path", path, "err", err)
		os.Exit(1)
	}

	newenv, err = dotenv.Parse(string(data))
	if err != nil {
		slog.Error("parsing env file failed", "path", path, "err", err)
		os.Exit(1)
	}

	for key, value := range newenv {
		if strings.Contains(value, "projects/") {
			slog.Warn("mock fetching gsm secret", "key", key, "value", value)
			newenv[key] = "***"
		}
	}

	str, err := newenv.ToShell(shell)
	if err != nil {
		slog.Error("creating environment for shell failed", "shell", shell.Name(), "path", path, "err", err)
		os.Exit(1)
	}

	fmt.Println(str)
}
