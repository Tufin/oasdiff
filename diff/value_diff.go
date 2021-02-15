package diff

type ValueDiff struct {
	OldValue interface{} `json:"oldValue"`
	NewValue interface{} `json:"newValue"`
}

func getValueDiff(value1 interface{}, value2 interface{}) *ValueDiff {
	if value1 == value2 {
		return nil
	}

	return &ValueDiff{
		OldValue: value1,
		NewValue: value2,
	}
}

func getFloat64RefDiff(valueRef1 *float64, valueRef2 *float64) *ValueDiff {
	if valueRef1 == nil && valueRef2 == nil {
		return nil
	}

	if valueRef1 == nil && valueRef2 != nil {
		return &ValueDiff{
			OldValue: nil,
			NewValue: *valueRef2,
		}
	}

	if valueRef1 != nil && valueRef2 == nil {
		return &ValueDiff{
			OldValue: *valueRef1,
			NewValue: nil,
		}
	}

	return getValueDiff(*valueRef1, *valueRef2)
}

func getBoolRefDiff(valueRef1 *bool, valueRef2 *bool) *ValueDiff {
	if valueRef1 == nil && valueRef2 == nil {
		return nil
	}

	if valueRef1 == nil && valueRef2 != nil {
		return &ValueDiff{
			OldValue: nil,
			NewValue: *valueRef2,
		}
	}

	if valueRef1 != nil && valueRef2 == nil {
		return &ValueDiff{
			OldValue: *valueRef1,
			NewValue: nil,
		}
	}

	return getValueDiff(*valueRef1, *valueRef2)
}
