package diff

import "reflect"

type EnumDiff struct {
	Added   []interface{} `json:"added,omitempty"`
	Deleted []interface{} `json:"deleted,omitempty"`
}

func newEnumDiff() *EnumDiff {
	return &EnumDiff{
		Added:   []interface{}{},
		Deleted: []interface{}{},
	}
}

func (enumDiff *EnumDiff) empty() bool {
	return len(enumDiff.Added) == 0 &&
		len(enumDiff.Deleted) == 0
}

func getEnumDiff(enum1 []interface{}, enum2 []interface{}) *EnumDiff {

	if enum1 == nil && enum2 == nil {
		return nil
	}

	diff := newEnumDiff()

	for _, v1 := range enum1 {
		if !findValue(v1, enum2) {
			diff.Deleted = append(diff.Deleted, v1)
		}
	}

	for _, v2 := range enum2 {
		if !findValue(v2, enum1) {
			diff.Added = append(diff.Added, v2)
		}
	}

	if diff.empty() {
		return nil
	}

	return diff
}

func findValue(value interface{}, enum []interface{}) bool {
	for _, other := range enum {
		if reflect.DeepEqual(value, other) {
			return true
		}
	}
	return false
}
