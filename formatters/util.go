package formatters

import (
	"regexp"
)

var stripAnsiEscapesRegex = regexp.MustCompile(`\\u001b\[[0-9]*m|\\e\[[0-9]*m`)

// StripANSIEscapeCodes removes ANSI escape codes
// TODO: remove this function once there is a better way to prevent ANSI escape codes in messages for json/yaml output
func StripANSIEscapeCodes(data []byte) []byte {
	return stripAnsiEscapesRegex.ReplaceAll(data, nil)
}

// StripANSIEscapeCodesStr removes ANSI escape codes
// TODO: remove this function once there is a better way to prevent ANSI escape codes in messages for json/yaml output
func StripANSIEscapeCodesStr(data string) string {
	return stripAnsiEscapesRegex.ReplaceAllString(data, "")
}
