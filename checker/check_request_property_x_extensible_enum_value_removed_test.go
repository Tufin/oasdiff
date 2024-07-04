package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// BC: Deleting a value from an x-extensible-enum property is breaking
func TestRequestPropertyXExtensibleEnumValueRemovedCheck(t *testing.T) {
	s1, err := open("../data/checker/request_property_extensible_enum_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/request_property_extensible_enum_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.RequestPropertyXExtensibleEnumValueRemovedCheck), d, osm)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestPropertyXExtensibleEnumValueRemovedId, errs[0].GetId())
}
