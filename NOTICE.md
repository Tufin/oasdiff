## OASDiff for chargebee
- We use this tool to generate a diff in JSON format between two OpenAPI documents.
- This tool, by default, do not include diffs for custom extensions. However, with some code modification, it can support this, which is what we have done.
    - We just need to add the following to `func WithCheckBreaking()` in `internal/flags.go`
```
config.IncludeExtensions.Add("x-cb-ui-key")
config.IncludeExtensions.Add("x-cb-obs-attributes")
config.IncludeExtensions.Add("x-cb-spec-domain")
config.IncludeExtensions.Add("x-cb-is-eap")
```
- Build the Docker Image and push it to ECR, in order to use it.