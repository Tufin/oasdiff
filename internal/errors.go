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
		fmt.Errorf("failed to load %s spec from %s: %w", what, source.Out(), err),
		102,
	)
}

func getErrFailedToLoadSpecs(what string, path string, err error) *ReturnError {
	return getError(
		fmt.Errorf("failed to load %s specs from glob %q: %w", what, path, err),
		103,
	)
}

func getErrDiffFailed(err error) *ReturnError {
	return getError(
		fmt.Errorf("diff failed: %w", err),
		104,
	)
}

func getErrFailedPrint(what string, err error) *ReturnError {
	return getError(
		fmt.Errorf("failed to print %q: %w", what, err),
		105,
	)
}

func getErrUnsupportedFormat(format, cmd string, code int) *ReturnError {
	return getError(
		fmt.Errorf("format %q is not supported by %q", format, cmd),
		code,
	)
}

func getErrUnsupportedDiffFormat(format string) *ReturnError {
	return getErrUnsupportedFormat(format, "diff", 109)
}

func getErrUnsupportedSummaryFormat(format string) *ReturnError {
	return getErrUnsupportedFormat(format, "summary", 110)
}

func getErrUnsupportedChangelogFormat(format string) *ReturnError {
	return getErrUnsupportedFormat(format, "changelog", 111)
}

func getErrUnsupportedChecksFormat(format string) *ReturnError {
	return getErrUnsupportedFormat(format, "checks", 113)
}

func getErrInvalidColorMode(err error) *ReturnError {
	return getError(
		err,
		114,
	)
}

func getErrCantProcessIgnoreFile(what string, err error) *ReturnError {
	return getError(
		fmt.Errorf("can't process %s ignore file: %w", what, err),
		121,
	)
}

func getError(err error, code int) *ReturnError {
	return &ReturnError{err, code}
}
