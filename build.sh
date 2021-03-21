#!/usr/bin/env bash
set -e
set -o xtrace

cd "$(dirname $0)"


go build -o build/private/decrypt cmd/decrypt/decrypt.go # Building personal decryptor

go build -o build/public/encrypt cmd/encrypt/encrypt.go # Building encryptor for the current platform


# Build encryptor for other platforms
GOOS=linux GOARCH=amd64 go build -ldflags "-w" -o build/public/encrypt-linux-amd64 cmd/encrypt/encrypt.go
# GOOS=linux GOARCH=arm64 go build -ldflags "-w" -o build/public/encrypt-linux-arm64 cmd/encrypt/encrypt.go
GOOS=darwin GOARCH=amd64 go build -ldflags "-w" -o build/public/encrypt-darwin-amd64 cmd/encrypt/encrypt.go
GOOS=darwin GOARCH=arm64 go build -ldflags "-w" -o build/public/encrypt-darwin-arm64 cmd/encrypt/encrypt.go
# GOOS=windows GOARCH=amd64 go build -ldflags "-w" -o build/public/encrypt-windows-amd64 cmd/encrypt/encrypt.go

GOOS=js GOARCH=wasm go build -o web/browser-me-a-secret/public/registerEncryptor.wasm cmd/wasm/registerEncryptor.go

