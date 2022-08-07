#!/bin/bash

set -euxo pipefail

go mod download && go mod tidy && go mod verify
go vet ./...
go fmt ./...
go test -v -run TestRaceyPatternSchema -race ./...

# check if breaking-changes example's doc is up to date,
# by creating a new breaking-change's doc and run a diff
./scripts/doc_breaking_changes.sh > breaking-changes.md
git --no-pager diff --exit-code -- breaking_changes.md
