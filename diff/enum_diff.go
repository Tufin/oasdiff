package diff

func getEnumDiff(enum1 []interface{}, enum2 []interface{}) bool {

	if enum1 == nil && enum2 == nil {
		return false
	}

	if enum1 != nil && enum2 != nil {
		return !contained(enum1, enum2) || !contained(enum2, enum1)
	}

	return true
}

func contained(enum1 []interface{}, enum2 []interface{}) bool {
	for _, v1 := range enum1 {
		if !findValue(v1, enum2) {
			return false
		}
	}

	return true
}

func findValue(value interface{}, enum []interface{}) bool {
	for _, other := range enum {
		if value == other {
			return true
		}
	}
	return false
}
