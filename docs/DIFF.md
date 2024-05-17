## Diff
The `oasdiff diff` command displays the raw changes between OpenAPI specs.  
Output is fully detailed, typically in yaml or json but also available in text, markdown and html formats.  
This commmand is typically used to generate a structured diff report which can be consumed by other tools but it can also be viewed ny humans.

### Output Formats
The default diff output format, YAML, provides a full view of all diff details.  
Note that no output in the YAML format signifies that the diff is empty, or, in other words, there are no changes.  
Other formats: text, markdown, and HTML, are designed to be more user-friendly by providing only the most important parts of the diff, in a simplified format.  
If you wish to include additional details in non-YAML formats, please open an issue.  
The JSON format works only with `-exclude-elements endpoints` and is intended as a workaround for YAML complex mapping keys which aren't supported by some libraries (see comment at end of next section for more details).

### Preventing Changes
A common way to use `oasdiff diff` is by running it as a step the CI/CD pipeline to detect changes.  
In order to prevent changes, `oasdiff diff` can be configured to return an error if changes are found.  
To exit with return code 1 if any changes are found, add the `--fail-on-diff` flag.  

### Paths vs. Endpoints
OpenAPI Specification has a hierarchical model of [Paths](https://swagger.io/specification/#paths-object) and [Operations](https://swagger.io/specification/#operation-object) (HTTP methods).  
Oasdiff follows this hierarchy and displays a hierarchical diff with path changes: added, deleted and modified, and within the latter, "modified" section, another set of operation changes: added, deleted and modified. For example:
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
Oasdiff also outputs an alternate simplified diff per "endpoint" which is a combination of Path + Operation, for example:
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

When using output format `json`, oasdiff excludes `endpoints` automatically.

### OpenAPI Extensions
Oasdiff tracks changes to [OpenAPI extensions](https://swagger.io/docs/specification/openapi-extensions/) by default. To disable this, see [Excluding Specific Kinds of Changes](#excluding-specific-kinds-of-changes).  
The diff format for OpenAPI extensions conforms with [JavaScript Object Notation (JSON) Patch](https://datatracker.ietf.org/doc/html/rfc6902#section-4.4f), for example:
```
endpoints:
    modified:
        ?   method: POST
            path: /example/callback
        :   extensions:
                modified:
                    x-amazon-apigateway-integration:
                        - oldValue: "201"
                          value: "200"
                          op: replace
                          from: ""
                          path: /responses/default/statusCode
                        - oldValue: http://api.example.com/v1/example/callback
                          value: http://api.example.com/v1/example/calllllllllback
                          op: replace
                          from: ""
                          path: /uri
```

### Excluding Specific Kinds of Changes 
You can use the `--exclude-elements` flag to exclude certain kinds of changes:
- Use `--exclude-elements examples` to exclude [Examples](https://swagger.io/specification/#example-object)
- Use `--exclude-elements extensions` to exclude [Extensions](https://swagger.io/specification/#specification-extensions)
- Use `--exclude-elements description` to exclude description fields
- Use `--exclude-elements title` to exclude title fields
- Use `--exclude-elements summary` to exclude summary fields
- Use `--exclude-elements endpoints` to exclude the [endpoints diff](#paths-vs-endpoints)

### Additional Options
- [Merging AllOf Schemas](ALLOF.md)
- [Merging common parameters from the path level into the operation level](COMMON-PARAMS.md)
- [Excluding some endpoints](EXCLUDING-ENDPOINTS.md)
- [Path parameter renaming](PATH-PARAM-RENAME.md)
- [Case-insensitive header comparison](HEADER-DIFF.md)
- [Comparing multiple specs](COMPOSED.md)
- [Running from docker](DOCKER.md)
- [Embedding in your go program](GO.md)
