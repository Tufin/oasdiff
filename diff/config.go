package diff

// Config includes various settings to control the diff
type Config struct {
	ExcludeExamples         bool
	ExcludeDescription      bool
	IncludeExtensions       StringSet
	PathFilter              string
	FilterExtension         string
	PathPrefixBase          string
	PathPrefixRevision      string
	PathStripPrefixBase     string
	PathStripPrefixRevision string
	BreakingOnly            bool
	DeprecationDays         int
	ExcludeEndpoints        bool
}

// NewConfig returns a default configuration
func NewConfig() *Config {
	return &Config{
		ExcludeExamples:    false,
		ExcludeDescription: false,
		IncludeExtensions:  StringSet{},
		BreakingOnly:       false,
		ExcludeEndpoints:   false,
	}
}
