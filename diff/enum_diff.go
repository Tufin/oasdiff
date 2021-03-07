package diff

import "reflect"

// EnumDiff is the diff between two enums
type EnumDiff struct {
	Added   EnumValues `json:"added,omitempty"`
	Deleted EnumValues `json:"deleted,omitempty"`
}

// EnumValues is a list of enum values
type EnumValues []interface{}

func newEnumDiff() *EnumDiff {
	return &EnumDiff{
		Added:   EnumValues{},
		Deleted: EnumValues{},
	}
}

func (enumDiff *EnumDiff) empty() bool {
	if enumDiff == nil {
		return true
	}
	
	return len(enumDiff.Added) == 0 &&
		len(enumDiff.Deleted) == 0
}

func getEnumDiff(enum1, enum2 EnumValues) *EnumDiff {

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

func findValue(value interface{}, enum EnumValues) bool {
	for _, other := range enum {
		if reflect.DeepEqual(value, other) {
			return true
		}
	}
	return false
}
