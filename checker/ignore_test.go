package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

func TestIgnore(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 3)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Equal(t, 6, len(errs))

	errs, err = checker.ProcessIgnoredBackwardCompatibilityErrors(checker.ERR, errs, "../data/ignore-err-example.txt")
	require.NoError(t, err)
	require.Equal(t, 5, len(errs))
}

func TestIgnoreSubpath(t *testing.T) {
	s1 := l(t, 6)
	s2 := l(t, 7)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Equal(t, 3, len(errs))

	errs, err = checker.ProcessIgnoredBackwardCompatibilityErrors(checker.ERR, errs, "../data/ignore-err-example-2.txt")
	require.NoError(t, err)
	require.Equal(t, 0, len(errs))
}

func TestIgnoreOnlyIncludedSubpaths(t *testing.T) {
	s1 := l(t, 8)
	s2 := l(t, 7)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Equal(t, 2, len(errs)) // detect new and newest were deleted

	errs, err = checker.ProcessIgnoredBackwardCompatibilityErrors(checker.ERR, errs, "../data/ignore-err-example-3.txt")
	require.NoError(t, err)
	require.Equal(t, 1, len(errs))
	require.IsType(t, checker.ApiChange{}, errs[0])
	e0 := errs[0].(checker.ApiChange)
	require.Contains(t, e0.Path, "/resource/new") //see that new breaking change was kept even though it is a substring of newest
}
