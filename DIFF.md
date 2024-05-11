## Diff
Displays the raw changes between a pair of OpenAPI specs.  
Output is fully detailed, typically in yaml or json but also available in text/markdown or html formats.  
Typically used as input for other tools but may also be viewed ny humans.

## Diff Output Formats
The default diff output format, YAML, provides a full view of all diff details.  
Note that no output in the YAML format signifies that the diff is empty, or, in other words, there are no changes.  
Other formats: text, markdown, and HTML, are designed to be more user-friendly by providing only the most important parts of the diff, in a simplified format.  
The JSON format works only with `-exclude-elements endpoints` and is intended as a workaround for YAML complex mapping keys which aren't supported by some libraries (see comment at end of next section for more details).
If you wish to include additional details in non-YAML formats, please open an issue.

## OpenAPI Extensions
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

## Excluding Specific Kinds of Changes 
You can use the `--exclude-elements` flag to exclude certain kinds of changes:
- Use `--exclude-elements examples` to exclude [Examples](https://swagger.io/specification/#example-object)
- Use `--exclude-elements extensions` to exclude [Extensions](https://swagger.io/specification/#specification-extensions)
- Use `--exclude-elements description` to exclude description fields
- Use `--exclude-elements title` to exclude title fields
- Use `--exclude-elements summary` to exclude summary fields
- Use `--exclude-elements endpoints` to exclude the [endpoints diff](#paths-vs-endpoints)

