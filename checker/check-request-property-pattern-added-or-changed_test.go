package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: Changing of requ
func TestRequestPropertyPatternChanged(t *testing.T) {
	s1, err := open("../data/checker/request_property_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_pattern_added_or_changed_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyPatternAddedOrChangedCheck), d, osm, checker.WARN)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-property-pattern-changed",
		Text:        "changed the pattern for the request property 'name' from '^\\w+$' to to '^[\\w\\s]+$'",
		Comment:     "This is a warning because it is difficult to automatically analyze if the new pattern is a superset of the previous pattern(e.g. changed from '[0-9]+' to '[0-9]*')",
		Level:       checker.WARN,
		Operation:   "POST",
		Path:        "/test",
		Source:      "../data/checker/request_property_pattern_added_or_changed_revision.yaml",
		OperationId: "",
	}, errs[0])
}

// CL: Pattern of request property 'name' added
func TestRequestPropertyPatternAdded(t *testing.T) {
	s1, err := open("../data/checker/request_property_pattern_added_or_changed_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyPatternAddedOrChangedCheck), d, osm, checker.WARN)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-property-pattern-added",
		Text:        "pattern of request property 'name' added: '^[\\w\\s]+$'",
		Comment:     "This change may impact clients using the API. Please review the pattern change.",
		Level:       checker.WARN,
		Operation:   "POST",
		Path:        "/test",
		Source:      "../data/checker/request_property_pattern_added_or_changed_base.yaml",
		OperationId: "",
	}, errs[0])
}
