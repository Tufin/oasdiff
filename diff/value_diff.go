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

func (diff *ValueDiff) minBreakingFloat64(direction direction) bool {
	if diff.Empty() {
		return false
	}

	from, _ := diff.From.(float64)
	to, _ := diff.To.(float64)

	return minBreaking(direction, diff.From == nil, diff.To == nil, from, to)
}

type numeric interface {
	uint64 |
		float64
}

func lessThan[N numeric](aNil, bNil bool, a, b N) bool {
	if aNil {
		return true
	}
	if bNil {
		return false
	}
	return b < a
}

func maxBreaking[N numeric](direction direction, fromNil, toNil bool, from, to N) bool {
	switch direction {
	case directionRequest:
		return lessThan(fromNil, toNil, from, to)
	case directionResponse:
		return lessThan(toNil, fromNil, to, from)
	}

	return false
}

func minBreaking[N numeric](direction direction, fromNil, toNil bool, from, to N) bool {
	switch direction {
	case directionRequest:
		return lessThan(toNil, fromNil, to, from)
	case directionResponse:
		return lessThan(fromNil, toNil, from, to)
	}

	return false
}

func (diff *ValueDiff) minBreakingUInt64(direction direction) bool {
	if diff.Empty() {
		return false
	}

	from, _ := diff.From.(uint64)
	to, _ := diff.To.(uint64)

	return minBreaking(direction, diff.From == nil, diff.To == nil, from, to)
}

func (diff *ValueDiff) maxBreakingFloat64(direction direction) bool {
	if diff.Empty() {
		return false
	}

	from, _ := diff.From.(float64)
	to, _ := diff.To.(float64)

	return maxBreaking(direction, diff.From == nil, diff.To == nil, from, to)
}

func (diff *ValueDiff) maxBreakingUInt64(direction direction) bool {
	if diff.Empty() {
		return false
	}

	from, _ := diff.From.(uint64)
	to, _ := diff.To.(uint64)

	return maxBreaking(direction, diff.From == nil, diff.To == nil, from, to)
}

func getValueWithDefault(value interface{}, defaultValue interface{}) interface{} {

	if value == nil {
		return defaultValue
	}
	return value
}

func getValueDiff(value1, value2 interface{}) *ValueDiff {

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

func getValueDiffConditional(exclude bool, value1, value2 interface{}) *ValueDiff {
	if exclude {
		return nil
	}

	return getValueDiff(value1, value2)
}

func getFloat64RefDiff(valueRef1, valueRef2 *float64) *ValueDiff {
	return getValueDiff(derefFloat64(valueRef1), derefFloat64(valueRef2))
}

func getBoolRefDiff(valueRef1, valueRef2 *bool) *ValueDiff {
	return getValueDiff(derefBool(valueRef1), derefBool(valueRef2))
}

func getStringRefDiffConditional(exclude bool, valueRef1, valueRef2 *string) *ValueDiff {
	return getValueDiffConditional(exclude, derefString(valueRef1), derefString(valueRef2))
}

func getUInt64RefDiff(valueRef1, valueRef2 *uint64) *ValueDiff {
	return getValueDiff(derefUInt64(valueRef1), derefUInt64(valueRef2))
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
