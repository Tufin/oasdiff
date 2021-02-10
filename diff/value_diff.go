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
