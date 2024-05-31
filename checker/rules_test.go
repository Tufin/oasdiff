package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
)

func TestGetOptionalRuleIds(t *testing.T) {
	require.Len(t, checker.GetOptionalRuleIds(), 7)
}
