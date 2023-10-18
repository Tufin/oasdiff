package formatters

type Output string

const (
	OutputDiff      Output = "diff"
	OutputSummary   Output = "summary"
	OutputBreaking  Output = "breaking"
	OutputChangelog Output = "changelog"
	OutputChecks    Output = "checks"
)
