# API Changelog {{ .BaseVersion }}{{ if ne .BaseVersion .RevisionVersion }} vs. {{ .RevisionVersion }}{{ end }}

{{ range $endpoint, $changes := .APIChanges }}
## {{ $endpoint.Operation }} {{ $endpoint.Path }}

{{ range $changes }}
- **{{ .Level }}**: {{ .Text }}
{{ end }}
{{ end }}
