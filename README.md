[![CI](https://github.com/Tufin/oasdiff/workflows/go/badge.svg)](https://github.com/Tufin/oasdiff/actions)
[![codecov](https://codecov.io/gh/tufin/oasdiff/branch/main/graph/badge.svg?token=Y8BM6X77JY)](https://codecov.io/gh/tufin/oasdiff)
[![Go Report Card](https://goreportcard.com/badge/github.com/tufin/oasdiff)](https://goreportcard.com/report/github.com/tufin/oasdiff)
[![GoDoc](https://godoc.org/github.com/tufin/oasdiff?status.svg)](https://godoc.org/github.com/tufin/oasdiff)
[![Docker Image Version](https://img.shields.io/docker/v/tufin/oasdiff?sort=semver)](https://hub.docker.com/r/tufin/oasdiff/tags)

# OpenAPI Diff
A diff tool for [OpenAPI Spec 3](https://swagger.io/specification/).

## Features 
- Generate a diff report in YAML, Text/Markdown or HTML
- [Run from Docker](#running-with-docker)
- [Embed in your go program](#embedding-into-your-go-program)
- Compare specs from the file system or over http/s
- Compare specs in YAML or JSON format
- Comprehensive diff including all aspects of [OpenAPI Specification](https://swagger.io/specification/): paths, operations, parameters, request bodies, responses, schemas, enums, callbacks, security etc.
- Detect breaking changes (Beta feature. Please report issues)

## Build
```
git clone https://github.com/Tufin/oasdiff.git
cd oasdiff
go build
```

## Running from the command-line
```bash
./oasdiff -base data/openapi-test1.yaml -revision data/openapi-test2.yaml
```

## Running with Docker

### Text/Markdown output (http):

```bash
docker run --rm -t tufin/oasdiff -format text -base https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml -revision https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```
Text and markdown output provide only a subset of the full diff.
To see the full diff, use the default format: YAML.

### HTML output (http):

```bash
docker run --rm -t tufin/oasdiff -format html -base https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml -revision https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```
HTML output provides only a subset of the full diff.
To see the full diff, use the default format: YAML.

### YAML output (http):

```bash
docker run --rm -t tufin/oasdiff -base https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml -revision https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```
This is the default output format.  
Note that no output means an empty diff (no changes).

### Breaking changes (http):

```bash
docker run --rm -t tufin/oasdiff -breaking-only -base https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml -revision https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```

### YAML output (local files):

```bash
docker run --rm -t -v $(pwd)/data:/data:ro tufin/oasdiff -base /data/openapi-test1.yaml -revision /data/openapi-test3.yaml
```

Replace `$(pwd)/data` by the path that contains your files.  
Add the `-format` flag to generate other formats (text or html).

## Usage

```./oasdiff -help```
```
-base string
    path of original OpenAPI spec in YAML or JSON format
-breaking-only
    display breaking changes only
-exclude-description
    exclude changes to descriptions
-exclude-examples
    exclude changes to examples
-fail-on-diff
    fail with exit code 1 if a difference is found
-filter string
    regex to filter result paths (optional)
-format string
    output format: yaml, text or html (default "yaml")
-prefix string
    path prefix that exists in base spec but not the revision (optional)
-revision string
    path of revised OpenAPI spec in YAML or JSON format
-summary
    display a summary of the changes instead of the full diff
```
All arguments can be passed with one or two leading minus signs.  
For example ```-breaking-only``` and ```--breaking-only``` are equivalent.

## Output example - Text/Markdown
```./oasdiff -format text -base data/openapi-test1.yaml -revision data/openapi-test5.yaml```

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
- Modified query param: image
- Modified query param: token
  - Schema changed
    - MaxLength changed from 29 to <nil>
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

## Output example - YAML
```./oasdiff -base data/openapi-test1.yaml -revision data/openapi-test5.yaml```

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

## Embedding into your Go program
```go
diff.Get(&diff.Config{}, spec1, spec2)
```
See full example: [main.go](main.go)

## Notes
1. Output Formats  
   - The default output format, YAML, provides a full view of all diff details.  
     Note that no output in YAML format signifies that the diff is empty, or, in other words, there are no changes.
   - Other formats: text, markdown and HTML, are designed to be more user-friendly by providing only the most important parts of the diff, in a simplified format.  
     If you wish to include additional details in non-YAML formats, please open an issue.

2. Paths vs. Endpoints  
OpenAPI Specification has a hierarchial model of [Paths](https://swagger.io/specification/#paths-object) and [Operations](https://swagger.io/specification/#operation-object) (HTTP methods).  
oasdiff respects this heirarchy and displays a hierarchial diff with path changes: added, deleted and modified, and within the latter, "modified" section, another set of operation changes: added, deleted and modified. For example:
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

## Notes for Developers
1. oasdiff expects [OpenAPI References](https://swagger.io/docs/specification/using-ref/) to be resolved.  
References are normally resolved automatically when you load the spec. In other cases you can resolve refs using [Loader.ResolveRefsIn](https://pkg.go.dev/github.com/getkin/kin-openapi/openapi3#Loader.ResolveRefsIn).

2. Use [configuration](diff/config.go) to exclude certain types of changes:
   - [Examples](https://swagger.io/specification/#example-object) 
   - Descriptions
   - [Extensions](https://swagger.io/specification/#specification-extensions) are excluded by default

## Work in progress
1. Patch support: currently supports Descriptions and a few fields in Schema 

## Requests for enhancements
1. OpenAPI 3.1 support: see https://github.com/Tufin/oasdiff/issues/52

## Credits
This project relies on the excellent implementation of OpenAPI 3.0 for Go: [kin-openapi](https://github.com/getkin/kin-openapi) 
