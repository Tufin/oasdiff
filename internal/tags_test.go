package internal_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/internal"
)

func Test_ChecksNoTags(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff checks -l ru"), io.Discard, io.Discard))
}

func Test_ChecksTagsDirection(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff checks -l ru --tags request"), io.Discard, io.Discard))
	require.Zero(t, internal.Run(cmdToArgs("oasdiff checks -l ru --tags response"), io.Discard, io.Discard))
}

func Test_ChecksTagsAction(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff checks -l ru --tags add"), io.Discard, io.Discard))
	require.Zero(t, internal.Run(cmdToArgs("oasdiff checks -l ru --tags remove"), io.Discard, io.Discard))
	require.Zero(t, internal.Run(cmdToArgs("oasdiff checks -l ru --tags change"), io.Discard, io.Discard))
	require.Zero(t, internal.Run(cmdToArgs("oasdiff checks -l ru --tags generalize"), io.Discard, io.Discard))
	require.Zero(t, internal.Run(cmdToArgs("oasdiff checks -l ru --tags specialize"), io.Discard, io.Discard))
	require.Zero(t, internal.Run(cmdToArgs("oasdiff checks -l ru --tags increase"), io.Discard, io.Discard))
	require.Zero(t, internal.Run(cmdToArgs("oasdiff checks -l ru --tags decrease"), io.Discard, io.Discard))
	require.Zero(t, internal.Run(cmdToArgs("oasdiff checks -l ru --tags set"), io.Discard, io.Discard))
}

func Test_ChecksTagsLocation(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff checks -l ru --tags body"), io.Discard, io.Discard))
	require.Zero(t, internal.Run(cmdToArgs("oasdiff checks -l ru --tags parameters"), io.Discard, io.Discard))
	require.Zero(t, internal.Run(cmdToArgs("oasdiff checks -l ru --tags properties"), io.Discard, io.Discard))
	require.Zero(t, internal.Run(cmdToArgs("oasdiff checks -l ru --tags headers"), io.Discard, io.Discard))
	require.Zero(t, internal.Run(cmdToArgs("oasdiff checks -l ru --tags security"), io.Discard, io.Discard))
	require.Zero(t, internal.Run(cmdToArgs("oasdiff checks -l ru --tags components"), io.Discard, io.Discard))
}
