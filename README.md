[![codecov](https://codecov.io/gh/Tufin/oasdiff/branch/master/graph/badge.svg?token=Y8BM6X77JY)](https://codecov.io/gh/Tufin/oasdiff)
[![CircleCI](https://circleci.com/gh/Tufin/oasdiff.svg?style=svg)](https://circleci.com/gh/Tufin/oasdiff)

# OpenAPI Spec Diff
A diff tool for OpenAPI Spec 3.  

## Unique features vs. other OAS3 diff tools
- go module
- deep diff into paths, parameters, responses, schemas, enums etc.

## Build
```
git clone https://github.com/Tufin/oasdiff.git
cd oasdiff
go build
```

## Running from the command-line
```
./oasdiff -base data/openapi-test1.yaml -revision data/openapi-test2.yaml
```

## Help
```
./oasdiff --help
```

## Embedding into your Go program
See main.go