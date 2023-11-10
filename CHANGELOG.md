## Changelog
As your API evolves, it will undergo changes. Some of these changes may be "breaking" while others are not.  
The changelog provides a list of all significant changes between two versions of the OpenAPI specification, including non-breaking changes.

To generate the changelog between two specs run oasdiff with the `changelog` command:
```
oasdiff changelog https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```

The changes are categorized into three levels:
- `INFO` - Changes which are not breaking yet still important to know about
- `WARN` - Warnings are potential breaking changes which developers should be aware of, but cannot be confirmed programmatically
- `ERR` - Errors are definite breaking changes which should be avoided

The changelog is actually an extension of the breaking changes output with additional `INFO`-level changes.  
See [the breaking-changes documentation](BREAKING-CHANGES.md) for additional options that can be also be used with changelog command.

To see the full list of supported changes, run:
```
oasdiff checks
```

### Output Formats
The default changelog format is human-readable text.  
You can specify the `--format` flag to output the changelog as json or yaml.

### Customizing the Changelog
If you encounter a change that isn't logged by oasdiff you may add a [custom check](CUSTOMIZING-CHECKS.md).
