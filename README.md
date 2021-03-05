[![codecov](https://codecov.io/gh/Tufin/oasdiff/branch/master/graph/badge.svg?token=Y8BM6X77JY)](https://codecov.io/gh/Tufin/oasdiff)
[![CircleCI](https://circleci.com/gh/Tufin/oasdiff.svg?style=svg)](https://circleci.com/gh/Tufin/oasdiff)
[![Go Report Card](https://goreportcard.com/badge/github.com/Tufin/oasdiff)](https://goreportcard.com/report/github.com/Tufin/oasdiff)
[![GoDoc](https://godoc.org/github.com/Tufin/oasdiff?status.svg)](https://godoc.org/github.com/Tufin/oasdiff)

# OpenAPI Spec Diff
A diff tool for OpenAPI Spec 3 written in [Go](https://golang.org).

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
```
json.MarshalIndent(diff.Get(&diff.Config{}, spec1, spec2), "", " ")
```
See full example: [main.go](main.go)

## Notes
1. oasdiff expects [OpenAPI References](https://swagger.io/docs/specification/using-ref/) to be resolved.  
You can resolve refs using [this function](https://pkg.go.dev/github.com/getkin/kin-openapi/openapi3#SwaggerLoader.ResolveRefsIn) from the openapi3 package.

2. oasdiff ignores changes to [Examples](https://swagger.io/specification/#example-object) and [Extensions](https://swagger.io/specification/#specification-extensions) by default. You can change this behavior through [configuration](diff/config.go).

## Documentation
https://pkg.go.dev/github.com/tufin/oasdiff/diff

