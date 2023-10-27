package formatters

type Check struct {
	Id          string `json:"id" yaml:"id"`
	Level       string `json:"level" yaml:"level"`
	Description string `json:"description" yaml:"description"`
	Required    bool   `json:"reuired" yaml:"reuired"`
}
