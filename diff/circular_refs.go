package diff

type visitedRefs map[string]struct{}

func (v visitedRefs) add(refName string) {
	v[refName] = struct{}{}
}

func (v visitedRefs) remove(refName string) {
	delete(v, refName)
}

func (v visitedRefs) isVisited(refName string) bool {
	_, ok := v[refName]
	return ok
}
