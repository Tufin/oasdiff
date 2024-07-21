# API Changelog {{ .BaseVersion }}{{ if ne .BaseVersion .RevisionVersion }} vs. {{ .RevisionVersion }}{{ end }}
{{ range $endpoint, $changes := .APIChanges }}
## {{ $endpoint.Operation }} {{ $endpoint.Path }}
{{ range $changes }}- {{ if .IsBreaking }}:warning:{{ end }} {{ .Text }}
{{ end }}
{{ end }}
