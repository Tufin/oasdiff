package diff_test

import (
	"fmt"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func getXOfFile(file string) string {
	return fmt.Sprintf("../data/x-of/%s", file)
}

func TestAllOf_SingleRef(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getXOfFile("single-ref-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getXOfFile("single-ref-revision.yaml"))
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	require.Equal(t, diff.StringList{"sku"}, dd.PathsDiff.Modified["/api"].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["200"].ContentDiff.MediaTypeModified["application/json"].SchemaDiff.AllOfDiff.Modified["#/components/schemas/ProductDto"].PropertiesDiff.Added)
}

func TestOneOf_TwoRefs(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getXOfFile("two-refs-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getXOfFile("two-refs-revision.yaml"))
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	require.Equal(t, diff.StringList{"guard"}, dd.PathsDiff.Modified["/pets"].OperationsDiff.Modified["PATCH"].RequestBodyDiff.ContentDiff.MediaTypeModified["application/json"].SchemaDiff.OneOfDiff.Modified["#/components/schemas/Dog"].AllOfDiff.Modified["#2"].PropertiesDiff.Added)
}

func TestOneOf_ChangeBoth(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getXOfFile("two-refs-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getXOfFile("two-refs-both-changed-revision.yaml"))
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	require.Equal(t, diff.StringList{"miao"}, dd.PathsDiff.Modified["/pets"].OperationsDiff.Modified["PATCH"].RequestBodyDiff.ContentDiff.MediaTypeModified["application/json"].SchemaDiff.OneOfDiff.Modified["#/components/schemas/Cat"].AllOfDiff.Modified["#2"].PropertiesDiff.Added)
	require.Equal(t, diff.StringList{"guard"}, dd.PathsDiff.Modified["/pets"].OperationsDiff.Modified["PATCH"].RequestBodyDiff.ContentDiff.MediaTypeModified["application/json"].SchemaDiff.OneOfDiff.Modified["#/components/schemas/Dog"].AllOfDiff.Modified["#2"].PropertiesDiff.Added)
}

func TestOneOf_TwoInline(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getXOfFile("two-inline-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getXOfFile("two-inline-revision.yaml"))
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	require.Equal(t, 1, dd.PathsDiff.Modified["/api"].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["200"].ContentDiff.MediaTypeModified["application/json"].SchemaDiff.OneOfDiff.Deleted)
}
