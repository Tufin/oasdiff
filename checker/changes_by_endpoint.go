package checker

type Endpoint struct {
	Path      string
	Operation string
}

type ChangesByEndpoint map[Endpoint]*Changes

func GroupChanges(changes Changes) ChangesByEndpoint {

	apiChanges := ChangesByEndpoint{}

	for _, change := range changes {
		switch change.(type) {
		case ApiChange:
			ep := Endpoint{Path: change.GetPath(), Operation: change.GetOperation()}
			if c, ok := apiChanges[ep]; ok {
				*c = append(*c, change)
			} else {
				apiChanges[ep] = &Changes{change}
			}
		}
	}

	return apiChanges
}
