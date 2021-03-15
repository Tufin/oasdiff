[![codecov](https://codecov.io/gh/tufin/oasdiff/branch/master/graph/badge.svg?token=Y8BM6X77JY)](https://codecov.io/gh/tufin/oasdiff)
[![CircleCI](https://circleci.com/gh/Tufin/oasdiff.svg?style=svg)](https://circleci.com/gh/Tufin/oasdiff)
[![Go Report Card](https://goreportcard.com/badge/github.com/tufin/oasdiff)](https://goreportcard.com/report/github.com/tufin/oasdiff)
[![GoDoc](https://godoc.org/github.com/tufin/oasdiff?status.svg)](https://godoc.org/github.com/tufin/oasdiff)

# OpenAPI Diff Go Module
This [Go](https://golang.org) module provides a diff utility for [OpenAPI Spec 3](https://swagger.io/specification/).

## Unique features vs. other diff tools
- This is a [go module](https://blog.golang.org/using-go-modules) - it can be embedded into other Go programs
- Comprehensive diff: covers just about every aspect of OpenAPI Spec see [work in progress](#Notes) for some limitations
- Deep diff into paths, operations, parameters, request bodies, responses, schemas, enums, callbacks, security etc.

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

## Output example

```yaml
spec:
  paths:
      deleted:
          - /subscribe
          - /api/{domain}/{project}/install-command
          - /register
      modified:
          /api/{domain}/{project}/badges/security-score:
              operations:
                  added:
                      - POST
                  modified:
                      GET:
                          tags:
                              deleted:
                                  - security
                          parameters:
                              deleted:
                                  cookie:
                                      - test
                                  header:
                                      - user
                                      - X-Auth-Name
                              modified:
                                  path:
                                      domain:
                                          schema:
                                              type:
                                                  from: string
                                                  to: integer
                                              format:
                                                  from: hyphen-separated list
                                                  to: non-negative integer
                                              description:
                                                  from: Hyphen-separated list of lowercase string
                                                  to: Non-negative integers (including zero)
                                              min:
                                                  from: null
                                                  to: 7
                                              pattern:
                                                  from: ^(?:([a-z]+-)*([a-z]+)?)$
                                                  to: ^(?:\d+)$
                                  query:
                                      filter:
                                          content:
                                              schema:
                                                  properties:
                                                      modified:
                                                          color:
                                                              type:
                                                                  from: string
                                                                  to: number
                                      image:
                                          explode:
                                              from: null
                                              to: true
                                          schema:
                                              description:
                                                  from: alphanumeric
                                                  to: alphanumeric with underscore, dash, period, slash and colon
                                      token:
                                          schema:
                                              anyOf: true
                                              type:
                                                  from: string
                                                  to: ""
                                              format:
                                                  from: uuid
                                                  to: ""
                                              description:
                                                  from: RFC 4122 UUID
                                                  to: ""
                                              pattern:
                                                  from: ^(?:[0-9a-f]{8}-[0-9a-f]{4}-[0-5][0-9a-f]{3}-[089ab][0-9a-f]{3}-[0-9a-f]{12})$
                                                  to: ""
                          responses:
                              added:
                                  - default
                              deleted:
                                  - "200"
                                  - "201"
              parameters:
                  deleted:
                      path:
                          - domain
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
  components:
      schemas:
          deleted:
              - rules
              - network-policies
      parameters:
          deleted:
              header:
                  - network-policies
      headers:
          deleted:
              - testc
              - new
              - test
      requestBodies:
          deleted:
              - reuven
      responses:
          deleted:
              - OK
      securitySchemes:
          deleted:
              - OAuth
              - bearerAuth
summary:
  diff: true
  details:
      headers:
          deleted: 3
      parameters:
          deleted: 1
      paths:
          deleted: 3
          modified: 1
      requestBodies:
          deleted: 1
      responses:
          deleted: 1
      schemas:
          deleted: 2
      security:
          deleted: 1
      securitySchemes:
          deleted: 2
      servers:
          deleted: 1
      tags:
          deleted: 2
```

## Embedding into your Go program
```
diff.Get(&diff.Config{}, spec1, spec2)
```
See full example: [main.go](main.go)

## Notes
1. oasdiff expects [OpenAPI References](https://swagger.io/docs/specification/using-ref/) to be resolved.  
You can resolve refs using [this function](https://pkg.go.dev/github.com/getkin/kin-openapi/openapi3#SwaggerLoader.ResolveRefsIn).

2. oasdiff ignores changes to [Examples](https://swagger.io/specification/#example-object) and [Extensions](https://swagger.io/specification/#specification-extensions) by default. You can change this behavior by [configuration](diff/config.go).

3. Work in progress
While most aspects of OpenAPI Spec are already supported by this diff tool, some are still missing, notably: Examples, ExternalDocs, Links, Variables and a couple more.  
Pull requests are welcome!
