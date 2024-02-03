#!/usr/bin/env bash
set -e
set -o xtrace

cd "$(dirname $0)"

go build -o build/send-me-a-secret cmd/send-me-a-secret/*.go

GOOS=js GOARCH=wasm go build -o web/browser-me-a-secret/public/registerEncryptor.wasm cmd/wasm/registerEncryptor.go
cp "$(go env GOROOT)"/misc/wasm/wasm_exec.js web/browser-me-a-secret/public/wasm_exec.js
