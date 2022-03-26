package diff

// Config includes various settings to control the diff
type Config struct {
	ExcludeExamples    bool
	ExcludeDescription bool
	IncludeExtensions  StringSet
	PathFilter         string
	FilterExtension    string
	PathPrefix         string
	BreakingOnly       bool
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
