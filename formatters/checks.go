package formatters

type Check struct {
	Id          string `json:"id" yaml:"id"`
	Level       string `json:"level" yaml:"level"`
	Description string `json:"description" yaml:"description"`
}

type Checks []Check

func (checks Checks) Len() int {
	return len(checks)
}

func (checks Checks) Less(i, j int) bool {
	return checks[i].Id < checks[j].Id
}

func (checks Checks) Swap(i, j int) {
	checks[i], checks[j] = checks[j], checks[i]
}
