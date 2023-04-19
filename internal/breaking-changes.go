package internal

import (
	"fmt"

	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/checker/localizations"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/utils"
)

func HandleBreakingChanges(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap,
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
		if err := PrintYAML(errs); err != nil {
			return false, GetErrFailedToPrint("breaking changes YAML", err)
		}
	case FormatJSON:
		if err := PrintJSON(errs); err != nil {
			return false, GetErrFailedToPrint("breaking changes JSON", err)
		}
	case FormatText:
		if len(errs) > 0 {
			fmt.Printf(c.Localizer.Get("messages.total-errors"), len(errs))
		}

		for _, bcerr := range errs {
			fmt.Printf("%s\n\n", bcerr.PrettyErrorText(c.Localizer))
		}
	default:
		return false, GetErrUnsupportedBreakingChangesFormat(format)
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
			return nil, GetErrCantProcessIgnoreFile("warn", err)
		}
	}

	if errIgnoreFile != "" {
		var err error
		errs, err = checker.ProcessIgnoredBackwardCompatibilityErrors(checker.ERR, errs, errIgnoreFile)
		if err != nil {
			return nil, GetErrCantProcessIgnoreFile("err", err)
		}
	}

	return errs, nil
}
