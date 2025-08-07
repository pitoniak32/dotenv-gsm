build:
  goreleaser build --snapshot --clean
  

build-run:
  goreleaser build --snapshot --clean && \
  ./dist/dotenv_gsm_darwin_arm64_v8.0/dotenv_gsm bash .env.secret

build-all:
  ./hack/build.sh darwin arm64
  ./hack/build.sh linux amd64 

# Test all packages
test flags="": 
  go test {{flags}} ./...

# Lint all packages
lint: fmt
  golangci-lint run

# Format all packages
fmt:
  ./hack/fmt.sh
