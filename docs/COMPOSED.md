## Composed Mode
Composed mode compares two collections of OpenAPI specs instead of a pair of specs in the default mode.
The collections are specified using a [glob](https://en.wikipedia.org/wiki/Glob_(programming)).
This can be useful when your APIs are defined across multiple files, for example, when multiple services, each one with its own spec, are exposed behind an API gateway, and you want to check changes across all the specs at once.

Notes: 
1. Composed mode compares only [paths and endpoints](DIFF.md#paths-vs-endpoints), other resources are compared only if referenced from the paths or endpoints.
2. Composed mode doesn't support [Path Prefix Modification](PATH-PREFIX.md) 
3. Learn more about how oasdiff [matches endpoints to each other](MATCHING-ENDPOINTS.md)
