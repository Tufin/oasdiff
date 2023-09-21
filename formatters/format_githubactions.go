package formatters

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

var githubActionsSeverity = map[checker.Level]string{
	checker.ERR:  "error",
	checker.WARN: "warning",
	checker.INFO: "notice",
}

type GitHubActionsFormatter struct {
	Localizer checker.Localizer
}

func (f GitHubActionsFormatter) RenderDiff(diff *diff.Diff, opts RenderOpts) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (f GitHubActionsFormatter) RenderSummary(diff *diff.Diff, opts RenderOpts) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (f GitHubActionsFormatter) RenderBreakingChanges(changes checker.Changes, opts RenderOpts) ([]byte, error) {
	var buf bytes.Buffer

	for _, change := range changes {
		// source file, line and column are optional
		var params = []string{
			"title=" + change.GetId(),
		}
		if change.GetSourceFile() != "" {
			params = append(params, "file="+change.GetSourceFile())
		}
		if change.GetSourceColumn() != 0 {
			params = append(params, "col="+strconv.Itoa(change.GetSourceColumn()+1))
		}
		if change.GetSourceColumnEnd() != 0 {
			params = append(params, "endColumn="+strconv.Itoa(change.GetSourceColumnEnd()+1))
		}
		if change.GetSourceLine() != 0 {
			params = append(params, "line="+strconv.Itoa(change.GetSourceLine()+1))
		}
		if change.GetSourceLineEnd() != 0 {
			params = append(params, "endLine="+strconv.Itoa(change.GetSourceLineEnd()+1))
		}

		// all annotated messages must be one-line, due to GitHub Actions limitations
		message := StripANSIEscapeCodes([]byte(strings.ReplaceAll(change.GetText(), "\n", "%0A")))

		buf.WriteString(fmt.Sprintf("::%s %s::%s\n", githubActionsSeverity[change.GetLevel()], strings.Join(params, ","), message))
	}

	return buf.Bytes(), nil
}

func (f GitHubActionsFormatter) RenderChangelog(changes checker.Changes, opts RenderOpts) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (f GitHubActionsFormatter) RenderChecks(rules []checker.BackwardCompatibilityRule, opts RenderOpts) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (f GitHubActionsFormatter) SupportedOutputs() []string {
	return []string{"breaking-changes"}
}
