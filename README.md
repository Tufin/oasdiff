
[![CI](https://github.com/Tufin/oasdiff/workflows/go/badge.svg)](https://github.com/Tufin/oasdiff/actions)
[![codecov](https://codecov.io/gh/tufin/oasdiff/branch/main/graph/badge.svg?token=Y8BM6X77JY)](https://codecov.io/gh/tufin/oasdiff)
[![Go Report Card](https://goreportcard.com/badge/github.com/tufin/oasdiff)](https://goreportcard.com/report/github.com/tufin/oasdiff)
[![GoDoc](https://godoc.org/github.com/tufin/oasdiff?status.svg)](https://godoc.org/github.com/tufin/oasdiff)
[![Docker Image Version](https://img.shields.io/docker/v/tufin/oasdiff?sort=semver)](https://hub.docker.com/r/tufin/oasdiff/tags)
[![Slack](https://img.shields.io/badge/slack-&#64;oasdiff-green.svg?logo=slack)](https://join.slack.com/t/oasdiff/shared_invite/zt-1wvo7wois-ttncNBmyjyRXqBzyg~P6oA)

![image](https://github.com/yonatanmgr/oasdiff/assets/31913495/4fc67b8f-a7d6-4417-bebe-9ab2c7935a1e)

Command-line and Go package to compare and detect breaking changes in OpenAPI specs.


## Try it
```
docker run --rm -t tufin/oasdiff diff https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml -f text
```
## Features 
- Detect [breaking changes](BREAKING-CHANGES.md)
- OpenAPI diff in YAML, JSON, Text/Markdown or HTML
- [Run from Docker](#running-with-docker)
- [Embed in your go program](#embedding-oasdiff-into-your-program)
- OpenAPI diff of local files system or remote files over http/s
- Compare specs in YAML or JSON format
- [Compare two collections of specs](#composed-mode)
- Comprehensive diff including all aspects of [OpenAPI Specification](https://swagger.io/specification/): paths, operations, parameters, request bodies, responses, schemas, enums, callbacks, security etc.
- [API deprecation](API-DEPRECATION.md)
- [Path prefix modification](#path-prefix-modification)
- [Path parameter renaming](#path-parameter-reanaming)
- [Excluding certain kinds of changes](#excluding-specific-kinds-of-changes)
- [Excluding endpoints](#excluding-specific-endpoints)
- [Extending breaking-changes with custom checks](CUSTOMIZING-CHECKS.md)
- Display a user-friendly changelog of all important API changes
- Localization: display breaking-changes and changelog messages in English or Russian (please submit an issue if you want to add another language)


## Install with Go
```bash
go install github.com/tufin/oasdiff@latest
```

## Install on macOS with Brew
```bash
brew tap tufin/homebrew-tufin
brew install oasdiff
```

## Install on macOS, Windows and Linux
Copy binaries from [latest release](https://github.com/Tufin/oasdiff/releases/)

## Wrappers
- [GitHub Action](https://github.com/oasdiff/oasdiff-action)
- [Cloud Service](#openapi-diff-and-breaking-changes-as-a-service)
- [OpenAPI Sync: Get notified when an API provider breaks the API](https://github.com/oasdiff/sync/)

## Usage Examples

### OpenAPI diff of local files in YAML
```bash
oasdiff diff data/openapi-test1.yaml data/openapi-test2.yaml
```
The default diff output format is YAML.  
No output means that the diff is empty, or, in other words, there are no changes.

### OpenAPI diff of local files in Text/Markdown 
```bash
oasdiff diff data/openapi-test1.yaml data/openapi-test2.yaml -f text
```
The Text/Markdown diff report provides a simplified and partial view of the changes.  
To view all details, use the default format: YAML.  
If you'd like to see additional details in the text/markdown report, please submit a [feature request](https://github.com/Tufin/oasdiff/issues/new?assignees=&labels=&template=feature_request.md&title=).

### OpenAPI diff of local files in HTML
```bash
oasdiff diff data/openapi-test1.yaml data/openapi-test2.yaml -f html 
```
The HTML diff report provides a simplified and partial view of the changes.  
To view all details, use the default format: YAML.  
If you'd like to see additional details in the HTML report, please submit a [feature request](https://github.com/Tufin/oasdiff/issues/new?assignees=&labels=&template=feature_request.md&title=).


### OpenAPI diff for remote files over http/s
```bash
oasdiff diff https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml -f text
```

### Breaking changes
```bash
oasdiff breaking https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```

### Breaking changes as YAML
```bash
oasdiff breaking https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml -f yaml
```

### Breaking changes across multiple specs with globs
```bash
oasdiff breaking "data/composed/base/*.yaml" "data/composed/revision/*.yaml" -c
```

### Fail with exit code 1 if any ERR-level breaking changes are found
```bash
oasdiff breaking "data/composed/base/*.yaml" "data/composed/revision/*.yaml" -c -o ERR
```

### Fail with exit code 1 if any change is found
```bash
oasdiff diff https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml -f text -o
```

### OpenAPI changelog 
```bash
oasdiff changelog https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```

### OpenAPI diff for endpoints containing "/api" in the path
```bash
oasdiff diff https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml -f text -p "/api"
```
Filters are applied recursively at all levels. For example, if a path contains a [callback](https://swagger.io/docs/specification/callbacks/), the filter will be applied both to the path itself and to the callback path. To include such a nested change, use a regular expression that contains both paths, for example ```-filter "path|callback-path"```

### Exclude paths and operations with extension "x-beta"
```bash
oasdiff diff https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml -f text --filter-extension "x-beta"
``` 
Notes:
1. [OpenAPI Extensions](https://swagger.io/docs/specification/openapi-extensions/) can be defined both at the [path](https://swagger.io/docs/specification/paths-and-operations/) level and at the [operation](https://swagger.io/docs/specification/paths-and-operations/) level. Both are matched and excluded with this flag.
2. If a path or operation has a matching extension only in one of the specs, but not in the other, it will appear as Added or Deleted.

### Ignore changes to description and examples
```bash
oasdiff diff data/openapi-test1.yaml data/openapi-test3.yaml --exclude-elements description,examples -f text
``` 

### Display change summary
```bash
oasdiff summary https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```

### OpenAPI Diff with Docker
To run with docker just replace the `oasdiff` command by `docker run --rm -t tufin/oasdiff`, for example:

```bash
docker run --rm -t tufin/oasdiff diff https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml -f text
```

### Breaking changes with Docker
```bash
docker run --rm -t tufin/oasdiff breaking https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```

### Comparing local files with Docker
```bash
docker run --rm -t -v $(pwd)/data:/data:ro tufin/oasdiff diff /data/openapi-test1.yaml /data/openapi-test3.yaml
```

Replace `$(pwd)/data` by the path that contains your files.  
Note that the spec paths must begin with `/`.  

## OpenAPI Diff and Breaking-Changes as a Service
You can use oasdiff as a service like this:
```
curl -X POST -F base=@spec1.yaml -F revision=@spec2.yaml https://api.oasdiff.com/diff
```
Or, to see breaking changes:
```
curl -X POST -F base=@spec1.yaml -F revision=@spec2.yaml https://api.oasdiff.com/breaking-changes
```
Service source code: https://github.com/oasdiff/oasdiff-service

## Diff Output Formats
The default diff output format, YAML, provides a full view of all diff details.  
Note that no output in the YAML format signifies that the diff is empty, or, in other words, there are no changes.  
Other formats: text, markdown, and HTML, are designed to be more user-friendly by providing only the most important parts of the diff, in a simplified format.  
The JSON format works only with `-exclude-elements endpoints` and is intended as a workaround for YAML complex mapping keys which aren't supported by some libraries (see comment at end of next section for more details).
If you wish to include additional details in non-YAML formats, please open an issue.

## Paths vs. Endpoints
OpenAPI Specification has a hierarchical model of [Paths](https://swagger.io/specification/#paths-object) and [Operations](https://swagger.io/specification/#operation-object) (HTTP methods).  
oasdiff respects this hierarchy and displays a hierarchical diff with path changes: added, deleted and modified, and within the latter, "modified" section, another set of operation changes: added, deleted and modified. For example:
```yaml
paths:
    deleted:
        - /register
        - /subscribe
    modified:
        /api/{domain}/{project}/badges/security-score:
            operations:
                modified:
                    GET:
```
oasdiff also outputs an alternate simplified diff per "endpoint" which is a combination of Path + Operation, for example:
```yaml
endpoints:
    deleted:
        - method: POST
          path: /subscribe
        - method: POST
          path: /register
    modified:
        ?   method: GET
            path: /api/{domain}/{project}/badges/security-score
        :   tags:
                deleted:
                    - security
```
The modified endpoints section has two items per key, method and path, this is called a [complex mapping key](https://stackoverflow.com/questions/33987316/what-is-a-complex-mapping-key-in-yaml) in YAML.  
Some YAML libraries don't support complex mapping keys:
- python PyYAML: see https://github.com/Tufin/oasdiff/issues/94#issuecomment-1087468450
- golang gopkg.in/yaml.v3 fails to unmarshal the oasdiff output. This package offers a solution: https://github.com/tliron/yamlkeys

In such cases, consider using `-exclude-elements endpoints` and `-format json` as a workaround.

## Composed Mode
Composed mode compares two collections of OpenAPI specs instead of a pair of specs in the default mode.
The collections are specified using a [glob](https://en.wikipedia.org/wiki/Glob_(programming)).
This can be useful when your APIs are defined across multiple files, for example, when multiple services, each one with its own spec, are exposed behind an API gateway, and you want to check changes across all the specs at once.

Notes: 
1. Composed mode compares only [paths and endpoints](#paths-vs-endpoints), other resources are compared only if referenced from the paths or endpoints.
2. Composed mode doesn't support [Path Prefix Modification](#path-prefix-modification) 
3. Learn more about how oasdiff [matches endpoints to each other](MATCHING-ENDPOINTS.md)

## Path Prefix Modification
Sometimes paths prefixes need to be modified, for example, to create a new version:
- /api/v1/...
- /api/v2/...

oasdiff allows comparison of API specifications with modified prefixes by stripping and/or prepending path prefixes.  
In the example above you could compare the files as follows:
```
oasdiff diff original.yaml new.yaml --strip-prefix-base /api/v1 --prefix-base /api/v2
```
or
```
oasdiff diff original.yaml new.yaml --strip-prefix-base /api/v1 --strip-prefix-revision /api/v2
```
Note that stripping precedes prepending.

## Path Parameter Reanaming
Sometimes developers decide to change names of path parameters, for example, in order to follow a certain naming convention.  
See [this](MATCHING-ENDPOINTS.md) to learn more about how oasdiff supports path parameter renaming.

## Excluding Specific Kinds of Changes 
You can use the `--exclude-elements` flag to exclude certain kinds of changes:
- Use `--exclude-elements examples` to exclude [Examples](https://swagger.io/specification/#example-object)
- Use `--exclude-elements description` to exclude description fields
- Use `--exclude-elements title` to exclude title fields
- Use `--exclude-elements summary` to exclude summary fields
- Use `--exclude-elements endpoints` to exclude the [endpoints diff](#paths-vs-endpoints)

You can ignore multiple elements with a comma-separated list of excluded elements as in [this example](#ignore-changes-to-description-and-examples).  
Note that [Extensions](https://swagger.io/specification/#specification-extensions) are always excluded from the diff.

## Excluding Specific Endpoints
You can filter endpoints in two ways:
1. By path name: use the `--match-path` option to exclude paths that don't match the given regular expression, see [example](#openapi-diff-for-endpoints-containing-api-in-the-path)
2. By extension: use the `--filter-extension` option to exclude paths and operations with an OpenAPI Extension matching the given regular expression, see [example](#exclude-paths-and-operations-with-extension-x-beta)

## Notes for Go Developers
### Embedding oasdiff into your program
```go
diff.Get(&diff.Config{}, spec1, spec2)
```

### Code Examples
- [diff](https://pkg.go.dev/github.com/tufin/oasdiff/diff#example-Get)
- [breaking changes](https://pkg.go.dev/github.com/tufin/oasdiff/diff#example-GetPathsDiff)


### OpenAPI References
oasdiff expects [OpenAPI References](https://swagger.io/docs/specification/using-ref/) to be resolved.  
References are normally resolved automatically when you load the spec. In other cases you can resolve refs using [Loader.ResolveRefsIn](https://pkg.go.dev/github.com/getkin/kin-openapi/openapi3#Loader.ResolveRefsIn).

## Requests for enhancements
1. OpenAPI 3.1 support: see https://github.com/Tufin/oasdiff/issues/52

If you have other ideas, please [let us know](https://github.com/Tufin/oasdiff/discussions/new?category=ideas).

## Credits
This project relies on the excellent implementation of OpenAPI 3.0 for Go: [kin-openapi](https://github.com/getkin/kin-openapi).
