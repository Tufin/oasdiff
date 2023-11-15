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
	checker.ComponentChange{
		Id:      "component-added",
		Text:    "component added",
		Comment: "",
		Level:   checker.INFO,
	},
	checker.SecurityChange{
		Id:      "security-added",
		Text:    "security added",
		Comment: "",
		Level:   checker.INFO,
	},
}

func TestChanges_Sort(t *testing.T) {
	sort.Sort(changes)
}

func TestChanges_IsBreaking(t *testing.T) {
	for _, c := range changes {
		require.True(t, c.IsBreaking() == (c.GetLevel() != checker.INFO))
	}
}

func TestChanges_Count(t *testing.T) {
	lc := changes.GetLevelCount()
	require.Equal(t, 3, lc[checker.INFO])
	require.Equal(t, 0, lc[checker.WARN])
	require.Equal(t, 1, lc[checker.ERR])
}

func TestChanges_Group(t *testing.T) {
	require.Contains(t, changes.Group().APIChanges, checker.Endpoint{Path: "/test", Operation: "GET"})
}
