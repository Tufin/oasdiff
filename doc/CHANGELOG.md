## Changelog
As your API evolves, it will undergo changes. Some of these changes may be "breaking" while others are not.  
The changelog provides a list of all significant changes between two versions of the OpenAPI specification, including non-breaking changes.  

To generate the changelog between two specs run oasdiff with the `changelog` command:
```
oasdiff changelog https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```

The changes are categorized into three levels:
- `ERR` - Errors are definite breaking changes which should be avoided
- `WARN` - Warnings are potential breaking changes which developers should be aware of, but cannot be confirmed programmatically
- `INFO` - Changes which are not breaking yet still important to know about

The changelog is an extension of the breaking changes output with additional `INFO`-level changes.  
See [the breaking changes documentation](BREAKING-CHANGES.md) for additional options that can be also be used with changelog command.

To see the full list of supported changes, run:
```
oasdiff checks
```

### Output Formats
By default, changes are displayed as human-readable text with [color](#color).  
You can specify the `--format` flag to output changes in other formats: `json`, `yaml`, [`html`](https://html-preview.github.io/?url=https://github.com/tufin/oasdiff/blob/main/changelog.html), `githubactions` or `junit`.  
An additional format `singleline` displays each change on a single line, this can be useful to prepare [ignore files](BREAKING-CHANGES.md#ignoring-specific-breaking-changes)

### Color
When outputting changes to a Unix terminal, oasdiff automatically adds colors with ANSI color escape sequences.  
If output is piped into another process or redirected to a file, oasdiff disables color.  
To control color manually, use the `--color` flag with `always` or `never`.

### Customizing the Changelog
If you encounter a change that isn't logged by oasdiff you may add a [custom check](CUSTOMIZING-CHECKS.md).

### Additional Options
- [Merging AllOf Schemas](ALLOF.md)
- [Merging common parameters from the path level into the operation level](COMMON-PARAMS.md)
- [Comparing multiple specs](COMPOSED.md)
- [Running from docker](DOCKER.md)
- [Embedding in your go program](GO.md)