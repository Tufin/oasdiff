package formatters

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

type NotImplementedFormatter struct{}

func (f NotImplementedFormatter) RenderDiff(*diff.Diff, RenderOpts) ([]byte, error) {
	return notImplemented()
}
func (f NotImplementedFormatter) RenderSummary(*diff.Diff, RenderOpts) ([]byte, error) {
	return notImplemented()
}
func (f NotImplementedFormatter) RenderBreakingChanges(checker.Changes, RenderOpts) ([]byte, error) {
	return notImplemented()
}
func (f NotImplementedFormatter) RenderChangelog(checker.Changes, RenderOpts) ([]byte, error) {
	return notImplemented()
}
func (f NotImplementedFormatter) RenderChecks([]Check, RenderOpts) ([]byte, error) {
	return notImplemented()
}
func (f NotImplementedFormatter) RenderFlatten(*openapi3.T, RenderOpts) ([]byte, error) {
	return notImplemented()
}

func notImplemented() ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}
