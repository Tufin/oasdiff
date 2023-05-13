## Matching Endpoints
oasdiff compares matching endpoints to each other.  
By default, matching **excludes** path parameter names.  
For example, the following endpoints will be compared to each other because they differ only by a path parameter names:
| File        | Method  | Path         |
| ----------- | ------- | ------------ |
| `-base`     | GET     | /pet/{petId} |
| `-resivion` | GET     | /pet/{id}    |

This capability allows oasdiff to compare matching endpoints even if their path parameters were renamed.

## Duplicate Endpoints
Because oasdiff compares matching endpoints to each other, it expects a single instance of each endpoint to appear in the `-base` spec or collection (see [Composed Mode](README.md#composed-mode)) and in the `-revision` spec or collection.
If duplicate matching endpoints are found in either `-base` or `-revision`, there are two options.

1. **Include** path parameter names in endpoint matching with the `-match-path-params` flag. In this case, the endpoints in the table above will be considered two different ones.

1. Use [`x-since-date`](#duplicate-endpoints-and-x-since-date)

## Duplicate Endpoints and `x-since-date`
If duplicate matching endpoints are found in either `-base` or `-resivion`, then oasdiff uses the endpoint with the most recent `x-since-date` value.

- The `x-since-date` extension can be set at the Path or Operation level.
- `x-since-date` extensions set on the Operation level override the value set on Path level.
- If an endpoint doesn't have the `x-since-date` extension, its value is set to the default: "2000-01-01".
- Duplicate endpoints with the same `x-since-date` value will trigger an error.
- The format of the `x-since-date` is the RFC3339 full-date format.

Example of the `x-since-date` usage:
   ```
   /api/test:
    get:
     x-since-date: "2023-01-11"
   ```

