package checker_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
)

var changes = checker.Changes{
	checker.ApiChange{
		Id:        "api-deleted",
		Text:      "API deleted",
		Comment:   "",
		Level:     checker.ERR,
		Operation: "GET",
		Path:      "/test",
	},
	checker.ApiChange{
		Id:        "api-added",
		Text:      "API added",
		Comment:   "",
		Level:     checker.INFO,
		Operation: "GET",
		Path:      "/test",
	},
}

func TestChanges_Sort(t *testing.T) {
	sort.Sort(changes)
	require.False(t, changes[1].IsBreaking())
}

func TestChanges_Count(t *testing.T) {
	require.Equal(t, map[checker.Level]int{checker.INFO: 1, checker.ERR: 1}, changes.GetLevelCount())
}

func TestChanges_Group(t *testing.T) {
	require.Contains(t, changes.Group().APIChanges, checker.Endpoint{Path: "/test", Operation: "GET"})
}
