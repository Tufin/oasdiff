package diff

import (
	"fmt"
	"reflect"
)

// ValueDiff describes the changes between a pair of values
type ValueDiff struct {
	From interface{} `json:"from" yaml:"from"`
	To   interface{} `json:"to" yaml:"to"`
}

// Empty indicates whether a change was found in this element
func (diff *ValueDiff) Empty() bool {
	return diff == nil
}

// CompareWithDefault checks if value was changed from a specific value to another specific value
// For example: was the value changed from 'true' to 'false'?
// If the original value or the new value are not defined, the comparison uses the default value
func (diff *ValueDiff) CompareWithDefault(from, to, defaultValue interface{}) bool {
	if diff.Empty() {
		return false
	}

	return getValueWithDefault(diff.From, defaultValue) == from &&
		getValueWithDefault(diff.To, defaultValue) == to
}

func (diff *ValueDiff) MinBreaking() bool {
	if diff.Empty() {
		return false
	}

	if diff.From == nil {
		return true
	}

	if diff.To == nil {
		return false
	}

	if diff.From.(uint64) < diff.To.(uint64) {
		return true
	}

	return false
}

func (diff *ValueDiff) MaxBreaking() bool {
	if diff.Empty() {
		return false
	}

	if diff.From == nil {
		return true
	}

	if diff.To == nil {
		return false
	}

	if diff.To.(uint64) < diff.From.(uint64) {
		return true
	}

	return false
}

func getValueWithDefault(value interface{}, defaultValue interface{}) interface{} {

	if value == nil {
		return defaultValue
	}
	return value
}

func getValueDiff(config *Config, value1, value2 interface{}) *ValueDiff {

	diff := getValueDiffInternal(value1, value2)

	if diff.Empty() {
		return nil
	}

	return diff
}

func getValueDiffInternal(value1, value2 interface{}) *ValueDiff {
	if reflect.DeepEqual(value1, value2) {
		return nil
	}

	return &ValueDiff{
		From: value1,
		To:   value2,
	}
}

func getValueDiffConditional(config *Config, exclude bool, value1, value2 interface{}) *ValueDiff {
	if exclude {
		return nil
	}

	return getValueDiff(config, value1, value2)
}

func getFloat64RefDiff(config *Config, valueRef1, valueRef2 *float64) *ValueDiff {
	return getValueDiff(config, derefFloat64(valueRef1), derefFloat64(valueRef2))
}

func getBoolRefDiff(config *Config, valueRef1, valueRef2 *bool) *ValueDiff {
	return getValueDiff(config, derefBool(valueRef1), derefBool(valueRef2))
}

func getStringRefDiffConditional(config *Config, exclude bool, valueRef1, valueRef2 *string) *ValueDiff {
	return getValueDiffConditional(config, exclude, derefString(valueRef1), derefString(valueRef2))
}

func getUInt64RefDiff(config *Config, valueRef1, valueRef2 *uint64) *ValueDiff {
	return getValueDiff(config, derefUInt64(valueRef1), derefUInt64(valueRef2))
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

func derefUInt64(ref *uint64) interface{} {
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

// patchString applies the patch to a string value
func (diff *ValueDiff) patchString(value *string) error {
	return diff.patchStringCB(func(s string) { *value = s })
}

// patchUInt64Ref applies the patch to a *unit64 value
func (diff *ValueDiff) patchUInt64Ref(value **uint64) error {
	if diff.Empty() {
		return nil
	}

	if diff.To == nil {
		*value = nil
		return nil
	}

	switch t := diff.To.(type) {
	case uint64:
		*value = &t
	default:
		return fmt.Errorf("diff value type mismatch: uint64 vs. %q", reflect.TypeOf(diff.To))
	}

	return nil
}
