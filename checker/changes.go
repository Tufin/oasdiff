package checker

type Changes []Change

func (changes Changes) HasLevelOrHigher(level Level) bool {
	for _, change := range changes {
		if change.GetLevel() >= level {
			return true
		}
	}
	return false
}

func (changes Changes) GetLevelCount() map[Level]int {
	counts := map[Level]int{}
	for _, change := range changes {
		level := change.GetLevel()
		counts[level] = counts[level] + 1
	}
	return counts
}

func (changes Changes) Len() int {
	return len(changes)
}

func (changes Changes) Less(i, j int) bool {

	iv, jv := changes[i], changes[j]

	switch {
	case iv.GetLevel() != jv.GetLevel():
		return iv.GetLevel() > jv.GetLevel()
	case iv.GetPath() != jv.GetPath():
		return iv.GetPath() < jv.GetPath()
	case iv.GetOperation() != jv.GetOperation():
		return iv.GetOperation() < jv.GetOperation()
	case iv.GetId() != jv.GetId():
		return iv.GetId() < jv.GetId()
	case len(iv.GetArgs()) != len(jv.GetArgs()):
		return len(iv.GetArgs()) < len(jv.GetArgs())
	default:
		for i, arg := range iv.GetArgs() {
			ia := interfaceToString(arg)
			ja := interfaceToString(jv.GetArgs()[i])
			if ia != ja {
				return ia < ja
			}
		}
	}
	return true
}

func (changes Changes) Swap(i, j int) {
	changes[i], changes[j] = changes[j], changes[i]
}
