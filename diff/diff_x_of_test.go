package diff_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func TestAllOf_SingleRef(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile("../data/x-of/single-ref-base.yaml")
	require.NoError(t, err)

	s2, err := loader.LoadFromFile("../data/x-of/single-ref-revision.yaml")
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	require.Equal(t, diff.StringList{"sku"}, dd.PathsDiff.Modified["/api"].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["200"].ContentDiff.MediaTypeModified["application/json"].SchemaDiff.AllOfDiff.Modified.PropertiesDiff.Added)
}

func TestOneOf_TwoRefs(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile("../data/x-of/two-refs-base.yaml")
	require.NoError(t, err)

	s2, err := loader.LoadFromFile("../data/x-of/two-refs-revision.yaml")
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	require.Equal(t, diff.StringList{"guard"}, dd.PathsDiff.Modified["/pets"].OperationsDiff.Modified["PATCH"].RequestBodyDiff.ContentDiff.MediaTypeModified["application/json"].SchemaDiff.OneOfDiff.Modified.AllOfDiff.Modified.PropertiesDiff.Added)

}
