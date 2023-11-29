package checker

import (
	"fmt"

	"github.com/TwiN/go-color"
)

type ColorMode int

const (
	ColorAlways ColorMode = iota
	ColorNever
	ColorAuto
)

func NewColorMode(color string) ColorMode {
	switch color {
	case "always":
		return ColorAlways
	case "never":
		return ColorNever
	case "auto":
		return ColorAuto
	default:
		return ColorNever
	}
}

func isColorEnabled(colorMode ColorMode) bool {
	switch colorMode {
	case ColorAlways:
		return true
	case ColorNever:
		return false
	case ColorAuto:
		return !isPipedOutput()
	default:
		return false
	}
}

func colorizedValues(args []any) []any {
	result := make([]any, len(args))
	for i, arg := range args {
		result[i] = color.InBold(fmt.Sprintf("'%s'", interfaceToString(arg)))
	}
	return result
}

func quotedValues(args []any) []any {
	result := make([]any, len(args))
	for i, arg := range args {
		result[i] = fmt.Sprintf("'%s'", interfaceToString(arg))
	}
	return result
}
