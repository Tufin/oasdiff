package checker

import (
	"bufio"
	"os"
	"strings"

	"github.com/TwiN/go-color"
)

func ignoreLinePath(ignoreLine string) string {
	ignoreComponents := strings.Fields(ignoreLine)
	pathIndex := -1

	for i, component := range ignoreComponents {
		if strings.HasPrefix(component, "/") {
			pathIndex = i
			break
		}
	}

	if pathIndex == -1 {
		return ""
	}

	return ignoreComponents[pathIndex]
}

func ProcessIgnoredBackwardCompatibilityErrors(level int, errs []BackwardCompatibilityError, ignoreFile string) ([]BackwardCompatibilityError, error) {
	result := make([]BackwardCompatibilityError, 0)

	ignore, err := os.Open(ignoreFile)
	if err != nil {
		return nil, err
	}
	defer ignore.Close()
	ignoreScanner := bufio.NewScanner(ignore)

	ignoredErrs := make([]bool, len(errs))
	for ignoreScanner.Scan() {
		ignoreLine := strings.ToLower(ignoreScanner.Text())

		ignorePath := ignoreLinePath(ignoreLine)
		if ignorePath == "" {
			continue
		}

		for errIndex, err := range errs {
			if err.Level != level {
				continue
			}

			uncolorizedText := strings.ReplaceAll(err.Text, color.Bold, "")
			uncolorizedText = strings.ReplaceAll(uncolorizedText, color.Reset, "")

			if ignorePath == strings.ToLower(err.Path) &&
				strings.Contains(ignoreLine, strings.ToLower(err.Operation+" "+err.Path)) &&
				strings.Contains(ignoreLine, strings.ToLower(uncolorizedText)) {
				ignoredErrs[errIndex] = true
			}
		}
	}

	for errIndex, err := range errs {
		if !ignoredErrs[errIndex] {
			result = append(result, err)
		}
	}
	return result, nil
}
