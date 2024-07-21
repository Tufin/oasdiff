package formatters

import "github.com/tufin/oasdiff/checker"

type Format string

const (
	FormatYAML          Format = "yaml"
	FormatJSON          Format = "json"
	FormatText          Format = "text"
	FormatMarkup        Format = "markup"
	FormatSingleLine    Format = "singleline"
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
	ColorMode checker.ColorMode
}

func NewRenderOpts() RenderOpts {
	return RenderOpts{
		ColorMode: checker.ColorAuto,
	}
}
