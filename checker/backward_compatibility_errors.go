package checker

import "fmt"

type BackwardCompatibilityErrors []BackwardCompatibilityError

func (errs BackwardCompatibilityErrors) HasLevelOrHigher(level Level) bool {
	for _, e := range errs {
		if e.Level >= level {
			return true
		}
	}
	return false
}

func (bcErrors BackwardCompatibilityErrors) Len() int {
	return len(bcErrors)
}

func (bcErrors BackwardCompatibilityErrors) Less(i, j int) bool {
	iv, jv := bcErrors[i], bcErrors[j]

	switch {
	case iv.Level != jv.Level:
		return iv.Level > jv.Level
	case iv.Path != jv.Path:
		return iv.Path < jv.Path
	case iv.Operation != jv.Operation:
		return iv.Operation < jv.Operation
	case iv.Id != jv.Id:
		return iv.Id < jv.Id
	case iv.Text != jv.Text:
		return iv.Text < jv.Text
	default:
		return iv.Comment < jv.Comment
	}
}

func (bcErrors BackwardCompatibilityErrors) Swap(i, j int) {
	bcErrors[i], bcErrors[j] = bcErrors[j], bcErrors[i]
}

func (r *BackwardCompatibilityError) Error() string {
	var levelName string
	switch r.Level {
	case ERR:
		levelName = "error"
	case WARN:
		levelName = "warning"
	case INFO:
		levelName = "info"
	default:
		levelName = "issue"
	}
	return fmt.Sprintf("%s at %s, in API %s %s %s [%s]. %s", levelName, r.Source, r.Operation, r.Path, r.Text, r.Id, r.Comment)
}
