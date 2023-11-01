package checker

import (
	"strings"

	"github.com/TwiN/go-color"
)

func GetUncolorizedText(c Change) string {
	uncolorizedText := strings.ReplaceAll(c.GetText(), color.Bold, "")
	return strings.ReplaceAll(uncolorizedText, color.Reset, "")
}
