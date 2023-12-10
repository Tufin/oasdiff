#!/bin/bash

set -euxo pipefail

go mod download && go mod tidy && go mod verify
go vet ./...

# update breaking-changes examples doc
./scripts/doc_breaking_changes.sh > BREAKING-CHANGES-EXAMPLES.md

go fmt ./...

git add "BREAKING-CHANGES-EXAMPLES.md"