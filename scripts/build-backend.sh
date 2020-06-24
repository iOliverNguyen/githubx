#!/usr/bin/env bash
set -eo pipefail
: "${GITHUBX_IMAGE?Must set GITHUBX_IMAGE}"

rm -rf build/backend || true
mkdir -p build/backend
cd backend
CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o ../build/backend/githubx ./cmd/githubx
cd ..
docker build -t "$GITHUBX_IMAGE" -f scripts/Dockerfile build/backend
