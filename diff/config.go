package diff

import (
	"github.com/tufin/oasdiff/utils"
)

// Config includes various settings to control the diff
type Config struct {
	IncludeExtensions       utils.StringSet
	PathFilter              string
	FilterExtension         string
	PathPrefixBase          string
	PathPrefixRevision      string
	PathStripPrefixBase     string
	PathStripPrefixRevision string
	ExcludeElements         utils.StringSet
	IncludePathParams       bool
	Unchanged               bool
}

const (
	ExcludeExamplesOption    = "examples"
	ExcludeDescriptionOption = "description"
	ExcludeEndpointsOption   = "endpoints"
	ExcludeTitleOption       = "title"
	ExcludeSummaryOption     = "summary"
)

var ExcludeDiffOptions = []string{
	ExcludeExamplesOption,
	ExcludeDescriptionOption,
	ExcludeEndpointsOption,
	ExcludeTitleOption,
	ExcludeSummaryOption,
}

// NewConfig returns a default configuration
func NewConfig() *Config {
	return &Config{
		IncludeExtensions: utils.StringSet{},
		ExcludeElements:   utils.StringSet{},
	}
}

func (config *Config) WithExcludeElements(excludeElements []string) *Config {
	config.ExcludeElements = utils.StringList(excludeElements).ToStringSet()
	return config
}

func (config *Config) IsExcludeExamples() bool {
	return config.ExcludeElements.Contains(ExcludeExamplesOption)
}

func (config *Config) IsExcludeDescription() bool {
	return config.ExcludeElements.Contains(ExcludeDescriptionOption)
}

func (config *Config) IsExcludeEndpoints() bool {
	return config.ExcludeElements.Contains(ExcludeEndpointsOption)
}

func (config *Config) IsExcludeTitle() bool {
	return config.ExcludeElements.Contains(ExcludeTitleOption)
}

func (config *Config) IsExcludeSummary() bool {
	return config.ExcludeElements.Contains(ExcludeSummaryOption)
}

const (
	SunsetExtension          = "x-sunset"
	XStabilityLevelExtension = "x-stability-level"
	XExtensibleEnumExtension = "x-extensible-enum"
)

func (config *Config) WithCheckBreaking() *Config {
	config.IncludeExtensions.Add(XStabilityLevelExtension)
	config.IncludeExtensions.Add(SunsetExtension)
	config.IncludeExtensions.Add(XExtensibleEnumExtension)

	return config
}

func (config *Config) WithUnchanged() *Config {
	config.Unchanged = true
	return config
}
