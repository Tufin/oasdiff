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
      display breaking changes only (deprecated, use check-breaking instead)
  -check-breaking
      check diff for breaking changes with breaking changes checks
  -composed
      work in 'composed' mode, compare paths in all specs in the base and revision directories. In such mode the base and the revision parameters could be Globs instead of files
  -deprecation-days int
    	minimal number of days required between deprecating a resource and removing it without being considered 'breaking'
  -exclude-description
    	ignore changes to descriptions
  -exclude-examples
    	ignore changes to examples
  -fail-on-diff
    	fail with exit code 1 if a difference is found
  -fail-on-warns
      fail with exit code 1 if only WARN breaking changes found, the option used only with both -check-breaking and -fail-on-diff options
  -filter string
    	if provided, diff will include only paths that match this regular expression
  -filter-extension string
    	if provided, diff will exclude paths and operations with an OpenAPI Extension matching this regular expression
  -format string
    	output format: yaml, text or html (default "yaml")
  -max-circular-dep int
    	maximum allowed number of circular dependencies between objects in OpenAPI specs (default 5)
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
  -warn-ignore string
    the filename for check breaking ignore file for warnings
  -err-ignore string
    the filename for check breaking ignore file for errors
  -version
    	show version and quit
```
All arguments can be passed with one or two leading minus signs.  
For example, ```-breaking-only``` and ```--breaking-only``` are equivalent.

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

### Check for breaking changes
```bash
oasdiff -check-breaking -base https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml -revision https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```
See [breaking changes](#breaking-changes-checks)

### Check multiple specs' paths for breaking changes
```bash
oasdiff -check-breaking -base base/**/*.yaml -revision revision/**/.*yaml
```
See [breaking changes](#breaking-changes-checks)

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
Other formats: text, markdown, and HTML, are designed to be more user-friendly by providing only the most important parts of the diff, in a simplified format.  
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
Some YAML libraries don't support complex mapping keys:
- python PyYAML: see https://github.com/Tufin/oasdiff/issues/94#issuecomment-1087468450
- golang gopkg.in/yaml.v3 fails to unmarshal the oasdiff output. This package offers a solution: https://github.com/tliron/yamlkeys

### Breaking Changes
Breaking changes are changes that could break a client that is relying on the OpenAPI specification.  
[See some examples of breaking and non-breaking changes](breaking-changes.md).  
Note: this is a Beta feature. Please report issues.

### Breaking Changes Checks
An improved method for detecting breaking changes in specifications. It works differently form the original "Breaking changes" feature.
It uses the list of rules, each checks specific case and could generate more decent human-readable error messages.
To use it, run the oasdiff with the `-check-breaking` option, e.g.:
```
oasdiff -check-breaking -base data/deprecation/base.yaml -revision data/deprecation/deprecated-past.yaml
```

There are two levels of check errors:
- `WARN` - for important changes which should be noticed by developers but may not be confirmed programmatically
- `ERR` - for most of breaking changes which should be avoided

If you use the `-fail-on-diff` option with the `-check-breaking`, oasdiff will exit with return code 1 when any level of error occured.
If you want to make oasdiff exit with non-zero return code only when there are errors with ERR level, use the `-fail-on-warns` option.
It makes oasdiff fail with exit code 1 if only WARN breaking changes found. The option is usable only when both `-check-breaking` and `-fail-on-diff` options used.

#### Composed mode
The "composed mode" could be useful when you want to check for breaking changes multiple specs describing the single API.
For example, your application could provide multiple APIs decribed in multiple files. Or your API is combined from multiple APIs with API gateway and you want to check for breaking changes your API overally, not single specs.
When composed mode is used, oasdiff before computing diff (and breaking changes checks)

There are some problems that could be in such usage scenario:
- Different files could have different components (schemas, properties, etc.) with the same name. So to check for difference we should not only internalize all referencies, but replace local referencies with its contenttoo for each API description
- To check for backward compatibility, we should use only API (path+method) difference
- We should support some way to move the API from one file to another without breaking changes. Sometimes it could not be done atomically (e.g. when specifications are in different repositories maintained by different teams). So we should have some rules to order similar API descriptions to find the latest one.

The composed mode allows to do it.
It is helpfull when you want to find diff and check for breaking changes for API divided into multiple files.
If there are same paths in different OpenAPI objects, then function uses version of the path with the last x-since-date extension.
The `x-since-date` extension should be set on path or operations level. Extension set on the operations level overrides the value set on path level.
If such path doesn't have `the x-since-date` extension, its value is default "2000-01-01"
If there are same paths with the same x-since-date value, then function returns error.
The format of the `x-since-date` is the RFC3339 full-date format

Example of the `x-since-date` usage:
   ```
   /api/test:
    get:
     x-since-date: "2023-01-11"
   ```

#### Stability level
Breaking changes checks method supports `x-stability-level` extension for APIs which allows ignore breaking changes for unstable APIs.
There are 4 levels: `draft`->`alpha`->`beta`->`stable`.
Any changes for APIs with the levels `draft` or `alpha` are allowed.
Stability level may be increased, but not decreased, like this: `draft`->`alpha`->`beta`->`stable`
If there is no specified stability level for an API, breaking changes are disallowed, but the stability level can be set to any level in this case.

Example:
   ```
   /api/test:
    post:
     x-stability-level: "alpha"
   ```

#### Breaking changes journal and ignorance
Sometimes there are detected some false positive breaking changes which should be ignored until the problem fixed in oasdiff. Or there are some important changes which should be noticed by developers but may not be confirmed programmatically as broken changes.
By default oasdiff with the `-check-breaking` option exists with non-zero status when there are some errors. But you could add any detected breaking change to ignorance file to ignore it during breaking changes checks.
To do so, you should add to the end of ignorance file the lines for all errors (on per error) which contains two parts for each error:
- method and path
- description of the breaking change

Like this:
```
GET /api/{domain}/{project}/badges/security-score removed the success response with the status '200'
```

The line could additionally contains additional simbols in between of required parts, like this:
```
 - 12.01.2023 In the GET /api/{domain}/{project}/badges/security-score, we removed the success response with the status '200'
```

Note that for multiple errors, added lines should be in the same order which oasdiff reports errors and should be without any other lines in between.
To specify the filename that should be used as the ignorance file, use the `-warn-ignore` option for WARNINGS and the `-err-ignore` option for ERRORS.
The ignorance files could be of any text type, e.g. Markdown, so you can use them as breaking and important changes journals.

#### x-extensible-enum support
The breaking changes has rules specific to enum changes which recommends using the `x-extensible-enum`.
Using such enums allows to add new entries to enums used in responses which is very usable in many case, but requires clients to support fallback to default logic when they receives unknown value.
The `x-extensible-enum` was introduced by the [Zalando](https://opensource.zalando.com/restful-api-guidelines/#112) and picked up by the OpenAPI community.
Technically, it could be replaced with anyOf+classical enum but the `x-extensible-enum` is more explicit way to do it.
Most tools doesn't support the `x-extensible-enum` but the breaking changes checks do. In most case it treats it similar to enum values, except it allows adding new entries in messages sent to client (responses or callbacks).
If you don't use the `x-extensible-enum` in your OpenAPI specifications, nothing changed for you, but if you do oasdiff will find you breaking changes realted to `x-extensible-enum` parameters and properties.

#### Othe differencies from the original implementation
There are multiple differencies using this feature in comparison with the original implementation:
- output in human readable format.
- supports localization for error messages and ignored changes.
- the set of checks can be modified by developers usign oasdiff as library with their own specific checks by adding/removing checks from the slice of checks.
- fewer false positive errors by design.
- better support for type changes checks: allows changing integer->number for json/xml properties, allows changing parameters (e.g. query/header/path) to type string from number/integer/etc.
- allows removal of responses with non-success codes (e.g. 503, 504, 403).
- allows adding new content-type to request (with the kept current).
- easier to extend and customize
- will continue to be improved 

#### Current limitations (going to be fixed in the nearest future):
- there are no checks for `context` instead of `schema` for request parameters
- there are no checks for `callback`s
- not fixed false positive breaking change error when the path parameter renamed both in path and in parameters section to the same name, this can be mitigated with the checks errors ignore feature
- doesn't support Path Prefix Modification, this can be mitigated with checks errors ignore feature 

### Non-Breaking Removal of Deprecated Resources
Sometimes APIs need to be removed, for example, when we replace an old API by a new version.
As API owners, we want a process that will allow us to phase out the old API version and transition to the new one smoothly as possible and with minimal disruptions to business.

OpenAPI specification supports a ```deprecated``` flag which can be used to mark operations and other object types as deprecated.  
Normally, deprecation **is not** considered a breaking change since it doesn't break the client but only serves as an indication of an intent to remove something in the future, in contrast, the eventual removal of a resource **is** considered a breaking change.

oasdiff allows you to gracefully remove a resource without getting the ```breaking-change``` warning, as follows:
1. First, the resource is marked as ```deprecated``` and a [special extension](https://swagger.io/specification/#specification-extensions) ```x-sunset``` is added to announce the date at which the resource will be removed
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
oasdiff -deprecation-days=30 -breaking-only -base data/deprecation/base.yaml -revision data/deprecation/deprecated-past.yaml
```

Setting deprecation-days to 0 is equivalent to the default which allows non-breaking deprecation regardless of the sunset date.  
Note: this is a Beta feature. Please report issues.

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

For breaking changes checks:
```go
diffConfig := &diff.Config{}
diffConfig.IncludeExtensions.Add(checker.XStabilityLevelExtension)
diffConfig.IncludeExtensions.Add(diff.SunsetExtension)
diffConfig.IncludeExtensions.Add(checker.XExtensibleEnumExtension)

diffRes, operationsSources, err := diff.GetPathsDiff(diffConfig, baseContracts, revisionContracts)
if err != nil {
  // process
}

c := checker.DefaultChecks() // here could be added additional checks
c.Localizer = *localizations.New("en", "en")
errs := checker.CheckBackwardCompatibility(c, diffRes, operationsSources)

// process WARN ignorance file
errs, err := processIgnorance(bcWarningIgnoranceFile, checker.WARN, errs)
if err != nil {
  // process
}

// process ERR ignorance file
errs, err = processIgnorance(bcErrorIgnoranceFile, checker.ERR, errs)
if err != nil {
  // process
}

// pretty print breaking changes errors
if len(errs) > 0 {
  fmt.Printf(c.Localizer.Get("messages.total-errors"), len(errs))
  for _, bcerr := range errs {
    fmt.Printf("%s\n\n", bcerr.PrettyError(c.Localizer))
  }
}

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
