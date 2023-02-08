package checker

import (
	"bufio"
	"os"
	"strings"

	"github.com/TwiN/go-color"
)

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
		for errIndex, err := range errs {
			if err.Level != level {
				continue
			}

			uncolorizedText := strings.ReplaceAll(err.Text, color.Bold, "")
			uncolorizedText = strings.ReplaceAll(uncolorizedText, color.Reset, "")

			if strings.Contains(ignoreLine, strings.ToLower(err.Operation+" "+err.Path)) &&
				strings.Contains(ignoreLine, strings.ToLower(uncolorizedText)) {
				ignoredErrs[errIndex] = true
				break
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
