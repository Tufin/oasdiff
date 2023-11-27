package formatters

import "github.com/tufin/oasdiff/checker"

type Change struct {
	Id          string        `json:"id,omitempty" yaml:"id,omitempty"`
	Text        string        `json:"text,omitempty" yaml:"text,omitempty"`
	Comment     string        `json:"comment,omitempty" yaml:"comment,omitempty"`
	Level       checker.Level `json:"level" yaml:"level"`
	Operation   string        `json:"operation,omitempty" yaml:"operation,omitempty"`
	OperationId string        `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Path        string        `json:"path,omitempty" yaml:"path,omitempty"`
	Source      string        `json:"source,omitempty" yaml:"source,omitempty"`
	IsBreaking  bool          `json:"-" yaml:"-"`
}

type Changes []Change

func NewChanges(originalChanges checker.Changes, l checker.Localizer) Changes {
	changes := make(Changes, len(originalChanges))
	for i, change := range originalChanges {
		changes[i] = Change{
			Id:          change.GetId(),
			Text:        change.GetText(l),
			Comment:     change.GetComment(),
			Level:       change.GetLevel(),
			Operation:   change.GetOperation(),
			OperationId: change.GetOperationId(),
			Path:        change.GetPath(),
			Source:      change.GetSource(),
		}
	}
	return changes
}
