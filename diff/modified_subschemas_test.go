package diff_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func TestModifiedSubschemaString_TitlesSame(t *testing.T) {

	ms := diff.ModifiedSubschema{
		Base: diff.Subschema{
			Title: "Title 1",
		},
		Revision: diff.Subschema{
			Title: "Title 1",
		},
	}

	require.Equal(t, "subschema #1: Title 1", ms.String())
}

func TestModifiedSubschemaString_TitlesDiff(t *testing.T) {

	ms := diff.ModifiedSubschema{
		Base: diff.Subschema{
			Title: "Title 1",
		},
		Revision: diff.Subschema{
			Title: "Title 2",
		},
	}

	require.Equal(t, "subschema #1: Title 1 -> subschema #1: Title 2", ms.String())
}

func TestModifiedSubschemaString_ComponentsSame(t *testing.T) {

	ms := diff.ModifiedSubschema{
		Base: diff.Subschema{
			Component: "component1",
		},
		Revision: diff.Subschema{
			Component: "component1",
		},
	}

	require.Equal(t, "#/components/schemas/component1", ms.String())
}

func TestModifiedSubschemaString_ComponentsDiff(t *testing.T) {

	ms := diff.ModifiedSubschema{
		Base: diff.Subschema{
			Component: "component1",
		},
		Revision: diff.Subschema{
			Component: "component2",
		},
	}

	require.Equal(t, "#/components/schemas/component1 -> #/components/schemas/component2", ms.String())
}

func TestSubschemaString(t *testing.T) {
	ms := diff.Subschemas{
		diff.Subschema{
			Index: 0,
			Title: "Title 0",
		},
		diff.Subschema{
			Index: 1,
			Title: "Title 1",
		},
	}
	require.Equal(t, "subschema #1: Title 0, subschema #2: Title 1", ms.String())
}

func TestSubschemaString_Empty(t *testing.T) {
	ms := diff.Subschemas{
		diff.Subschema{
			Index: 0,
		},
		diff.Subschema{
			Index: 1,
		},
	}
	require.Equal(t, "subschema #1, subschema #2", ms.String())
}
