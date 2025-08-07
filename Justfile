build platform="darwin" arch="arm64":
  ./hack/build.sh {{platform}} {{arch}}

build-run platform="darwin" arch="arm64":
  ./hack/build.sh {{platform}} {{arch}}
  ./bin/dotenv_gsm-{{platform}}-{{arch}} bash .env.secret

build-all:
  ./hack/build.sh darwin arm64
  ./hack/build.sh linux amd64 

# Run the entry point
run:
  go run . bash .env.secret

# Test all packages
test flags="": 
  go test {{flags}} ./...

# Lint all packages
lint: fmt
  golangci-lint run

# Format all packages
fmt:
  ./hack/fmt.sh
