package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"

	intdirenv "github.com/pitoniak32/dotenv_gsm/internal/direnv"
	"github.com/pitoniak32/dotenv_gsm/internal/version"

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

	slog.Debug("version info", "details", version.VersionInfo)

	args := os.Args

	var shell intdirenv.Shell
	var newenv intdirenv.Env
	var target string

	if len(args) == 2 && (args[1] == "--version" || args[1] == "version") {
		fmt.Println(fmt.Sprintf("%#+v", version.VersionInfo))
		os.Exit(0)
	}

	if len(args) > 1 {
		shell = intdirenv.DetectShell(args[1])
	} else {
		shell = intdirenv.Bash
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
		if !strings.Contains(value, "projects/") {
			slog.Error("found non GSM secretId value... ensure you only use gsm paths in this file... exiting", "key", key, "value", value)
			os.Exit(1)
		}
	}

	fetchSecrets(newenv)

	str, err := newenv.ToShell(shell)
	if err != nil {
		slog.Error("creating environment for shell failed", "shell", shell.Name(), "path", path, "err", err)
		os.Exit(1)
	}

	fmt.Println(str)
}

func fetchSecrets(secrets intdirenv.Env) (intdirenv.Env, error) {
	ctx := context.Background()

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		slog.Error("failed to create secret manager client: %v", err)
		os.Exit(1)
	}
	defer client.Close()

	type SecretResult struct {
		Key      string
		SecretId string
		Value    string
		Error    error
	}

	results := make(chan SecretResult, len(secrets))
	var wg sync.WaitGroup

	for key, secretId := range secrets {
		wg.Add(1)
		go func(secretId string) {
			defer wg.Done()
			if !strings.Contains(secretId, "versions/") {
				slog.Debug("adding '/versions/latest' to secret", "key", key, "secretId", secretId)
				secretId = fmt.Sprintf("%s/versions/latest", secretId)
			}
			req := &secretmanagerpb.AccessSecretVersionRequest{Name: secretId}

			slog.Debug("fetching gsm secret", "key", key, "secretId", secretId)
			result, err := client.AccessSecretVersion(ctx, req)
			if err != nil {
				results <- SecretResult{Key: key, SecretId: secretId, Value: "", Error: err}
				return
			}

			secretData := string(result.Payload.Data)
			results <- SecretResult{Key: key, SecretId: secretId, Value: secretData}
		}(secretId)
	}

	wg.Wait()
	close(results)

	for res := range results {
		if res.Error != nil {
			slog.Error("accessing secret failed", "key", res.Key, "secretId", res.SecretId, "err", res.Error)
		} else {
			slog.Debug("setting secret", "key", res.Key, "secretId", res.SecretId)
			secrets[res.Key] = res.Value
		}
	}

	return secrets, nil
}
