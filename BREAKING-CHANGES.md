## Breaking Changes
A breaking change is a change to a component, such as a server, that could break a dependent component, such as a client, for example deleting an endpoint. 
When working with OpenAPI, breaking-changes can be caught by monitoring changes to the specification.

**oasdiff detects over 100 kinds of breaking changes**

To detect breaking-changes between two specs run oasdiff with the `breaking` command:
```
oasdiff breaking https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```

There are two levels of breaking changes:
- `WARN` - Warnings are potential breaking changes which developers should be aware of, but cannot be confirmed programmatically
- `ERR` - Errors are definite breaking changes which should be avoided

To exit with return code 1 when any ERR-level breaking changes are found, add the `--fail-on ERR` flag.  
To exit with return code 1 even if only WARN-level breaking changes are found, add the `--fail-on WARN` flag.

### Output Formats
The default output format is human-readable text.  
You can specify the `--format` flag to output breaking-changes in json, yaml, githubactions or junit formats.

### API Stability Levels
When a new API is introduced, you may want to allow developers to change its behavior without triggering a breaking-change error.  
oasdiff provides this feature through the `x-stability-level` extension.  
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

### Deprecating APIs
oasdiff allows you to [deprecate APIs gracefully](API-DEPRECATION.md) without triggering a breaking-change error.

### Ignoring Specific Breaking Changes
Sometimes, you may want to ignore certain breaking changes.  
The new Breaking Changes method allows you define breaking changes that you want to ignore in a configuration file.  
You can specify the configuration file name in the oasdiff command-line with the `--warn-ignore` flag for WARNINGS or the `--err-ignore` flag for ERRORS.  
Each line in the configuration file should contain two parts:
1. method and path (the first field in the line beginning with slash)
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
You can use the `--include-checks` flag to include optional checks for breaking-changes.  
Use the `oasdiff checks --required false` command to see the list of optional checks.  

For example:
```
oasdiff breaking data/openapi-test1.yaml data/openapi-test3.yaml --include-checks response-non-success-status-removed
```

### Localization
To display changes in other languages use the `--lang` flag.  
Currently English and Russian are supported.  
[Please improve oasdiff by adding your own language](https://github.com/Tufin/oasdiff/issues/383).

### Customizing Breaking-Changes Checks
If you encounter a change that isn't considered breaking by oasdiff you may:
1. Check if the change is already available as an [optional breaking-changes check](#optional-breaking-changes-checks).  
2. Add a [custom check](CUSTOMIZING-CHECKS.md)

### Examples
[Here are some examples of breaking and non-breaking changes that oasdiff supports](BREAKING-CHANGES-EXAMPLES.md).  
This document is automatically generated from oasdiff unit tests.

### Known Limitations
- no checks for `context` instead of `schema` for request parameters
- no checks for `callback`s
