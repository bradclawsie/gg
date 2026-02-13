set shell := ["bash", "-c"]
set dotenv-load := true
set dotenv-filename := ".env-unit"
set dotenv-required := true

default:
    @just --list

pkgs:
    go install golang.org/x/vuln/cmd/govulncheck@latest
    go install github.com/securego/gosec/v2/cmd/gosec@latest
    go install github.com/google/capslock/cmd/capslock@latest

build:
    @mkdir -p out
    go build -o out/gg cmd/gg/main.go

lint:
    golangci-lint --timeout=24h run pkg/... && staticcheck ./... && go vet ./... && govulncheck ./... && gosec ./...

test:
    go test -race -v ./...

test-pkg PKG:
    pushd {{ PKG }} && go test -count=1 -race -v ./...
