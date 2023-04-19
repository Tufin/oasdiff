package internal

import "fmt"

type ReturnError struct {
	Err  error
	Code int
}

func GetErrInvalidFlags(err error) *ReturnError {
	return &ReturnError{
		Err:  err,
		Code: 101,
	}
}

func GetErrFailedToLoadSpec(what string, path string, err error) *ReturnError {
	return &ReturnError{
		Err:  fmt.Errorf("failed to load %s spec from %q with %v", what, path, err),
		Code: 102,
	}
}

func GetErrDiffFailed(err error) *ReturnError {
	return &ReturnError{
		Err:  fmt.Errorf("diff failed with %v", err),
		Code: 104,
	}
}

func GetErrFailedPrint(what string, err error) *ReturnError {
	return &ReturnError{
		Err:  fmt.Errorf("failed to print %q with %v", what, err),
		Code: 105,
	}
}

func GetErrFailedToPrint(what string, err error) *ReturnError {
	return &ReturnError{
		Err:  fmt.Errorf("failed to print %s with %v", what, err),
		Code: 106,
	}
}

func GetErrFailedGenerateHTML(err error) *ReturnError {
	return &ReturnError{
		Err:  fmt.Errorf("failed to generate HTML diff report with %v", err),
		Code: 107,
	}
}

func GetErrUnsupportedDiffFormat(format string) *ReturnError {
	return &ReturnError{
		Err:  fmt.Errorf("unsupported format %q", format),
		Code: 108,
	}
}

func GetErrUnsupportedBreakingChangesFormat(format string) *ReturnError {
	return &ReturnError{
		Err:  fmt.Errorf("format %q is not supported with \"-check-breaking\"", format),
		Code: 109,
	}
}
func GetErrCantProcessIgnoreFile(what string, err error) *ReturnError {
	return &ReturnError{
		Err:  fmt.Errorf("can't process %s ignore file %v", what, err),
		Code: 121,
	}
}
