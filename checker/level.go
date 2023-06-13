package checker

import "fmt"

type Level int

const (
	ERR  Level = 3
	WARN Level = 2
	INFO Level = 1
	NONE Level = 0
)

func (level Level) IsValid() bool {
	return level != NONE
}

func (level Level) HigherOrEqual(other Level) bool {
	return level >= other
}

func (level *Level) Set(levelStr string) error {
	switch levelStr {
	case "ERR":
		*level = ERR
	case "WARN":
		*level = WARN
	case "INFO":
		*level = INFO
	default:
		*level = NONE
		return fmt.Errorf("invalid level: %s", levelStr)
	}
	return nil
}

func (level Level) String() string {
	switch level {
	case ERR:
		return "ERR"
	case WARN:
		return "WARN"
	case INFO:
		return "INFO"
	}
	return "NONE"
}

func (level Level) Type() string {
	return "int64"
}
