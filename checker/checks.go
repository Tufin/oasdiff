package checker

import (
	"github.com/tufin/oasdiff/diff"
)

// BackwardCompatibilityCheck, or a check, is a function that receives a diff report and returns a list of changes
type BackwardCompatibilityCheck func(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes

type BackwardCompatibilityChecks []BackwardCompatibilityCheck
