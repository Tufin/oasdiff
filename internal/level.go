package internal

import "errors"

type Level string

const (
	LevelErr  = "ERR"
	LevelWarn = "WARN"
	LevelInfo = "INFO"
)

// String is used both by fmt.Print and by Cobra in help text
func (level *Level) String() string {
	return string(*level)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (level *Level) Set(v string) error {
	switch v {
	case LevelErr, LevelWarn, LevelInfo:
		*level = Level(v)
		return nil
	default:
		return errors.New(`must be "ERR", "WARN", or "INFO"`)
	}
}

// Type is only used in help text
func (level *Level) Type() string {
	return "level"
}
