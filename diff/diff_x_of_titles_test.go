package diff_test

import (
	"fmt"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
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
	require.Len(t, anyOfDiff.Added, 1)
	require.Len(t, anyOfDiff.Deleted, 1)
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
	require.Len(t, anyOfDiff.Added, 1)
	require.Empty(t, anyOfDiff.Deleted)
	require.Len(t, anyOfDiff.Modified, 1)
}
