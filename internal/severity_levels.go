package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/utils"
)

// processSeverityLevels reads a file with severity levels and returns a map of severity levels
func processSeverityLevels(file string) (map[string]checker.Level, *ReturnError) {

	if file == "" {
		return nil, nil
	}

	result := map[string]checker.Level{}

	f, err := os.Open(file)
	if err != nil {
		return nil, getErrFailedToLoadSeverityLevels(file, err)
	}
	defer f.Close()

	validIds := utils.StringList(checker.GetAllRuleIds()).ToStringSet()

	scanner := bufio.NewScanner(f)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		frags := strings.Fields(line)

		if len(frags) != 2 {
			return nil, getErrFailedToLoadSeverityLevels(file, fmt.Errorf("invalid line #%d: %s", lineNum, line))
		}

		id := frags[0]
		if !validIds.Contains(id) {
			return nil, getErrFailedToLoadSeverityLevels(file, fmt.Errorf("invalid rule id %q on line %d", id, lineNum))
		}

		level, err := checker.NewLevel(frags[1])
		if err != nil {
			return nil, getErrFailedToLoadSeverityLevels(file, fmt.Errorf("invalid level %q on line %d", frags[1], lineNum))
		}

		result[id] = level
	}

	return result, nil
}
