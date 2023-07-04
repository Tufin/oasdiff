package checker

type IBackwardCompatibilityErrors []IBackwardCompatibilityError

func (errs IBackwardCompatibilityErrors) HasLevelOrHigher(level Level) bool {
	for _, e := range errs {
		if e.GetLevel() >= level {
			return true
		}
	}
	return false
}

func (bcErrors IBackwardCompatibilityErrors) Len() int {
	return len(bcErrors)
}

func (bcErrors IBackwardCompatibilityErrors) Less(i, j int) bool {

	iv, jv := bcErrors[i], bcErrors[j]

	switch {
	case iv.GetLevel() != jv.GetLevel():
		return iv.GetLevel() > jv.GetLevel()
	case iv.GetPath() != jv.GetPath():
		return iv.GetPath() < jv.GetPath()
	case iv.GetOperation() != jv.GetOperation():
		return iv.GetOperation() < jv.GetOperation()
	case iv.GetId() != jv.GetId():
		return iv.GetId() < jv.GetId()
	case iv.GetText() != jv.GetText():
		return iv.GetText() < jv.GetText()
	default:
		return iv.GetComment() < jv.GetComment()
	}
}

func (bcErrors IBackwardCompatibilityErrors) Swap(i, j int) {
	bcErrors[i], bcErrors[j] = bcErrors[j], bcErrors[i]
}
