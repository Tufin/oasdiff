package diff

import "reflect"

// ValueDiff describes the diff between two values
type ValueDiff struct {
	From interface{} `json:"from"`
	To   interface{} `json:"to"`
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
