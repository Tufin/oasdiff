package diff_test

import (
	"fmt"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/utils"
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
	require.Equal(t, utils.StringList{"sku"}, dd.PathsDiff.Modified["/api"].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["200"].ContentDiff.MediaTypeModified["application/json"].SchemaDiff.AllOfDiff.Modified["#/components/schemas/ProductDto"].PropertiesDiff.Added)
}

func TestOneOf_TwoRefs(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getXOfFile("two-refs-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getXOfFile("two-refs-revision.yaml"))
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	require.Equal(t, utils.StringList{"guard"}, dd.PathsDiff.Modified["/pets"].OperationsDiff.Modified["PATCH"].RequestBodyDiff.ContentDiff.MediaTypeModified["application/json"].SchemaDiff.OneOfDiff.Modified["#/components/schemas/Dog"].AllOfDiff.Modified["#2"].PropertiesDiff.Added)
}

func TestOneOf_ChangeBoth(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getXOfFile("two-refs-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getXOfFile("two-refs-both-changed-revision.yaml"))
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	require.Equal(t, utils.StringList{"miao"}, dd.PathsDiff.Modified["/pets"].OperationsDiff.Modified["PATCH"].RequestBodyDiff.ContentDiff.MediaTypeModified["application/json"].SchemaDiff.OneOfDiff.Modified["#/components/schemas/Cat"].AllOfDiff.Modified["#2"].PropertiesDiff.Added)
	require.Equal(t, utils.StringList{"guard"}, dd.PathsDiff.Modified["/pets"].OperationsDiff.Modified["PATCH"].RequestBodyDiff.ContentDiff.MediaTypeModified["application/json"].SchemaDiff.OneOfDiff.Modified["#/components/schemas/Dog"].AllOfDiff.Modified["#2"].PropertiesDiff.Added)
}

func TestOneOf_TwoInlineDuplicate(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getXOfFile("two-inline-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getXOfFile("two-inline-revision-duplicate.yaml"))
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	require.Equal(t, "name2", dd.PathsDiff.Modified["/api"].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["200"].ContentDiff.MediaTypeModified["application/json"].SchemaDiff.OneOfDiff.Modified["#1"].PropertiesDiff.Added[0])
	require.Equal(t, "name1", dd.PathsDiff.Modified["/api"].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["200"].ContentDiff.MediaTypeModified["application/json"].SchemaDiff.OneOfDiff.Modified["#1"].PropertiesDiff.Deleted[0])
}

func TestOneOf_TwoInlineOneModified(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getXOfFile("two-inline-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getXOfFile("two-inline-revision-one-modified.yaml"))
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	require.Equal(t, "name4", dd.PathsDiff.Modified["/api"].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["200"].ContentDiff.MediaTypeModified["application/json"].SchemaDiff.OneOfDiff.Modified["#1"].PropertiesDiff.Added[0])
	require.Equal(t, "name1", dd.PathsDiff.Modified["/api"].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["200"].ContentDiff.MediaTypeModified["application/json"].SchemaDiff.OneOfDiff.Modified["#1"].PropertiesDiff.Deleted[0])
}

func TestOneOf_MultiRefs(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getXOfFile("multi-refs-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getXOfFile("multi-refs-revision.yaml"))
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	require.Equal(t, "bark", dd.PathsDiff.Modified["/pets"].OperationsDiff.Modified["GET"].RequestBodyDiff.ContentDiff.MediaTypeModified["application/json"].SchemaDiff.OneOfDiff.Modified["#/components/schemas/Dog"].PropertiesDiff.Added[0])
	require.Equal(t, "name", dd.PathsDiff.Modified["/pets"].OperationsDiff.Modified["GET"].RequestBodyDiff.ContentDiff.MediaTypeModified["application/json"].SchemaDiff.OneOfDiff.Modified["#/components/schemas/Dog"].PropertiesDiff.Deleted[0])
}

func TestAnyOf_IncludeDescriptions(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getXOfFile("anyof-base-openapi.yml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getXOfFile("anyof-rev-openapi.yml"))
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	anyOfDiff := dd.PathsDiff.Modified["/test"].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["200"].ContentDiff.MediaTypeModified["application/json"].SchemaDiff.AnyOfDiff
	require.ElementsMatch(t, []string{"RevisionSchema[0]", "RevisionSchema[2]"}, anyOfDiff.Added)
	require.ElementsMatch(t, []string{"BaseSchema[0]"}, anyOfDiff.Deleted)
	require.Empty(t, anyOfDiff.Modified)
}

func TestAnyOf_ExcludeDescriptions(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getXOfFile("anyof-base-openapi.yml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getXOfFile("anyof-rev-openapi.yml"))
	require.NoError(t, err)

	dd, err := diff.Get(diff.NewConfig().WithExcludeElements([]string{diff.ExcludeDescriptionOption}), s1, s2)
	require.NoError(t, err)
	anyOfDiff := dd.PathsDiff.Modified["/test"].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["200"].ContentDiff.MediaTypeModified["application/json"].SchemaDiff.AnyOfDiff
	require.ElementsMatch(t, []string{"RevisionSchema[2]"}, anyOfDiff.Added)
	require.ElementsMatch(t, []string{}, anyOfDiff.Deleted)
	require.Empty(t, anyOfDiff.Modified)
}
