package internal

import (
	"fmt"
)

type ReturnError struct {
	Err  error
	Code int
}

func getErrNone() *ReturnError {
	return &ReturnError{
		Err:  nil, // don't print an error in this case
		Code: 0,
	}
}

func getErrParseFlags() *ReturnError {
	return &ReturnError{
		Err:  nil, // don't print an error in this case
		Code: 100,
	}
}

func getErrInvalidFlags(err error) *ReturnError {
	return &ReturnError{
		Err:  err,
		Code: 101,
	}
}

func getErrFailedToLoadSpec(what string, path string, err error) *ReturnError {
	return &ReturnError{
		Err:  fmt.Errorf("failed to load %s spec from %q with %v", what, path, err),
		Code: 102,
	}
}

func getErrDiffFailed(err error) *ReturnError {
	return &ReturnError{
		Err:  fmt.Errorf("diff failed with %v", err),
		Code: 104,
	}
}

func getErrFailedPrint(what string, err error) *ReturnError {
	return &ReturnError{
		Err:  fmt.Errorf("failed to print %q with %v", what, err),
		Code: 105,
	}
}

func getErrFailedGenerateHTML(err error) *ReturnError {
	return &ReturnError{
		Err:  fmt.Errorf("failed to generate HTML diff report with %v", err),
		Code: 107,
	}
}

func getErrUnsupportedDiffFormat(format string) *ReturnError {
	return &ReturnError{
		Err:  fmt.Errorf("unsupported format %q", format),
		Code: 108,
	}
}

func getErrUnsupportedBreakingChangesFormat(format string) *ReturnError {
	return &ReturnError{
		Err:  fmt.Errorf("format %q is not supported with \"-check-breaking\"", format),
		Code: 109,
	}
}
func getErrCantProcessIgnoreFile(what string, err error) *ReturnError {
	return &ReturnError{
		Err:  fmt.Errorf("can't process %s ignore file %v", what, err),
		Code: 121,
	}
}

func getErrLintFailed() *ReturnError {
	return &ReturnError{
		Err:  nil,
		Code: 130,
	}
}
