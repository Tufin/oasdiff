package checker

import (
	"fmt"
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
