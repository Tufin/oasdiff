package diff

import "github.com/tufin/oasdiff/utils"

// StringsDiff describes the changes between a pair of lists of strings
type StringsDiff struct {
	Added   utils.StringList `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted utils.StringList `json:"deleted,omitempty" yaml:"deleted,omitempty"`
}

func newStringsDiff() *StringsDiff {
	return &StringsDiff{
		Added:   utils.StringList{},
		Deleted: utils.StringList{},
	}
}

// Empty indicates whether a change was found in this element
func (stringsDiff *StringsDiff) Empty() bool {
	if stringsDiff == nil {
		return true
	}

	return len(stringsDiff.Added) == 0 &&
		len(stringsDiff.Deleted) == 0
}

func getStringsDiff(strings1, strings2 utils.StringList) *StringsDiff {
	diff := getStringsDiffInternal(strings1, strings2)

	if diff.Empty() {
		return nil
	}

	return diff
}

func getStringsDiffInternal(strings1, strings2 utils.StringList) *StringsDiff {
	result := newStringsDiff()

	s1 := strings1.ToStringSet()
	s2 := strings2.ToStringSet()

	result.Added = s2.Minus(s1).ToStringList()
	result.Deleted = s1.Minus(s2).ToStringList()

	return result
}
