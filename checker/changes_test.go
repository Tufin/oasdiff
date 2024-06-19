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
		Level:     checker.ERR,
		Operation: "GET",
		Path:      "/test",
	},
	checker.ApiChange{
		Id:        "api-added",
		Level:     checker.INFO,
		Operation: "GET",
		Path:      "/test",
	},
	checker.ComponentChange{
		Id:    "component-added",
		Level: checker.INFO,
	},
	checker.SecurityChange{
		Id:    "security-added",
		Level: checker.INFO,
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

func TestIsEmpty_EmptyIncludeWarns(t *testing.T) {
	bcErrors := checker.Changes{}
	require.False(t, bcErrors.HasLevelOrHigher(checker.WARN))
}

func TestIsEmpty_EmptyExcludeWarns(t *testing.T) {
	bcErrors := checker.Changes{}
	require.False(t, bcErrors.HasLevelOrHigher(checker.ERR))
}

func TestIsEmpty_OneErrIncludeWarns(t *testing.T) {
	bcErrors := checker.Changes{
		checker.ApiChange{Level: checker.ERR},
	}
	require.True(t, bcErrors.HasLevelOrHigher(checker.WARN))
}

func TestIsEmpty_OneErrExcludeWarns(t *testing.T) {
	bcErrors := checker.Changes{
		checker.ApiChange{Level: checker.ERR},
	}
	require.True(t, bcErrors.HasLevelOrHigher(checker.ERR))
}

func TestIsEmpty_OneWarnIncludeWarns(t *testing.T) {
	bcErrors := checker.Changes{
		checker.ApiChange{Level: checker.WARN},
	}
	require.True(t, bcErrors.HasLevelOrHigher(checker.WARN))
}

func TestIsEmpty_OneWarnExcludeWarns(t *testing.T) {
	bcErrors := checker.Changes{
		checker.ApiChange{Level: checker.WARN},
	}
	require.False(t, bcErrors.HasLevelOrHigher(checker.ERR))
}
