package checker

import (
	"bufio"
	"os"
	"strings"
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

func ProcessIgnoredBackwardCompatibilityErrors(level Level, errs IBackwardCompatibilityErrors, ignoreFile string) (IBackwardCompatibilityErrors, error) {
	result := make(IBackwardCompatibilityErrors, 0)

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
			if err.GetLevel() != level {
				continue
			}

			if err.MatchIgnore(ignorePath, ignoreLine) {
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
