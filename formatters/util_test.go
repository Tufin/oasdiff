package formatters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripANSIEscapeCodes(t *testing.T) {
	input := "deleted the \\u001b[1m'header'\\u001b[0m request parameter \\u001b[1m'network-policies'\\u001b[0m"
	output := StripANSIEscapeCodes([]byte(input))
	expectedOutput := "deleted the 'header' request parameter 'network-policies'"
	assert.Equal(t, expectedOutput, string(output))
}
