package lint

type Config struct {
	Checks []Check
}

func NewConfig(checks []Check) *Config {
	return &Config{Checks: checks}
}

func DefaultConfig() *Config {
	return &Config{
		Checks: defaultChecks(),
	}
}

func defaultChecks() []Check {
	return []Check{
		SchemaCheck,
		PathParamsCheck,
		RequiredParamsCheck,
		InfoCheck,
	}
}
