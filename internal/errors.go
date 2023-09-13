package internal

import (
	"fmt"
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

func getErrFailedToLoadSpec(what string, path string, err error) *ReturnError {
	return &ReturnError{
		error: fmt.Errorf("failed to load %s spec from %q with %v", what, path, err),
		Code:  102,
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

func getErrFailedGenerateHTML(err error) *ReturnError {
	return &ReturnError{
		error: fmt.Errorf("failed to generate HTML diff report with %v", err),
		Code:  107,
	}
}

func getErrUnsupportedFormat(format string) *ReturnError {
	return &ReturnError{
		error: fmt.Errorf("unsupported format %q", format),
		Code:  108,
	}
}

func getErrCantProcessIgnoreFile(what string, err error) *ReturnError {
	return &ReturnError{
		error: fmt.Errorf("can't process %s ignore file %v", what, err),
		Code:  121,
	}
}

func getErrFailedToFlattenSpec(what string, path string, err error) *ReturnError {
	return &ReturnError{
		error: fmt.Errorf("failed to flatten %s spec from %q with %v", what, path, err),
		Code:  102,
	}
}
