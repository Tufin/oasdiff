package diff

// Config includes various settings to control the diff
type Config struct {
	IncludeExamples   bool      // whether to include examples in diff or not (ignored by default)
	IncludeExtensions StringSet // which extensions to include in the diff (default is none) - see https://swagger.io/specification/#specification-extensions
	Filter            string    // diff will only include paths that match this regex (optional)
	Prefix            string    // a prefix that exists in s1 paths but not in s2 (optional)
}

// NewConfig returns a default configuration
func NewConfig() *Config {
	return &Config{
		IncludeExtensions: StringSet{},
	}
}
