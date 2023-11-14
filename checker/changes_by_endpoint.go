package checker

type Endpoint struct {
	Path      string
	Operation string
}

type ChangesByEndpoint map[Endpoint]*Changes

type GroupedChanges struct {
	APIChanges       ChangesByEndpoint
	ComponentChanges Changes
	SecurityChanges  Changes
}

func newGroupedChanges() GroupedChanges {
	return GroupedChanges{
		APIChanges: ChangesByEndpoint{},
	}
}

func groupChanges(changes Changes) GroupedChanges {

	result := newGroupedChanges()

	for _, change := range changes {
		switch change.(type) {
		case ApiChange:
			ep := Endpoint{Path: change.GetPath(), Operation: change.GetOperation()}
			if c, ok := result.APIChanges[ep]; ok {
				*c = append(*c, change)
			} else {
				result.APIChanges[ep] = &Changes{change}
			}
		}
	}
	return result
}
