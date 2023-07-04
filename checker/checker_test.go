package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
)

func TestIsEmpty_EmptyIncludeWarns(t *testing.T) {
	bcErrors := checker.IBackwardCompatibilityErrors{}
	require.False(t, bcErrors.HasLevelOrHigher(checker.WARN))
}

func TestIsEmpty_EmptyExcludeWarns(t *testing.T) {
	bcErrors := checker.IBackwardCompatibilityErrors{}
	require.False(t, bcErrors.HasLevelOrHigher(checker.ERR))
}

func TestIsEmpty_OneErrIncludeWarns(t *testing.T) {
	bcErrors := checker.IBackwardCompatibilityErrors{
		checker.BackwardCompatibilityError{Level: checker.ERR},
	}
	require.True(t, bcErrors.HasLevelOrHigher(checker.WARN))
}

func TestIsEmpty_OneErrExcludeWarns(t *testing.T) {
	bcErrors := checker.IBackwardCompatibilityErrors{
		checker.BackwardCompatibilityError{Level: checker.ERR},
	}
	require.True(t, bcErrors.HasLevelOrHigher(checker.ERR))
}

func TestIsEmpty_OneWarnIncludeWarns(t *testing.T) {
	bcErrors := checker.IBackwardCompatibilityErrors{
		checker.BackwardCompatibilityError{Level: checker.WARN},
	}
	require.True(t, bcErrors.HasLevelOrHigher(checker.WARN))
}

func TestIsEmpty_OneWarnExcludeWarns(t *testing.T) {
	bcErrors := checker.IBackwardCompatibilityErrors{
		checker.BackwardCompatibilityError{Level: checker.WARN},
	}
	require.False(t, bcErrors.HasLevelOrHigher(checker.ERR))
}
