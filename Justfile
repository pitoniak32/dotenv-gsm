build:
  goreleaser build --snapshot --clean

build-run:
  goreleaser build --snapshot --clean && \
  echo "export TESTING=$'projects/test/secrets/testing';export TEST=$'projects/test/secrets/test/versions/1';" | LOG_LEVEL=debug ./dist/dotenv-gsm_darwin_arm64_v8.0/dotenv-gsm -

install-local location="~/.local/bin/": build
  cp ./dist/dotenv-gsm_darwin_arm64_v8.0/dotenv-gsm {{location}}

# Test all packages
test flags="": 
  go test {{flags}} ./...

# Lint all packages
lint: fmt
  golangci-lint run

# Format all packages
fmt:
  ./hack/fmt.sh

