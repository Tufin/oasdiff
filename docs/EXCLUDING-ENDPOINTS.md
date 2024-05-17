## Excluding Specific Endpoints
You can filter endpoints in two ways:
1. By path name: use the `--match-path` option to exclude paths that don't match the given regular expression, see [example](#openapi-diff-for-endpoints-containing-api-in-the-path)
2. By extension: use the `--filter-extension` option to exclude paths and operations with an OpenAPI Extension matching the given regular expression, see [example](#exclude-paths-and-operations-with-extension-x-beta)
