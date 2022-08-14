[![CI](https://github.com/Tufin/oasdiff/workflows/go/badge.svg)](https://github.com/Tufin/oasdiff/actions)
[![codecov](https://codecov.io/gh/tufin/oasdiff/branch/main/graph/badge.svg?token=Y8BM6X77JY)](https://codecov.io/gh/tufin/oasdiff)
[![Go Report Card](https://goreportcard.com/badge/github.com/tufin/oasdiff)](https://goreportcard.com/report/github.com/tufin/oasdiff)
[![GoDoc](https://godoc.org/github.com/tufin/oasdiff?status.svg)](https://godoc.org/github.com/tufin/oasdiff)
[![Docker Image Version](https://img.shields.io/docker/v/tufin/oasdiff?sort=semver)](https://hub.docker.com/r/tufin/oasdiff/tags)

# OpenAPI Diff
A diff tool for [OpenAPI Spec 3](https://swagger.io/specification/).

## Try it
```
docker run --rm -t tufin/oasdiff -format text -base https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml -revision https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```

## Features 
- Generate a diff report in YAML, Text/Markdown or HTML
- [Run from Docker](#running-with-docker)
- [Embed in your go program](#embedding-into-your-go-program)
- Compare specs from the file system or over http/s
- Compare specs in YAML or JSON format
- Comprehensive diff including all aspects of [OpenAPI Specification](https://swagger.io/specification/): paths, operations, parameters, request bodies, responses, schemas, enums, callbacks, security etc.
- Detect [breaking changes](#breaking-changes)
- Allow [non-breaking removal of deprecated resources](#non-breaking-removal-of-deprecated-resources)
- Support [path prefix modification](#path-prefix-modification)
- [GitHub Action](https://github.com/marketplace/actions/openapi-spec-diff)

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

## Usage

```oasdiff -help```
```
Usage of oasdiff:
  -base string
    	path of original OpenAPI spec in YAML or JSON format
  -breaking-only
    	display breaking changes only
  -deprecation-days int
    	minimal number of days required between deprecating a resource and removing it without being considered 'breaking'
  -exclude-description
    	ignore changes to descriptions
  -exclude-examples
    	ignore changes to examples
  -fail-on-diff
    	fail with exit code 1 if a difference is found
  -filter string
    	if provided, diff will include only paths that match this regular expression
  -filter-extension string
    	if provided, diff will exclude paths and operations with an OpenAPI Extension matching this regular expression
  -format string
    	output format: yaml, text or html (default "yaml")
  -prefix string
    	deprecated. use -prefix-revision instead
  -prefix-base string
    	if provided, paths in original (base) spec will be prefixed with the given prefix before comparison
  -prefix-revision string
    	if provided, paths in revised (revision) spec will be prefixed with the given prefix before comparison
  -revision string
    	path of revised OpenAPI spec in YAML or JSON format
  -strip-prefix-base string
    	if provided, this prefix will be stripped from paths in original (base) spec before comparison
  -strip-prefix-revision string
    	if provided, this prefix will be stripped from paths in revised (revision) spec before comparison
  -summary
    	display a summary of the changes instead of the full diff
  -version
    	show version and quit
```
All arguments can be passed with one or two leading minus signs.  
For example ```-breaking-only``` and ```--breaking-only``` are equivalent.

## Usage Examples

### YAML diff of local files
```bash
oasdiff -base data/openapi-test1.yaml -revision data/openapi-test2.yaml
```
The default output format is YAML.  
No output means that the diff is empty, or, in other words, there are no changes.

### Text/Markdown diff of local files
```bash
oasdiff -format text -base data/openapi-test1.yaml -revision data/openapi-test2.yaml
```
The Text/Markdown diff report provides a simplified and partial view of the changes.  
To view all details, use the default format: YAML.  
If you'd like to see additional details in the text/markdown report, please submit a [feature request](https://github.com/Tufin/oasdiff/issues/new?assignees=&labels=&template=feature_request.md&title=).

### HTML diff of local files
```bash
oasdiff -format text -base data/openapi-test1.yaml -revision data/openapi-test2.yaml
```
The HTML diff report provides a simplified and partial view of the changes.  
To view all details, use the default format: YAML.  
If you'd like to see additional details in the HTML report, please submit a [feature request](https://github.com/Tufin/oasdiff/issues/new?assignees=&labels=&template=feature_request.md&title=).


### Diff files over http/s
```bash
oasdiff -format text -base https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml -revision https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```

### Display breaking changes only
```bash
oasdiff -breaking-only -format text -base https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml -revision https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```
See [breaking changes](#breaking-changes)

### Fail with exit code 1 if a change is found
```bash
oasdiff -fail-on-diff -format text -base https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml -revision https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```

### Fail with exit code 1 if a breaking change is found
```bash
oasdiff -fail-on-diff -breaking-only -format text -base https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml -revision https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```

### Display changes to endpoints containing "/api" in the path
```bash
oasdiff -format text -filter "/api" -base https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml -revision https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```
Filters are applied recursively at all levels. For example, if a path contains a [callback](https://swagger.io/docs/specification/callbacks/), the filter will be applied both to the path itself and to the callback path. To include such a nested change, use a regular expression that contains both paths, for example ```-filter "path|callback-path"```

### Exclude paths and operations with extension "x-beta"
```bash
oasdiff -format text -filter-extension "x-beta" -base https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml -revision https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
``` 
Notes:
1. [OpenAPI Extensions](https://swagger.io/docs/specification/openapi-extensions/) can be defined both at the [path](https://swagger.io/docs/specification/paths-and-operations/) level and at the [operation](https://swagger.io/docs/specification/paths-and-operations/) level. Both are matched and excluded with this flag.
2. If a path or operation has a matching extension only in one of the specs, but not in the other, it will appear as Added or Deleted.

### Ignore changes to descriptions and examples
```bash
oasdiff -exclude-description -exclude-examples -format text -base data/openapi-test1.yaml -revision data/openapi-test3.yaml
``` 

### Display change summary
```bash
oasdiff -summary -base https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml -revision https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```

### Running with Docker
To run with docker just replace the `oasdiff` command by `docker run --rm -t tufin/oasdiff`, for example:

```bash
docker run --rm -t tufin/oasdiff -format text -base https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml -revision https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```

### Comparing local files with Docker
```bash
docker run --rm -t -v $(pwd)/data:/data:ro tufin/oasdiff -base /data/openapi-test1.yaml -revision /data/openapi-test3.yaml
```

Replace `$(pwd)/data` by the path that contains your files.  
Note that the `-base` and `-revision` paths must begin with `/`.  

## Output example - Text/Markdown
```
oasdiff -format text -base data/openapi-test1.yaml -revision data/openapi-test5.yaml
```

### New Endpoints: None
-----------------------

### Deleted Endpoints: 2
------------------------
POST /register
POST /subscribe

### Modified Endpoints: 2
-------------------------
GET /api/{domain}/{project}/badges/security-score
- Modified query param: filter
  - Content changed
    - Modified media type: application/json
      - Schema changed
        - Required changed
          - New required property: type
- Modified query param: image
  - Examples changed
    - Deleted example: 0
- Modified query param: token
  - Schema changed
    - MaxLength changed from 29 to null
- Modified header param: user
  - Schema changed
    - Schema added
  - Content changed
    - Deleted media type: application/json
- Modified cookie param: test
  - Content changed
    - Modified media type: application/json
      - Schema changed
        - Type changed from 'object' to 'string'
- Responses changed
  - New response: default
  - Deleted response: 200
  - Modified response: 201
    - Content changed
      - Modified media type: application/xml
        - Schema changed
          - Type changed from 'string' to 'object'

GET /api/{domain}/{project}/install-command
- Deleted header param: network-policies
- Responses changed
  - Modified response: default
    - Description changed from 'Tufin1' to 'Tufin'
    - Headers changed
      - Deleted header: X-RateLimit-Limit
- Servers changed
  - New server: https://www.tufin.io/securecloud

Security Requirements changed
- Deleted security requirements: bearerAuth

Servers changed
- Deleted server: tufin.com

## Output example - YAML
```
oasdiff -base data/openapi-test1.yaml -revision data/openapi-test5.yaml
```

```yaml
info:
  title:
    from: Tufin
    to: Tufin1
  contact:
    added: true
  license:
    added: true
  version:
    from: 1.0.0
    to: 1.0.1
paths:
  deleted:
    - /register
    - /subscribe
  modified:
    /api/{domain}/{project}/badges/security-score:
      operations:
        modified:
          GET:
            tags:
              deleted:
                - security
            parameters:
              modified:
                cookie:
                  test:
                    content:
                      mediaTypeModified:
                        application/json:
                          schema:
                            type:
                              from: object
                              to: string
                header:
                  user:
                    schema:
                      schemaAdded: true
                    content:
                      mediaTypeDeleted:
                        - application/json
                query:
                  filter:
                    content:
                      mediaTypeModified:
                        application/json:
                          schema:
                            required:
                              stringsdiff:
                                added:
                                  - type
                  image:
                    examples:
                      deleted:
                        - "0"
                  token:
                    schema:
                      maxLength:
                        from: 29
                        to: null
            responses:
              added:
                - default
              deleted:
                - "200"
              modified:
                "201":
                  content:
                    mediaTypeModified:
                      application/xml:
                        schema:
                          type:
                            from: string
                            to: object
      parameters:
        deleted:
          path:
            - domain
    /api/{domain}/{project}/install-command:
      operations:
        modified:
          GET:
            parameters:
              deleted:
                header:
                  - network-policies
            responses:
              modified:
                default:
                  description:
                    from: Tufin1
                    to: Tufin
                  headers:
                    deleted:
                      - X-RateLimit-Limit
            servers:
              added:
                - https://www.tufin.io/securecloud
endpoints:
  deleted:
    - method: POST
      path: /register
    - method: POST
      path: /subscribe
  modified:
    ? method: GET
      path: /api/{domain}/{project}/install-command
    : parameters:
        deleted:
          header:
            - network-policies
      responses:
        modified:
          default:
            description:
              from: Tufin1
              to: Tufin
            headers:
              deleted:
                - X-RateLimit-Limit
      servers:
        added:
          - https://www.tufin.io/securecloud
    ? method: GET
      path: /api/{domain}/{project}/badges/security-score
    : tags:
        deleted:
          - security
      parameters:
        modified:
          cookie:
            test:
              content:
                mediaTypeModified:
                  application/json:
                    schema:
                      type:
                        from: object
                        to: string
          header:
            user:
              schema:
                schemaAdded: true
              content:
                mediaTypeDeleted:
                  - application/json
          query:
            filter:
              content:
                mediaTypeModified:
                  application/json:
                    schema:
                      required:
                        stringsdiff:
                          added:
                            - type
            image:
              examples:
                deleted:
                  - "0"
            token:
              schema:
                maxLength:
                  from: 29
                  to: null
      responses:
        added:
          - default
        deleted:
          - "200"
        modified:
          "201":
            content:
              mediaTypeModified:
                application/xml:
                  schema:
                    type:
                      from: string
                      to: object
security:
  deleted:
    - bearerAuth
servers:
  deleted:
    - tufin.com
tags:
  deleted:
    - security
    - reuven
externalDocs:
  deleted: true
components:
  schemas:
    added:
      - requests
    modified:
      network-policies:
        additionalPropertiesAllowed:
          from: true
          to: false
      rules:
        additionalPropertiesAllowed:
          from: null
          to: false
  parameters:
    deleted:
      header:
        - network-policies
  headers:
    deleted:
      - new
    modified:
      test:
        schema:
          additionalPropertiesAllowed:
            from: true
            to: false
      testc:
        content:
          mediaTypeModified:
            application/json:
              schema:
                type:
                  from: object
                  to: string
  requestBodies:
    deleted:
      - reuven
  responses:
    added:
      - default
    deleted:
      - OK
  securitySchemes:
    deleted:
      - OAuth
    modified:
      AccessToken:
        type:
          from: http
          to: oauth2
        scheme:
          from: bearer
          to: ""
        OAuthFlows:
          added: true
```

## Notes
### Output Formats
The default output format, YAML, provides a full view of all diff details.  
Note that no output in YAML format signifies that the diff is empty, or, in other words, there are no changes.  
Other formats: text, markdown and HTML, are designed to be more user-friendly by providing only the most important parts of the diff, in a simplified format.  
If you wish to include additional details in non-YAML formats, please open an issue.

### Paths vs. Endpoints
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
Some YAML libraries don't support complex mapping keys, for example, python's PyYAML. [Here's possible solution](https://github.com/Tufin/oasdiff/issues/94#issuecomment-1087468450).

### Breaking Changes
Breaking changes are changes that could break a client that is relying on the OpenAPI specification.  
[See some examples of breaking and non-breaking changes](breaking-changes.md).  
Note: this is a Beta feature. Please report issues.

### Non Breaking Removal of Deprecated Resources
OpenAPI specification supports a ```deprecated``` flag which can be used to mark operations and other object types as deprecated.  
Normally, deprecation **is not** considered a breaking change since it doesn't break the client but only serves as an indication of an intent to remove something in the future, in contrast, the eventual removal of a resource **is** considered a breaking change.

oasdiff allows you to gracefully remove a resource without getting the ```breaking-change``` warning, as follows:
1.   First, the resource is marked as ```deprecated``` and a [special extension](https://swagger.io/specification/#specification-extensions) ```x-sunset``` is added to announce the date at which the resource will be removed
   ```
   /api/test:
    get:
     deprecated: true
     x-sunset: "2022-08-10"
   ```
2. At the sunset date or anytime later, the resource can be removed without triggering a ```breaking-change``` warning. An earlier removal will be considered a breaking change.

In addition, oasdiff also allows you to control the minimal number of days required between deprecating a resource and removing it with the ```deprecation-days``` flag.  
For example, the following command requires any deprecation to be accompanied by an ```x-sunset``` extension with a date which is at least 30 days away, otherwise the deprecation itself will be considered a breaking change:
```
oasdiff -deprecation-days=10 -breaking-only -base data/deprecation/deprecated-base.yaml -revision data/deprecation/deprecated-past.yaml
```

Setting deprecation-days to 0 is equivalent to the default which allows non-breaking deprecation regardless of the sunset date.


### Path Prefix Modification
Sometimes paths prefixes need to be modified, for example, to create a new version:
- /api/v1/...
- /api/v2/...

oasdiff allows comparison of API specifications with modified prefixes by stripping and/or prepending path prefixes.  
In the example above you could compare the files as follows:
```
oasdiff -base original.yaml -revision new.yaml -strip-prefix-base /api/v1 -prefix-base /api/v2
```
or
```
oasdiff -base original.yaml -revision new.yaml -strip-prefix-base /api/v1 -strip-prefix-revision /api/v2
```
Note that stripping precedes prepending.

## Notes for Go Developers
### Embedding oasdiff into your program
```go
diff.Get(&diff.Config{}, spec1, spec2)
```
See full example: [main.go](main.go)

### OpenAPI References
oasdiff expects [OpenAPI References](https://swagger.io/docs/specification/using-ref/) to be resolved.  
References are normally resolved automatically when you load the spec. In other cases you can resolve refs using [Loader.ResolveRefsIn](https://pkg.go.dev/github.com/getkin/kin-openapi/openapi3#Loader.ResolveRefsIn).

### Excluding Changes to Examples etc.
Use [configuration](diff/config.go) to exclude certain types of changes:
- [Examples](https://swagger.io/specification/#example-object) 
- Descriptions
- [Extensions](https://swagger.io/specification/#specification-extensions) are excluded by default

## Requests for enhancements
1. OpenAPI 3.1 support: see https://github.com/Tufin/oasdiff/issues/52

If you have other ideas, please [let us know](https://github.com/Tufin/oasdiff/discussions/new?category=ideas).

## Credits
This project relies on the excellent implementation of OpenAPI 3.0 for Go: [kin-openapi](https://github.com/getkin/kin-openapi) 
