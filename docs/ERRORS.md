## Errors

### Invalid Specs
Oasdiff expects valid OpenAPI 3 specs as input.  
The specs can be written in JSON or YAML.  
Oasdiff may return an error when given invalid specs, for example:
```
Error: failed to load base spec from "spec.yaml": error converting YAML to JSON: yaml: line 2: mapping values are not allowed in this context
```
The reason for this error is that the underlying library, [kin-openapi3](https://github.com/getkin/kin-openapi), converts YAML specs to JSON before parsing them.

### Circular Schema References
Schemas may reference themselves directly or indirectly.  
If the circle is too complex, you may receive this kind of error:
```
Error: failed to load base spec from "data/circular3.yaml": kin-openapi bug found: circular schema reference not handled with length 4 - #/components/schemas/circular2 -> #/components/schemas/circular3 -> #/components/schemas/circular1 -> #/components/schemas/circular2
```
To mitigate this problem, try increasing the value of `--max-circular-dep` (the default is 5).