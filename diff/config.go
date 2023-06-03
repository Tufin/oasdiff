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
	BreakingOnly            bool
	DeprecationDays         int
	ExcludeElements         utils.StringSet
	MatchPathParams         bool
}

const (
	excludeExamplesOption    = "examples"
	excludeDescriptionOption = "description"
	excludeEndpointsOption   = "endpoints"
	excludeTitleOption       = "title"
	excludeSummaryOption     = "summary"
)

var excludeDiffOptions = utils.StringSet{
	excludeExamplesOption:    {},
	excludeDescriptionOption: {},
	excludeEndpointsOption:   {},
	excludeTitleOption:       {},
	excludeSummaryOption:     {},
}

// NewConfig returns a default configuration
func NewConfig() *Config {
	return &Config{
		IncludeExtensions: utils.StringSet{},
		BreakingOnly:      false,
		ExcludeElements:   utils.StringSet{},
	}
}

func (config *Config) SetExcludeElements(excludeElements utils.StringSet, excludeExamples, excludeDescription, excludeEndpoints bool) {
	config.ExcludeElements = excludeElements

	// backwards compatibility for deprecated flags
	if excludeExamples {
		config.ExcludeElements.Add(excludeExamplesOption)
	}
	if excludeDescription {
		config.ExcludeElements.Add(excludeDescriptionOption)
	}
	if excludeEndpoints {
		config.ExcludeElements.Add(excludeEndpointsOption)
	}
}

func (config *Config) IsExcludeExamples() bool {
	return config.ExcludeElements.Contains(excludeExamplesOption)
}

func (config *Config) IsExcludeDescription() bool {
	return config.ExcludeElements.Contains(excludeDescriptionOption)
}

func (config *Config) IsExcludeEndpoints() bool {
	return config.ExcludeElements.Contains(excludeEndpointsOption)
}

func (config *Config) IsExcludeTitle() bool {
	return config.ExcludeElements.Contains(excludeTitleOption)
}

func (config *Config) IsExcludeSummary() bool {
	return config.ExcludeElements.Contains(excludeSummaryOption)
}

func ValidateExcludeElements(input utils.StringList) utils.StringList {
	return input.ToStringSet().Minus(excludeDiffOptions).ToStringList().Sort()
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
