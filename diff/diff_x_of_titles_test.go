package diff_test

import (
	"fmt"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/utils"
)

func getXOfTitlesFile(file string) string {
	return fmt.Sprintf("../data/x-of-titles/%s", file)
}

func TestXOfTitles_Identical(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getXOfTitlesFile("spec1.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getXOfTitlesFile("spec1.yaml"))
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, dd)
}

func TestXOfTitles_TitleNameChanged(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getXOfTitlesFile("spec1.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getXOfTitlesFile("spec2.yaml"))
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, dd)
	anyOfDiff := dd.PathsDiff.Modified["/test"].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["200"].ContentDiff.MediaTypeModified["application/json"].SchemaDiff.AnyOfDiff
	require.ElementsMatch(t, utils.StringList{"RevisionSchema[1]:Title 3"}, anyOfDiff.Added)
	require.ElementsMatch(t, utils.StringList{"BaseSchema[1]:Title 2"}, anyOfDiff.Deleted)
	require.Empty(t, anyOfDiff.Modified)
}

func TestXOfTitles_TitlesModified(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getXOfTitlesFile("spec1.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getXOfTitlesFile("spec3.yaml"))
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, dd)
	anyOfDiff := dd.PathsDiff.Modified["/test"].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["200"].ContentDiff.MediaTypeModified["application/json"].SchemaDiff.AnyOfDiff
	require.Empty(t, anyOfDiff.Added)
	require.Empty(t, anyOfDiff.Deleted)
	require.Len(t, anyOfDiff.Modified, 2)
	require.Equal(t, diff.ValueDiff{From: "string", To: "boolean"}, *anyOfDiff.Modified["Title 1"].TypeDiff)
	require.Equal(t, diff.ValueDiff{From: "string", To: "number"}, *anyOfDiff.Modified["Title 2"].TypeDiff)
}

func TestXOfTitles_TitlesModifiedAndAdded(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getXOfTitlesFile("spec3.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getXOfTitlesFile("spec4.yaml"))
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, dd)
	anyOfDiff := dd.PathsDiff.Modified["/test"].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["200"].ContentDiff.MediaTypeModified["application/json"].SchemaDiff.AnyOfDiff
	require.ElementsMatch(t, utils.StringList{"RevisionSchema[2]:Title 3"}, anyOfDiff.Added)
	require.Empty(t, anyOfDiff.Deleted)
	require.Len(t, anyOfDiff.Modified, 1)
	require.Equal(t, diff.ValueDiff{From: "boolean", To: "string"}, *anyOfDiff.Modified["Title 1"].TypeDiff)
}

func TestXOfTitles_DuplicateTitles(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getXOfTitlesFile("spec3.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getXOfTitlesFile("spec5.yaml"))
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, dd)
	anyOfDiff := dd.PathsDiff.Modified["/test"].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["200"].ContentDiff.MediaTypeModified["application/json"].SchemaDiff.AnyOfDiff
	require.ElementsMatch(t, utils.StringList{"RevisionSchema[0]:Title 2", "RevisionSchema[2]:Title 3"}, anyOfDiff.Added)
	require.ElementsMatch(t, utils.StringList{"BaseSchema[0]:Title 1"}, anyOfDiff.Deleted)
	require.Empty(t, anyOfDiff.Modified, 0)
}
