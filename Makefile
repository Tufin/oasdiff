# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

VERSION=$(shell git describe --always --tags | cut -d "v" -f 2)
LINKER_FLAGS=-s -w -X github.com/tufin/oasdiff/build.Version=${VERSION}
GOLANGCI_VERSION=v1.50.1

.PHONY: test
test:
	scripts/test.sh

.PHONY: build
build:
	@echo "==> Building oasdiff binary"
	go build -ldflags "$(LINKER_FLAGS)" -o ./bin/oasdiff $(MCLI_SOURCE_FILES)

.PHONY: deps
deps:  ## Download go module dependencies
	@echo "==> Installing go.mod dependencies..."
	go mod download
	go mod tidy


.PHONY: lint
lint: ## Run linter
	golangci-lint run

.PHONY: devtools
devtools:  ## Install dev tools
	@echo "==> Installing dev tools..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_VERSION)
