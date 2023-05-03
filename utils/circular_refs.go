package utils

type VisitedRefs map[string]struct{}

func (v VisitedRefs) Add(refName string) {
	v[refName] = struct{}{}
}

func (v VisitedRefs) Remove(refName string) {
	delete(v, refName)
}

func (v VisitedRefs) IsVisited(refName string) bool {
	_, ok := v[refName]
	return ok
}
