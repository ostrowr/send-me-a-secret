#!/usr/bin/env bash
set -e

go build -o build/private/decrypt cmd/decrypt/decrypt.go
go build -o build/public/encrypt cmd/encrypt/encrypt.go

GOOS=js GOARCH=wasm go build -o web/browser-me-a-secret/public/registerEncryptor.wasm cmd/wasm/registerEncryptor.go

