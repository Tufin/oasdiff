package formatters

type Format string

const (
	FormatYAML          Format = "yaml"
	FormatJSON          Format = "json"
	FormatText          Format = "text"
	FormatHTML          Format = "html"
	FormatGithubActions Format = "githubactions"
	FormatJUnit         Format = "junit"
	FormatSarif         Format = "sarif"
)

// FormatterOpts can be used to pass properties to the formatter (e.g. colors)
type FormatterOpts struct {
	Language string
}

// RenderOpts can be used to pass properties to the renderer method
type RenderOpts struct {
}
