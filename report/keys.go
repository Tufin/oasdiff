package report

import (
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/utils"
)

type DiffT interface {
	*diff.ExampleDiff |
		*diff.ServerDiff |
		*diff.ParameterDiff |
		*diff.VariableDiff |
		*diff.SchemaDiff |
		*diff.ResponseDiff |
		*diff.MediaTypeDiff |
		*diff.HeaderDiff |
		diff.SecurityScopesDiff |
		*diff.StringsDiff
}

func getKeys[diff DiffT](m map[string]diff) utils.StringList {
	keys := make(utils.StringList, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys.Sort()
}
