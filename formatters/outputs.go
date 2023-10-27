package formatters

type Output int

const (
	OutputDiff Output = iota
	OutputSummary
	OutputBreaking
	OutputChangelog
	OutputChecks
	OutputFlatten
)
