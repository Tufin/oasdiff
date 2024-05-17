## Filtering Endpoints

You can filter endpoints in two ways:

### By path name
Use the `--match-path` option to exclude paths that don't match the given regular expression.  
For example, this diff includes only endpoints containing "/api" in the path:
```
oasdiff diff https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml -f text --match-path "/api"
```
Note, that the filter is applied recursively at all levels. For example, if a path contains a callback, the filter will be applied both to the path itself and to the callback path.  
To include such a nested change, use a regular expression that contains both paths, for example -filter "path|callback-path"
   
### By extension
Use the `--filter-extension` option to exclude paths and operations with an OpenAPI Extension matching the given regular expression.  
For example, this diff excludes paths and operations with extension "x-beta"
```
oasdiff diff https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml -f text --filter-extension "x-beta"
```
Notes:
1. OpenAPI Extensions can be defined both at the path level and at the operation level. Both are matched and excluded with this flag.
2. If a path or operation has a matching extension only in one of the specs, but not in the other, it will appear as Added or Deleted.