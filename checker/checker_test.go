package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
)

func TestIsEmpty_EmptyIncludeWarns(t *testing.T) {
	bcErrors := checker.BackwardCompatibilityErrors{}
	require.True(t, bcErrors.IsEmpty(true))
}

func TestIsEmpty_EmptyExcludeWarns(t *testing.T) {
	bcErrors := checker.BackwardCompatibilityErrors{}
	require.True(t, bcErrors.IsEmpty(false))
}

func TestIsEmpty_OneErrIncludeWarns(t *testing.T) {
	bcErrors := checker.BackwardCompatibilityErrors{
		checker.BackwardCompatibilityError{Level: checker.ERR},
	}
	require.False(t, bcErrors.IsEmpty(true))
}

func TestIsEmpty_OneErrExcludeWarns(t *testing.T) {
	bcErrors := checker.BackwardCompatibilityErrors{
		checker.BackwardCompatibilityError{Level: checker.ERR},
	}
	require.False(t, bcErrors.IsEmpty(false))
}

func TestIsEmpty_OneWarnIncludeWarns(t *testing.T) {
	bcErrors := checker.BackwardCompatibilityErrors{
		checker.BackwardCompatibilityError{Level: checker.WARN},
	}
	require.False(t, bcErrors.IsEmpty(true))
}

func TestIsEmpty_OneWarnExcludeWarns(t *testing.T) {
	bcErrors := checker.BackwardCompatibilityErrors{
		checker.BackwardCompatibilityError{Level: checker.WARN},
	}
	require.True(t, bcErrors.IsEmpty(false))
}
