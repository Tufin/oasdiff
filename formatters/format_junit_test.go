package formatters_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/formatters"
)

var jUnitFormatter = formatters.JUnitFormatter{
	Localizer: MockLocalizer,
}

func TestJUnitLineLookup(t *testing.T) {
	f, err := formatters.Lookup(string(formatters.FormatJUnit), formatters.DefaultFormatterOpts())
	require.NoError(t, err)
	require.IsType(t, formatters.JUnitFormatter{}, f)
}

func TestJUnitFormatter_RenderBreakingChanges_OneFailure(t *testing.T) {
	testChanges := checker.Changes{
		checker.ComponentChange{
			Id:    "change_id",
			Level: checker.ERR,
		},
	}

	// check output
	output, err := jUnitFormatter.RenderBreakingChanges(testChanges, formatters.NewRenderOpts())
	assert.NoError(t, err)
	expectedOutput := `<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
  <testsuite package="com.oasdiff" time="0" tests="1" errors="0" failures="1" name="OASDiff">
    <testcase name="change_id" classname="OASDiff" time="0">
      <failure message="Breaking change detected">This is a breaking change.</failure>
    </testcase>
  </testsuite>
</testsuites>`
	assert.Equal(t, expectedOutput, string(output))
}

func TestJUnitFormatter_RenderBreakingChanges_Success(t *testing.T) {
	testChanges := checker.Changes{}

	// check output
	output, err := jUnitFormatter.RenderBreakingChanges(testChanges, formatters.NewRenderOpts())
	assert.NoError(t, err)
	expectedOutput := `<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
  <testsuite package="com.oasdiff" time="0" tests="0" errors="0" failures="0" name="OASDiff">
    <testcase name="no breaking changes detected" classname="OASDiff" time="0"></testcase>
  </testsuite>
</testsuites>`
	assert.Equal(t, expectedOutput, string(output))
}

func TestJUnitFormatter_NotImplemented(t *testing.T) {

	var err error
	_, err = jUnitFormatter.RenderDiff(nil, formatters.NewRenderOpts())
	assert.Error(t, err)

	_, err = jUnitFormatter.RenderSummary(nil, formatters.NewRenderOpts())
	assert.Error(t, err)

	_, err = jUnitFormatter.RenderChecks(nil, formatters.NewRenderOpts())
	assert.Error(t, err)

	_, err = jUnitFormatter.RenderFlatten(nil, formatters.NewRenderOpts())
	assert.Error(t, err)
}
