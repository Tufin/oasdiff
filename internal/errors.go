package internal

import (
	"fmt"

	"github.com/tufin/oasdiff/load"
)

type ReturnError struct {
	error
	Code int
}

const generalExecutionErr = 100

func getErrInvalidFlags(err error) *ReturnError {
	return getError(
		err,
		101,
	)
}

func getErrFailedToLoadSpec(what string, source *load.Source, err error) *ReturnError {
	return getError(
		fmt.Errorf("failed to load %s spec from %s with %v", what, source.Out(), err),
		102,
	)
}

func getErrFailedToLoadSpecs(what string, path string, err error) *ReturnError {
	return getError(
		fmt.Errorf("failed to load %s specs from glob %q with %v", what, path, err),
		103,
	)
}

func getErrDiffFailed(err error) *ReturnError {
	return getError(
		fmt.Errorf("diff failed with %v", err),
		104,
	)
}

func getErrFailedPrint(what string, err error) *ReturnError {
	return getError(
		fmt.Errorf("failed to print %q with %v", what, err),
		105,
	)
}

func getErrUnsupportedDiffFormat(format string) *ReturnError {
	return getError(
		fmt.Errorf("format %q is not supported by \"diff\"", format),
		109,
	)
}

func getErrUnsupportedSummaryFormat(format string) *ReturnError {
	return getError(
		fmt.Errorf("format %q is not supported by \"summary\"", format),
		110,
	)
}

func getErrUnsupportedChangelogFormat(format string) *ReturnError {
	return getError(
		fmt.Errorf("format %q is not supported by \"changelog\"", format),
		111,
	)
}

func getErrUnsupportedBreakingChangesFormat(format string) *ReturnError {
	return getError(
		fmt.Errorf("format %q is not supported by \"breaking\"", format),
		112,
	)
}

func getErrUnsupportedChecksFormat(format string) *ReturnError {
	return getError(
		fmt.Errorf("format %q is not supported with \"checks\"", format),
		113,
	)
}

func getErrInvalidColorMode(err error) *ReturnError {
	return getError(
		err,
		114,
	)
}

func getErrCantProcessIgnoreFile(what string, err error) *ReturnError {
	return getError(
		fmt.Errorf("can't process %s ignore file %v", what, err),
		121,
	)
}

func getError(err error, code int) *ReturnError {
	return &ReturnError{err, code}
}
