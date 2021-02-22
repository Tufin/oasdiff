package diff

// ValueDiff describes the diff between two values
type ValueDiff struct {
	From interface{} `json:"from"`
	To   interface{} `json:"to"`
}

func getValueDiff(value1, value2 interface{}) *ValueDiff {
	if value1 == value2 {
		return nil
	}

	return &ValueDiff{
		From: value1,
		To:   value2,
	}
}

func getFloat64RefDiff(valueRef1, valueRef2 *float64) *ValueDiff {
	if valueRef1 == nil && valueRef2 == nil {
		return nil
	}

	if valueRef1 == nil && valueRef2 != nil {
		return &ValueDiff{
			From: nil,
			To:   *valueRef2,
		}
	}

	if valueRef1 != nil && valueRef2 == nil {
		return &ValueDiff{
			From: *valueRef1,
			To:   nil,
		}
	}

	return getValueDiff(*valueRef1, *valueRef2)
}

func getBoolRefDiff(valueRef1, valueRef2 *bool) *ValueDiff {
	if valueRef1 == nil && valueRef2 == nil {
		return nil
	}

	if valueRef1 == nil && valueRef2 != nil {
		return &ValueDiff{
			From: nil,
			To:   *valueRef2,
		}
	}

	if valueRef1 != nil && valueRef2 == nil {
		return &ValueDiff{
			From: *valueRef1,
			To:   nil,
		}
	}

	return getValueDiff(*valueRef1, *valueRef2)
}

func getStringRefDiff(valueRef1, valueRef2 *string) *ValueDiff {
	if valueRef1 == nil && valueRef2 == nil {
		return nil
	}

	if valueRef1 == nil && valueRef2 != nil {
		return &ValueDiff{
			From: nil,
			To:   *valueRef2,
		}
	}

	if valueRef1 != nil && valueRef2 == nil {
		return &ValueDiff{
			From: *valueRef1,
			To:   nil,
		}
	}

	return getValueDiff(*valueRef1, *valueRef2)
}
