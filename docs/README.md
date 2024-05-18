
[![CI](https://github.com/Tufin/oasdiff/workflows/go/badge.svg)](https://github.com/Tufin/oasdiff/actions)
[![codecov](https://codecov.io/gh/tufin/oasdiff/branch/main/graph/badge.svg?token=Y8BM6X77JY)](https://codecov.io/gh/tufin/oasdiff)
[![Go Report Card](https://goreportcard.com/badge/github.com/tufin/oasdiff)](https://goreportcard.com/report/github.com/tufin/oasdiff)
[![GoDoc](https://godoc.org/github.com/tufin/oasdiff?status.svg)](https://godoc.org/github.com/tufin/oasdiff)
[![Docker Image Version](https://img.shields.io/docker/v/tufin/oasdiff?sort=semver)](https://hub.docker.com/r/tufin/oasdiff/tags)
[![Slack](https://img.shields.io/badge/slack-&#64;oasdiff-green.svg?logo=slack)](https://join.slack.com/t/oasdiff/shared_invite/zt-1wvo7wois-ttncNBmyjyRXqBzyg~P6oA)

![oasdiff banner](https://github.com/yonatanmgr/oasdiff/assets/31913495/ac9b154e-72d1-4969-bc3b-f527bbe7751d)


Command-line and Go package to compare and detect breaking changes in OpenAPI specs.


## Try it
```
docker run --rm -t tufin/oasdiff changelog https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test5.yaml
```

## Features 
- Detect [breaking changes](BREAKING-CHANGES.md)
- Display a user-friendly [changelog](BREAKING-CHANGES.md) of all important API changes
- OpenAPI [diff](DIFF.md) in YAML, JSON, Text, Markdown or HTML
- [Run from Docker](DOCKER.md)
- [GitHub Action](https://github.com/oasdiff/oasdiff-action)
- [Cloud Service](OASDIFF-SERVICE.md)
- [OpenAPI Sync: Get notified when an API provider breaks the API](https://github.com/oasdiff/sync/)
- [Embed in your go program](GO.md)
- Compare local files or remote files over http/s
- Compare specs in YAML or JSON format
- [Compare two collections of specs](COMPOSED.md)
- Comprehensive diff including all aspects of [OpenAPI Specification](https://swagger.io/specification/): paths, operations, parameters, request bodies, responses, schemas, enums, callbacks, security etc.
- [API deprecation](API-DEPRECATION.md)
- [Multiple versions of the same endpoint](MATCHING-ENDPOINTS.md#duplicate-endpoints)
- [Merge allOf schemas](ALLOF.md)
- [Merge common parameters](COMMON-PARAMS.md)
- [Case-insensitive header comparison](HEADER-DIFF.md)
- [Path prefix modification](PATH-PREFIX.md)
- [Path parameter renaming](PATH-PARAM-RENAME.md)
- [Excluding certain kinds of changes](DIFF.md#excluding-specific-kinds-of-changes)
- [Tracking changes to OpenAPI Extensions](DIFF.md#openapi-extensions)
- [Filtering endpoints](FILTERING-ENDPOINTS.md)
- [Extending breaking changes with custom checks](CUSTOMIZING-CHECKS.md)
- Localization: display breaking changes and changelog messages in English or Russian ([please contribute support for your language](https://github.com/Tufin/oasdiff/issues/383))


## Demo
<img src="./demo.svg">

## Installation

### Install with Go
```bash
go install github.com/tufin/oasdiff@latest
```

### Install on macOS with Brew
```bash
brew tap tufin/homebrew-tufin
brew install oasdiff
```

### Install on macOS, Windows and Linux
Copy binaries from [latest release](https://github.com/Tufin/oasdiff/releases/)

## The main commands
- [diff](DIFF.md): the diff between OpenAPI specs, fully detailed
- [breaking](BREAKING-CHANGES.md): breaking changes between OpenAPI specs  
- [changelog](BREAKING-CHANGES.md): important changes between OpenAPI specs including breaking and non-breaking changes
- [flatten](ALLOF.md): replace all instances of allOf by a merged equivalent
- checks: displays the different checks that oasdiff runs to detect changes

## Credits
This project relies on the excellent implementation of OpenAPI 3.0 for Go: [kin-openapi](https://github.com/getkin/kin-openapi).

## Feedback
We welcome your feedback.  
If you have ideas for improvement or additional needs around APIs, please [let us know](https://github.com/Tufin/oasdiff/discussions/new?category=ideas).
