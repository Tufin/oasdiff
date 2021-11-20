package diff

// Config includes various settings to control the diff
type Config struct {
	ExcludeExamples    bool      // whether to exclude changes to examples (included by default)
	ExcludeDescription bool      // whether to exclude changes to descriptions (included by default)
	IncludeExtensions  StringSet // which extensions to include in the diff (default is none) - see https://swagger.io/specification/#specification-extensions
	PathFilter         string    // diff will only include paths that match this regex (optional)
	PathPrefix         string    // a prefix that exists in first spec paths but not in second one (optional)
	BreakingOnly       bool      // whether to calc breaking changes only

}

// NewConfig returns a default configuration
func NewConfig() *Config {
	return &Config{
		ExcludeExamples:    false,
		ExcludeDescription: false,
		IncludeExtensions:  StringSet{},
		BreakingOnly:       false,
	}
}
