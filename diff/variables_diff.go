package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/utils"
)

// VariablesDiff describes the changes between a pair of sets of server variable objects: https://swagger.io/specification/#server-variable-object
type VariablesDiff struct {
	Added    utils.StringList  `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted  utils.StringList  `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	Modified ModifiedVariables `json:"modified,omitempty" yaml:"modified,omitempty"`
}

// ModifiedVariables is map of variable names to their respective diffs
type ModifiedVariables map[string]*VariableDiff

// Empty indicates whether a change was found in this element
func (diff *VariablesDiff) Empty() bool {
	if diff == nil {
		return true
	}

	return len(diff.Added) == 0 &&
		len(diff.Deleted) == 0 &&
		len(diff.Modified) == 0
}

func newVariablesDiff() *VariablesDiff {
	return &VariablesDiff{
		Added:    utils.StringList{},
		Deleted:  utils.StringList{},
		Modified: ModifiedVariables{},
	}
}

func getVariablesDiff(config *Config, state *state, variables1, variables2 map[string]*openapi3.ServerVariable) *VariablesDiff {
	diff := getVariablesDiffInternal(config, state, variables1, variables2)

	if diff.Empty() {
		return nil
	}

	return diff
}

func getVariablesDiffInternal(config *Config, state *state, variables1, variables2 map[string]*openapi3.ServerVariable) *VariablesDiff {
	result := newVariablesDiff()

	for name1, var1 := range variables1 {
		var2, ok := variables2[name1]
		if !ok {
			result.Deleted = append(result.Deleted, name1)
			continue
		}

		if diff := getVariableDiff(config, state, var1, var2); !diff.Empty() {
			result.Modified[name1] = diff
		}
	}

	for name2 := range variables2 {
		_, ok := variables1[name2]
		if !ok {
			result.Added = append(result.Added, name2)
		}
	}

	return result
}
