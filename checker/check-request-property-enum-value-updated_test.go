package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: removing request property enum values
func TestRequestPropertyEnumValueRemovedCheck(t *testing.T) {
	s1, err := open("../data/checker/request_property_enum_value_updated_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_enum_value_updated_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/pets"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["category"].Value.Enum = []interface{}{"dog", "cat"}

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyEnumValueUpdatedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)

	require.Equal(t, checker.ApiChange{
		Id:          "request-property-enum-value-removed",
		Level:       checker.ERR,
		Text:        "removed the enum value 'bird' of the request property 'category'",
		Operation:   "POST",
		OperationId: "updatePet",
		Path:        "/pets",
		Source:      "../data/checker/request_property_enum_value_updated_base.yaml",
	}, errs[0])
}

// CL: adding request property enum values
func TestRequestPropertyEnumValueAddedCheck(t *testing.T) {
	s1, err := open("../data/checker/request_property_enum_value_updated_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_enum_value_updated_base.yaml")
	require.NoError(t, err)

	s1.Spec.Paths["/pets"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["category"].Value.Enum = []interface{}{"dog", "cat"}

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyEnumValueUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)

	require.Equal(t, checker.ApiChange{
		Id:          "request-property-enum-value-added",
		Level:       checker.INFO,
		Text:        "added the new 'bird' enum value to the request property 'category'",
		Operation:   "POST",
		OperationId: "updatePet",
		Path:        "/pets",
		Source:      "../data/checker/request_property_enum_value_updated_base.yaml",
	}, errs[0])
}
