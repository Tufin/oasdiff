package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: changing request body default value
func TestRequestBodyDefaultValueChanged(t *testing.T) {
	s1, err := open("../data/checker/request_body_default_value_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_body_default_value_changed_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyDefaultValueChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestBodyDefaultValueChangedId,
		Text:        "the request body 'text/plain' default value changed from 'Default' to 'NewDefault'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/products",
		Source:      "../data/checker/request_body_default_value_changed_revision.yaml",
		OperationId: "createProduct",
	}, errs[0])
}

// CL: changing request property default value
func TestRequestPropertyDefaultValueChanged(t *testing.T) {
	s1, err := open("../data/checker/request_property_default_value_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_default_value_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/products"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["price"].Value.Default = 20.0

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyDefaultValueChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyDefaultValueChangedId,
		Text:        "the 'price' request property default value changed from '10.00' to '20.00'",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/products",
		Source:      "../data/checker/request_property_default_value_changed_base.yaml",
		OperationId: "createProduct",
	}, errs[0])
}

// CL: adding request body default value or request property default value
func TestRequestBodyDefaultValueAdded(t *testing.T) {
	s1, err := open("../data/checker/request_body_default_value_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_body_default_value_changed_base.yaml")
	require.NoError(t, err)

	s1.Spec.Paths["/products"].Post.RequestBody.Value.Content["text/plain"].Schema.Value.Default = nil
	s1.Spec.Paths["/products"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["price"].Value.Default = nil

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyDefaultValueChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 2)
	require.ElementsMatch(t, []checker.ApiChange{{
		Id:          checker.RequestBodyDefaultValueAddedId,
		Text:        "the request body 'text/plain' default value 'Default' was added",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/products",
		Source:      "../data/checker/request_body_default_value_changed_base.yaml",
		OperationId: "createProduct",
	}, {
		Id:          checker.RequestPropertyDefaultValueAddedId,
		Text:        "the 'price' request property default value '10.00' was added",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/products",
		Source:      "../data/checker/request_body_default_value_changed_base.yaml",
		OperationId: "createProduct",
	}}, errs)
}

// CL: removing request body default value or request property default value
func TestRequestBodyDefaultValueRemoving(t *testing.T) {
	s1, err := open("../data/checker/request_body_default_value_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_body_default_value_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/products"].Post.RequestBody.Value.Content["text/plain"].Schema.Value.Default = nil
	s2.Spec.Paths["/products"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["price"].Value.Default = nil

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyDefaultValueChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 2)
	require.ElementsMatch(t, []checker.ApiChange{{
		Id:          checker.RequestBodyDefaultValueRemovedId,
		Text:        "the request body 'text/plain' default value 'Default' was removed",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/products",
		Source:      "../data/checker/request_body_default_value_changed_base.yaml",
		OperationId: "createProduct",
	}, {
		Id:          checker.RequestPropertyDefaultValueRemovedId,
		Text:        "the 'price' request property default value '10.00' was removed",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/products",
		Source:      "../data/checker/request_body_default_value_changed_base.yaml",
		OperationId: "createProduct",
	}}, errs)
}
