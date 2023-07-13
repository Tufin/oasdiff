package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: Adding an enum value to a response property
func TestResponsePropertyEnumValueAdded(t *testing.T) {
	s1, err := open("../data/checker/response_property_enum_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_enum_added_base.yaml")
	require.NoError(t, err)

	s2.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["typeEnum"].Value.Enum = []interface{}{"Test"}

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyEnumValueAddedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "response-property-enum-value-added",
		Text:        "added the new 'Test' enum value the 'data/typeEnum' response property for the response status '200'",
		Comment:     "Adding new enum values to response could be unexpected for clients, use x-extensible-enum instead.",
		Level:       checker.WARN,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_property_enum_added_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: Adding an enum value to a response write-only property
func TestResponseWriteOnlyPropertyEnumValueAdded(t *testing.T) {
	s1, err := open("../data/checker/response_property_enum_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_enum_added_base.yaml")
	require.NoError(t, err)

	s2.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["writeOnlyEnum"].Value.Enum = []interface{}{"Test"}

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyEnumValueAddedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "response-write-only-property-enum-value-added",
		Text:        "added the new 'Test' enum value the 'data/writeOnlyEnum' response write-only property for the response status '200'",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_property_enum_added_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}
