package internal

import (
	"fmt"

	"github.com/tufin/oasdiff/load"
)

type ReturnError struct {
	error
	Code int
}

func getErrInvalidFlags(err error) *ReturnError {
	return &ReturnError{
		error: err,
		Code:  101,
	}
}

func getErrFailedToLoadSpec(what string, source load.Source, err error) *ReturnError {
	return &ReturnError{
		error: fmt.Errorf("failed to load %s spec from %s with %v", what, source.Out(), err),
		Code:  102,
	}
}

func getErrFailedToLoadSpecs(what string, path string, err error) *ReturnError {
	return &ReturnError{
		error: fmt.Errorf("failed to load %s specs from glob %q with %v", what, path, err),
		Code:  103,
	}
}

func getErrDiffFailed(err error) *ReturnError {
	return &ReturnError{
		error: fmt.Errorf("diff failed with %v", err),
		Code:  104,
	}
}

func getErrFailedPrint(what string, err error) *ReturnError {
	return &ReturnError{
		error: fmt.Errorf("failed to print %q with %v", what, err),
		Code:  105,
	}
}

func getErrUnsupportedFormat(format string) *ReturnError {
	return &ReturnError{
		error: fmt.Errorf("unsupported format %q", format),
		Code:  108,
	}
}

func getErrUnsupportedDiffFormat(format string) *ReturnError {
	return &ReturnError{
		error: fmt.Errorf("format %q is not supported by \"diff\"", format),
		Code:  109,
	}
}

func getErrUnsupportedSummaryFormat(format string) *ReturnError {
	return &ReturnError{
		error: fmt.Errorf("format %q is not supported by \"summary\"", format),
		Code:  110,
	}
}

func getErrUnsupportedChangelogFormat(format string) *ReturnError {
	return &ReturnError{
		error: fmt.Errorf("format %q is not supported by \"changelog\"", format),
		Code:  111,
	}
}

func getErrUnsupportedBreakingChangesFormat(format string) *ReturnError {
	return &ReturnError{
		error: fmt.Errorf("format %q is not supported by \"breaking\"", format),
		Code:  112,
	}
}

func getErrUnsupportedChecksFormat(format string) *ReturnError {
	return &ReturnError{
		error: fmt.Errorf("format %q is not supported with \"checks\"", format),
		Code:  113,
	}
}

func getErrCantProcessIgnoreFile(what string, err error) *ReturnError {
	return &ReturnError{
		error: fmt.Errorf("can't process %s ignore file %v", what, err),
		Code:  121,
	}
}

func getErrFailedToFlattenSpec(what string, source load.Source, err error) *ReturnError {
	return &ReturnError{
		error: fmt.Errorf("failed to flatten %s spec from %s with %v", what, source.Out(), err),
		Code:  102,
	}
}
