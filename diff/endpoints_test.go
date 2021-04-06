package diff_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func TestEndpointsSort(t *testing.T) {
	endpoints := diff.Endpoints{
		{
			Method: "GET",
			Path:   "/b",
		},
		{
			Method: "GET",
			Path:   "/a",
		},
	}

	sort.Sort(endpoints)
	require.Equal(t, "/a", endpoints[0].Path)
}

func TestEndpointsSort_Methods(t *testing.T) {
	endpoints := diff.Endpoints{
		{
			Method: "POST",
			Path:   "/a",
		},
		{
			Method: "OPTIONS",
			Path:   "/a",
		},
	}

	sort.Sort(endpoints)
	require.Equal(t, "OPTIONS", endpoints[0].Method)
}
