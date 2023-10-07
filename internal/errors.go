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
	pathStr := source.Path
	if !source.Stdin {
		pathStr = fmt.Sprintf("%q", source.Path)
	}

	return &ReturnError{
		error: fmt.Errorf("failed to load %s spec from %s with %v", what, pathStr, err),
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

func getErrFailedToFlattenSpec(what string, source load.Source, err error) *ReturnError {
	pathStr := source.Path
	if !source.Stdin {
		pathStr = fmt.Sprintf("%q", source.Path)
	}

	return &ReturnError{
		error: fmt.Errorf("failed to flatten %s spec from %s with %v", what, pathStr, err),
		Code:  102,
	}
}
