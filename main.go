package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"

	slogenv "github.com/cbrewster/slog-env"
	"log/slog"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

func main() {

	slog.SetDefault(slog.New(slogenv.NewHandler(slog.NewTextHandler(os.Stderr, nil), slogenv.WithEnvVarName("LOG_LEVEL"))))

	slog.Debug("version info", "details", VersionInfo)

	if len(os.Args) == 2 && (os.Args[1] == "--version" || os.Args[1] == "version") {
		fmt.Println(fmt.Sprintf("%#+v", VersionInfo))
		os.Exit(0)
	}

	var environmentString string

	if len(os.Args) > 1 && os.Args[1] == "-" {
		slog.Debug("reading input from stdin")
		environmentBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			slog.Error("encountered error reading stdin", "err", err)
			os.Exit(1)
		}
		environmentString = string(environmentBytes)
	} else if len(os.Args) > 1 && strings.Contains(os.Args[1], "export ") {
		slog.Debug("reading input from first arg")
		environmentString = os.Args[1]
	} else {
		fmt.Fprintln(os.Stderr, "Usage: direnv-gsm [ - | <exports string> ]")
		return
	}

	newenv := parseExportString(string(environmentString))

	for key, value := range newenv {
		if !strings.Contains(value, "projects/") {
			slog.Error("found non GSM secretId value... ensure you only use gsm paths in this file... exiting", "key", key, "value", value)
			os.Exit(1)
		}
	}

	for key, value := range fetchSecrets(newenv) {
		environmentString = strings.ReplaceAll(environmentString, key, value)
		slog.Debug("replaced all instances of key with secret value", "key", key)
	}

	fmt.Println(environmentString)
}

type Env = map[string]string

func parseExportString(input string) Env {
	// Regex to match: export KEY=$'VALUE';
	re := regexp.MustCompile(`export\s+(\w+)=\$'([^']*)';`)
	matches := re.FindAllStringSubmatch(input, -1)

	result := make(map[string]string)
	for _, match := range matches {
		if len(match) == 3 {
			value := match[2]
			result[value] = value
		}
	}
	return result
}

func fetchSecrets(secrets Env) Env {
	ctx := context.Background()

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		slog.Error("failed to create secret manager client: %v", "err", err)
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

	return secrets
}

