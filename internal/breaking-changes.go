package internal

import (
	"fmt"
	"io"

	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/checker/localizations"
	"github.com/tufin/oasdiff/diff"
)

func handleBreakingChanges(stdout io.Writer, diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, inputFlags *InputFlags) (bool, *ReturnError) {

	c := checker.GetChecks(inputFlags.includeChecks)
	c.Localizer = *localizations.New(inputFlags.lang, "en")

	errs, returnErr := getBreakingChanges(c, diffReport, operationsSources, inputFlags.warnIgnoreFile, inputFlags.errIgnoreFile)
	if returnErr != nil {
		return false, returnErr
	}

	switch inputFlags.format {
	case FormatYAML:
		if err := printYAML(stdout, errs); err != nil {
			return false, getErrFailedToPrint("breaking changes YAML", err)
		}
	case FormatJSON:
		if err := printJSON(stdout, errs); err != nil {
			return false, getErrFailedToPrint("breaking changes JSON", err)
		}
	case FormatText:
		if len(errs) > 0 {
			fmt.Fprintf(stdout, c.Localizer.Get("messages.total-errors"), len(errs))
		}

		for _, bcerr := range errs {
			fmt.Fprintf(stdout, "%s\n\n", bcerr.PrettyErrorText(c.Localizer))
		}
	default:
		return false, getErrUnsupportedBreakingChangesFormat(inputFlags.format)
	}

	return errs.IsEmpty(inputFlags.failOnWarns), nil
}

func getBreakingChanges(c checker.BackwardCompatibilityCheckConfig, diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, warnIgnoreFile string, errIgnoreFile string) (checker.BackwardCompatibilityErrors, *ReturnError) {

	errs := checker.CheckBackwardCompatibility(c, diffReport, operationsSources)

	if warnIgnoreFile != "" {
		var err error
		errs, err = checker.ProcessIgnoredBackwardCompatibilityErrors(checker.WARN, errs, warnIgnoreFile)
		if err != nil {
			return nil, getErrCantProcessIgnoreFile("warn", err)
		}
	}

	if errIgnoreFile != "" {
		var err error
		errs, err = checker.ProcessIgnoredBackwardCompatibilityErrors(checker.ERR, errs, errIgnoreFile)
		if err != nil {
			return nil, getErrCantProcessIgnoreFile("err", err)
		}
	}

	return errs, nil
}
