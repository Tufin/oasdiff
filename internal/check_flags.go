package internal

type ChecksFlags struct {
	CommonFlags
}

func NewChecksFlags() *ChecksFlags {
	return &ChecksFlags{
		CommonFlags: NewCommonFlags(),
	}
}

func (flags *ChecksFlags) getSeverity() []string {
	return fixViperStringSlice(flags.getViper().GetStringSlice("severity"))
}

func (flags *ChecksFlags) getTags() []string {
	return fixViperStringSlice(flags.getViper().GetStringSlice("tags"))
}
