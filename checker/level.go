package checker

import (
	"errors"
	"fmt"
)

type Level int

const (
	ERR  Level = 3
	WARN Level = 2
	INFO Level = 1
	NONE Level = 0
)

type LevelArg string

const (
	LevelERR  LevelArg = "ERR"
	LevelWARN LevelArg = "WARN"
	LevelINFO LevelArg = "INFO"
)

// String is used both by fmt.Print and by Cobra in help text
func (level *LevelArg) String() string {
	return string(*level)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (level *LevelArg) Set(v string) error {
	switch v {
	case "ERR", "WARN", "INFO":
		*level = LevelArg(v)
		return nil
	default:
		return errors.New(`must be one of "ERR", "WARN", or "INFO"`)
	}
}

// Type is only used in help text
func (level *LevelArg) Type() string {
	return "level"
}

func (level Level) HigherOrEqual(other Level) bool {
	return level >= other
}

func (level LevelArg) ToLevel() (Level, error) {
	switch level {
	case "ERR":
		return ERR, nil
	case "WARN":
		return WARN, nil
	case "INFO":
		return INFO, nil
	}
	return INFO, fmt.Errorf("invalid level %s", level)
}
