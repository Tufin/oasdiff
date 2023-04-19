package internal

import (
	"fmt"
	"io"

	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/checker/localizations"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/utils"
)

func handleBreakingChanges(stdout io.Writer,
	diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap,
	includeChecks utils.StringList,
	format string,
	failOnWarns bool,
	lang string,
	warnIgnoreFile string,
	errIgnoreFile string) (bool, *ReturnError) {

	c := checker.GetChecks(includeChecks)
	c.Localizer = *localizations.New(lang, "en")

	errs, err := getBreakingChanges(c, diffReport, operationsSources, warnIgnoreFile, errIgnoreFile)
	if err != nil {
		return false, err
	}

	switch format {
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
		return false, getErrUnsupportedBreakingChangesFormat(format)
	}

	countErrs := len(errs) - errs.CountWarnings()

	diffEmpty := countErrs == 0
	if failOnWarns {
		diffEmpty = len(errs) == 0
	}

	return diffEmpty, nil
}

func getBreakingChanges(c checker.BackwardCompatibilityCheckConfig,
	diffReport *diff.Diff,
	operationsSources *diff.OperationsSourcesMap,
	warnIgnoreFile string,
	errIgnoreFile string) (checker.BackwardCompatibilityErrors, *ReturnError) {

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
