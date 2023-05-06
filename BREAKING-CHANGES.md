## Breaking Changes [Beta]
Breaking changes are changes that could break a client that is relying on the OpenAPI specification.  
[See some examples of breaking and non-breaking changes](BREAKING-CHANGES-EXAMPLES.md).  
Notes: 
1. This is a Beta feature, please report issues
2. There are two different methods for detecting breaking changes (see below)


### Old Method
The original implementation with the `-breaking-only` flag.
While this method is still supported, the new one will eventually replace it.

### New Method
An improved implementation for detecting breaking changes with the `-check-breaking` flag.
This method works differently from the old one: it is more accurate, generates nicer human-readable output, and is easier to maintain and extend.

To use it, run oasdiff with the `-check-breaking` flag, e.g.:
```
oasdiff -check-breaking -base https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml -revision https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```

There are two levels of breaking changes:
- `WARN` - Warning are potential breaking changes which developers should be aware of, but cannot be confirmed programmatically
- `ERR` - Errors are definite breaking changes which should be avoided

To exit with return code 1 when any ERR-level breaking changes are found, add the `-fail-on-diff` flag.  
To exit with return code 1 even if only WARN-level breaking changes are found, add the `-fail-on-diff` and `-fail-on-warns` flags.


### Stability Level
When a new API is introduced, you may want to allow developers to change its behavior without triggering a breaking-change error.  
The new Breaking Changes method provides this feature through the `x-stability-level` extension.  
There are four stability levels: `draft`->`alpha`->`beta`->`stable`.  
APIs with the levels `draft` or `alpha` can be changed freely without triggering a breaking-change error.  
Stability level may be increased, but not decreased, like this: `draft`->`alpha`->`beta`->`stable`.  
APIs with no stability level will trigger breaking changes errors upon relevant change.  
APIs with no stability level can be changed to any stability level.  

Example:
   ```
   /api/test:
    post:
     x-stability-level: "alpha"
   ```

### Ignoring Specific Breaking Changes
Sometimes, you may want to ignore certain breaking changes.  
The new Breaking Changes method allows you define breaking changes that you want to ignore in a configuration file.  
You can specify the configuration file name in the oasdiff command-line with the `-warn-ignore` flag for WARNINGS or the `-err-ignore` flag for ERRORS.  
Each line in the configuration file should contain two parts:
1. method and path
2. description of the breaking change

For example:
```
GET /api/{domain}/{project}/badges/security-score removed the success response with the status '200'
```

The line may contain additional info, like this:
```
 - 12.01.2023 In the GET /api/{domain}/{project}/badges/security-score, we removed the success response with the status '200'
```

The configuration files can be of any text type, e.g., Markdown, so you can use them to document breaking changes and other important changes.

### Breaking Changes to Enum Values
The new Breaking Changes method support rules for enum changes using the `x-extensible-enum` extension.  
This method allows adding new entries to enums used in responses which is very usable in many cases but requires clients to support a fallback to default logic when they receive an unknown value.
`x-extensible-enum` was introduced by [Zalando](https://opensource.zalando.com/restful-api-guidelines/#112) and picked up by the OpenAPI community. Technically, it could be replaced with anyOf+classical enum but the `x-extensible-enum` is a more explicit way to do it.  
In most cases the `x-extensible-enum` is similar to enum values, except it allows adding new entries in messages sent to the client (responses or callbacks).
If you don't use the `x-extensible-enum` in your OpenAPI specifications, nothing changes for you, but if you do, oasdiff will identify breaking changes related to `x-extensible-enum` parameters and properties.

### Optional Breaking-Changes Checks
You can use the `-include-checks` flag to include the following optional checks:
- response-non-success-status-removed
- api-operation-id-removed
- api-tag-removed
- response-property-enum-value-removed
- response-mediatype-enum-value-removed
- request-body-enum-value-removed

For example:
```
oasdiff -include-checks response-non-success-status-removed -check-breaking -base data/openapi-test1.yaml -revision data/openapi-test3.yaml
```


### Advantages of the New Breaking Changes Method 
- output is human readable
- supports localization for error messages and ignored changes
- [checks can be customized by developers](#customizing-breaking-changes-checks)
- fewer false-positive errors by design
- improved support for type changes: allows changing integer->number for json/xml properties, allows changing parameters (e.g. query/header/path) to type string from number/integer/etc.
- allows removal of responses with non-success codes (e.g., 503, 504, 403)
- allows adding new content-type to request
- easier to extend and customize
- will continue to be improved

### Limitations of the New Breaking Changes Method
- no checks for `context` instead of `schema` for request parameters
- no checks for `callback`s

## Customizing Breaking-Changes Checks
If you encounter a change that isn't considered breaking by oasdiff and you would like to consider it as a breaking-change you may add an [optional breaking-changes check](#optional-breaking-changes-checks).  
For more information, see [this guide](CUSTOMIZING-CHECKS.md) and this example of adding a custom check: https://github.com/Tufin/oasdiff/pull/208/files

