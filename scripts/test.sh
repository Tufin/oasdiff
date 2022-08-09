#!/bin/bash

set -euxo pipefail

go mod download && go mod tidy && go mod verify
go vet ./...
go fmt ./...
go test -v -run TestRaceyPatternSchema -race ./...

# update breaking-changes doc
./scripts/doc_breaking_changes.sh > breaking-changes.md
