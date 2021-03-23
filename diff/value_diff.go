package diff

import (
	"fmt"
	"reflect"
)

// ValueDiff describes the diff between two values
type ValueDiff struct {
	From interface{} `json:"from" yaml:"from"`
	To   interface{} `json:"to" yaml:"to"`
}

// Empty indicates whether a change was found in this element
func (diff *ValueDiff) Empty() bool {
	return diff == nil
}

func getValueDiff(value1, value2 interface{}) *ValueDiff {

	if reflect.DeepEqual(value1, value2) {
		return nil
	}

	return &ValueDiff{
		From: value1,
		To:   value2,
	}
}

func getFloat64RefDiff(valueRef1, valueRef2 *float64) *ValueDiff {
	return getValueDiff(derefFloat64(valueRef1), derefFloat64(valueRef2))
}

func getBoolRefDiff(valueRef1, valueRef2 *bool) *ValueDiff {
	return getValueDiff(derefBool(valueRef1), derefBool(valueRef2))
}

func getStringRefDiff(valueRef1, valueRef2 *string) *ValueDiff {
	return getValueDiff(derefString(valueRef1), derefString(valueRef2))
}

func derefString(ref *string) interface{} {
	if ref == nil {
		return nil
	}

	return *ref
}

func derefBool(ref *bool) interface{} {
	if ref == nil {
		return nil
	}

	return *ref
}

func derefFloat64(ref *float64) interface{} {
	if ref == nil {
		return nil
	}

	return *ref
}

func (diff *ValueDiff) patchStringCB(cb func(string)) error {
	if diff.Empty() {
		return nil
	}

	if diff.To == nil {
		return fmt.Errorf("diff value is nil instead of string")
	}

	switch diff.To.(type) {
	case string:
		cb(diff.To.(string))
	default:
		return fmt.Errorf("diff value type mismatch: string vs. %q", reflect.TypeOf(diff.To))
	}

	return nil
}

// PatchString applies the patch to a string value
func (diff *ValueDiff) PatchString(value *string) error {
	return diff.patchStringCB(func(s string) { *value = s })
}

// PatchUInt64Ref applies the patch to a *unit64 value
func (diff *ValueDiff) PatchUInt64Ref(value **uint64) error {
	if diff.Empty() {
		return nil
	}

	if diff.To == nil {
		return fmt.Errorf("diff value is nil instead of *uint64")
	}

	switch diff.To.(type) {
	case *uint64:
		*value = diff.To.(*uint64)
	default:
		return fmt.Errorf("diff value type mismatch: *uint64 vs. %q", reflect.TypeOf(diff.To))
	}

	return nil
}
