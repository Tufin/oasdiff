package formatters_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/formatters"
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

func TestChanges_Group(t *testing.T) {
	require.Contains(t, formatters.GroupChanges(changes, checker.NewDefaultLocalizer()), formatters.Endpoint{Path: "/test", Operation: "GET"})
}
