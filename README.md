[![CI](https://github.com/Tufin/oasdiff/workflows/go/badge.svg)](https://github.com/Tufin/oasdiff/actions)
[![codecov](https://codecov.io/gh/tufin/oasdiff/branch/main/graph/badge.svg?token=Y8BM6X77JY)](https://codecov.io/gh/tufin/oasdiff)
[![Go Report Card](https://goreportcard.com/badge/github.com/tufin/oasdiff)](https://goreportcard.com/report/github.com/tufin/oasdiff)
[![GoDoc](https://godoc.org/github.com/tufin/oasdiff?status.svg)](https://godoc.org/github.com/tufin/oasdiff)
[![Docker Image Version](https://img.shields.io/docker/v/tufin/oasdiff?sort=semver)](https://hub.docker.com/r/tufin/oasdiff/tags)

# OpenAPI Diff
This is a diff tool for [OpenAPI Spec 3](https://swagger.io/specification/).
It can be used in two ways:
1. Generate a diff report in [YAML](#Output-example---YAML) or [Text/Markdown](#Output-example---TEXT) from the cmd-line.  
2. Embed into a [Go](https://golang.org) program as a go module

## Unique features vs. other diff tools
- This is a [go module](https://blog.golang.org/using-go-modules) - it can be embedded into other Go programs
- Comprehensive diff: covers just about every aspect of OpenAPI Spec. See [work in progress](#Work-in-progress) for some limitations.
- Deep diff into paths, operations, parameters, request bodies, responses, schemas, enums, callbacks, security etc.
- Patch support is currently being added - see [work in progress](#Work-in-progress)

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

## Running with Docker

### Comparing public files (text output):

```
docker run --rm -t tufin/oasdiff -base https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml -revision https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml -format text
```

### Comparing public files (yaml output):

```
docker run --rm -t tufin/oasdiff -base https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml -revision https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```

### Comparing public files (summary):

```
docker run --rm -t tufin/oasdiff -base https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml -revision https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml -summary
```

### Comparing local files (yaml output):

```
docker run --rm -t -v $(pwd)/data:/data:ro tufin/oasdiff -base /data/openapi-test1.yaml -revision /data/openapi-test3.yaml
```

Replace "$(pwd)/data" by the path that contains your files.

### Help

```
docker run --rm -t tufin/oasdiff -help
```

## Output example - Text

### New Endpoints
-----------------

### Deleted Endpoints
---------------------
POST /subscribe  
POST /register  

### Modified Endpoints
----------------------
GET /api/{domain}/{project}/badges/security-score  
* Modified query param: filter
  - Content changed
    - Schema changed
* Modified query param: image
* Modified header param: user
  - Schema changed
    - Schema added
  - Content changed
* Modified cookie param: test
  - Content changed
* Response changed
  - New response: default
  - Deleted response: 200
  - Modified response: 201
    - Content changed
      - Schema changed
        - Type changed from 'string' to 'object'

GET /api/{domain}/{project}/install-command
* Deleted header param: network-policies
* Response changed
  - Modified response: default
    - Description changed from 'Tufin1' to 'Tufin'
    - Headers changed
      - Deleted header: X-RateLimit-Limit
        
## Output example - YAML

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
                      mediaType: true
                header:
                  user:
                    schema:
                      schemaAdded: true
                    content:
                      mediaTypeDeleted: true
                query:
                  filter:
                    content:
                      schema:
                        required:
                          added:
                            - type
            responses:
              added:
                - default
              deleted:
                - "200"
              modified:
                "201":
                  content:
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
      path: /subscribe
    - method: POST
      path: /register
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
                mediaType: true
          header:
            user:
              schema:
                schemaAdded: true
              content:
                mediaTypeDeleted: true
          query:
            filter:
              content:
                schema:
                  required:
                    added:
                      - type
      responses:
        added:
          - default
        deleted:
          - "200"
        modified:
          "201":
            content:
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
```

## Embedding into your Go program
```
diff.Get(&diff.Config{}, spec1, spec2)
```
See full example: [main.go](main.go)

## Notes
1. oasdiff expects [OpenAPI References](https://swagger.io/docs/specification/using-ref/) to be resolved.  
References are normally resolved automatically when you load the spec. In other cases you can resolve refs using [this function](https://pkg.go.dev/github.com/getkin/kin-openapi/openapi3#SwaggerLoader.ResolveRefsIn).

2. Use [configuration](diff/config.go) to exclude certain types of changes:
   - [Examples](https://swagger.io/specification/#example-object) 
   - Descriptions
  
3. [Extensions](https://swagger.io/specification/#specification-extensions) are excluded by default. Use [configuration](diff/config.go) to specify which ones to include.

4. Paths vs. Endpoints  
OpenAPI Specification has a hierarchial model of [Paths](https://swagger.io/specification/#paths-object) and [Operations](https://swagger.io/specification/#operation-object).  
oasdiff respects this heirarchy and displays a hierarchial diff with path changes: added, deleted and modified, and within the later "modified" section, another set of operation changes: added, deleted and modified.  
For example:
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
oasdiff also outputs an altrnate simplified view of "endpoints" which are a flattened combination of OpenAPI Paths and Operations, for example:
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

## Work in progress
1. While most aspects of OpenAPI Spec are already supported by this diff tool, one is still missing: server variables. 
2. Patch support: currently supports Descriptions and a few fields in Schema 

## Credits
This project relies on the excellent implementation of OpenAPI 3.0 for Go: [kin-openapi](https://github.com/getkin/kin-openapi) 

## Pull requests are welcome.
