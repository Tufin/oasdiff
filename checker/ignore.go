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

	foundFirst := false
	ignoreLine := ""

	for _, err := range errs {
		if err.Level != level {
			result = append(result, err)
			continue
		}

		uncolorizedText := strings.ReplaceAll(err.Text, color.Bold, "")
		uncolorizedText = strings.ReplaceAll(uncolorizedText, color.Reset, "")

		found := false
		for ignoreScanner.Scan() {
			ignoreLine = ignoreScanner.Text()
			if strings.Contains(strings.ToLower(ignoreLine), strings.ToLower(err.Operation+" "+err.Path)) &&
				strings.Contains(strings.ToLower(ignoreLine), strings.ToLower(uncolorizedText)) {
				found = true
				foundFirst = true
				break
			}

			// after we found the first ignorance, we require ignorances to be one after for all breaking changes another without any lines in between
			if foundFirst {
				break
			}
		}
		if !found {
			result = append(result, err)
		}
	}
	return result, nil
}
