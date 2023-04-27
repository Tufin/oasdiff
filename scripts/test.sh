#!/bin/bash

set -euxo pipefail

go mod download && go mod tidy && go mod verify
go vet ./...
go fmt ./...

# update breaking-changes examples doc
./scripts/doc_breaking_changes.sh > BREAKING-CHANGES.md

# check if any changes need to be pushed
git diff --exit-code --name-only
