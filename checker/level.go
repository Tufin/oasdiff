package checker

import (
	"fmt"

	"github.com/TwiN/go-color"
)

type Level int

const (
	ERR  Level = 3
	WARN Level = 2
	INFO Level = 1
)

func NewLevel(level string) (Level, error) {
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

func (level Level) String() string {
	switch level {
	case ERR:
		return "error"
	case WARN:
		return "warning"
	case INFO:
		return "info"
	default:
		return "issue"
	}
}

func (level Level) PrettyString() string {
	if isPipedOutput() {
		return level.String()
	}

	levelName := level.String()
	switch level {
	case ERR:
		return color.InRed(levelName)
	case WARN:
		return color.InPurple(levelName)
	case INFO:
		return color.InCyan(levelName)
	default:
		return color.InGray(levelName)
	}
}
