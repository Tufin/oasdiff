package internal_test

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/internal"
)

type ViperMock struct {
	viper.Viper
	ReadInConfigMock   func() error
	BindPFlagMock      func(key string, flag *pflag.Flag) error
	UnmarshalExactMock func(rawVal any, opts ...viper.DecoderConfigOption) error
}

func NewViperMock() *ViperMock {
	result := ViperMock{
		Viper: *viper.GetViper(),
	}
	result.ReadInConfigMock = result.Viper.ReadInConfig
	result.BindPFlagMock = result.Viper.BindPFlag
	result.UnmarshalExactMock = result.Viper.UnmarshalExact
	return &result
}

func (v *ViperMock) ReadInConfig() error {
	return v.ReadInConfigMock()
}

func (v *ViperMock) BindPFlag(key string, flag *pflag.Flag) error {
	return v.BindPFlagMock(key, flag)
}

func (v *ViperMock) UnmarshalExact(rawVal any, opts ...viper.DecoderConfigOption) error {
	return v.UnmarshalExactMock(rawVal, opts...)
}

func TestViper_ReadInConfigErr(t *testing.T) {
	v := NewViperMock()
	v.ReadInConfigMock = func() error { return errors.New("error") }

	cmd := cobra.Command{}
	require.EqualError(t, internal.RunViper(&cmd, v), "failed to load config file: read error: error \n")
}

func TestViper_BindPFlagErr(t *testing.T) {
	v := NewViperMock()
	v.BindPFlagMock = func(key string, flag *pflag.Flag) error {
		return errors.New("error")
	}

	cmd := cobra.Command{}
	cmd.PersistentFlags().BoolP("composed", "c", false, "work in 'composed' mode, compare paths in all specs matching base and revision globs")

	require.EqualError(t, internal.RunViper(&cmd, v), "failed to load config file: error binding flag \"composed\" to viper: error")
}

func TestViper_UnmarshalExactErr(t *testing.T) {
	v := NewViperMock()
	v.UnmarshalExactMock = func(rawVal any, opts ...viper.DecoderConfigOption) error {
		return errors.New("error")
	}

	cmd := cobra.Command{}

	require.EqualError(t, internal.RunViper(&cmd, v), "failed to load config file: validation error: error \n")
}

func TestViper_InvalidLang(t *testing.T) {
	v := NewViperMock()

	cmd := cobra.Command{}

	v.UnmarshalExactMock = func(rawVal any, opts ...viper.DecoderConfigOption) error {
		rawVal.(*internal.Config).Lang = "invalid"
		return nil
	}

	require.EqualError(t, internal.RunViper(&cmd, v), "failed to load config file: invalid lang \"invalid\", allowed values: en, ru")
}

func TestViper_InvalidColor(t *testing.T) {
	v := NewViperMock()

	cmd := cobra.Command{}

	v.UnmarshalExactMock = func(rawVal any, opts ...viper.DecoderConfigOption) error {
		rawVal.(*internal.Config).Color = "invalid"
		return nil
	}

	require.EqualError(t, internal.RunViper(&cmd, v), "failed to load config file: invalid color \"invalid\", allowed values: auto, always, never")
}

func TestViper_InvalidFormat(t *testing.T) {
	v := NewViperMock()

	cmd := cobra.Command{}

	v.UnmarshalExactMock = func(rawVal any, opts ...viper.DecoderConfigOption) error {
		rawVal.(*internal.Config).Format = "invalid"
		return nil
	}

	require.EqualError(t, internal.RunViper(&cmd, v), "failed to load config file: invalid format \"invalid\", allowed values: yaml, json, text, markup, singleline, html, githubactions, junit, sarif")
}

func TestViper_InvalidFailOn(t *testing.T) {
	v := NewViperMock()

	cmd := cobra.Command{}

	v.UnmarshalExactMock = func(rawVal any, opts ...viper.DecoderConfigOption) error {
		rawVal.(*internal.Config).FailOn = "invalid"
		return nil
	}

	require.EqualError(t, internal.RunViper(&cmd, v), "failed to load config file: invalid fail-on \"invalid\", allowed values: ERR, WARN, INFO")
}

func TestViper_InvalidLevel(t *testing.T) {
	v := NewViperMock()

	cmd := cobra.Command{}

	v.UnmarshalExactMock = func(rawVal any, opts ...viper.DecoderConfigOption) error {
		rawVal.(*internal.Config).Level = "invalid"
		return nil
	}

	require.EqualError(t, internal.RunViper(&cmd, v), "failed to load config file: invalid level \"invalid\", allowed values: ERR, WARN, INFO")
}

func TestViper_InvalidExcludeElements(t *testing.T) {
	v := NewViperMock()

	cmd := cobra.Command{}

	v.UnmarshalExactMock = func(rawVal any, opts ...viper.DecoderConfigOption) error {
		rawVal.(*internal.Config).ExcludeElements = []string{"invalid"}
		return nil
	}

	require.EqualError(t, internal.RunViper(&cmd, v), "failed to load config file: invalid exclude-elements \"invalid\", allowed values: examples, description, endpoints, title, summary, extensions")
}

func TestViper_InvalidSeverity(t *testing.T) {
	v := NewViperMock()

	cmd := cobra.Command{}

	v.UnmarshalExactMock = func(rawVal any, opts ...viper.DecoderConfigOption) error {
		rawVal.(*internal.Config).Severity = []string{"invalid"}
		return nil
	}

	require.EqualError(t, internal.RunViper(&cmd, v), "failed to load config file: invalid severity \"invalid\", allowed values: info, warn, error")
}

func TestViper_InvalidTags(t *testing.T) {
	v := NewViperMock()

	cmd := cobra.Command{}

	v.UnmarshalExactMock = func(rawVal any, opts ...viper.DecoderConfigOption) error {
		rawVal.(*internal.Config).Tags = []string{"invalid"}
		return nil
	}

	require.EqualError(t, internal.RunViper(&cmd, v), "failed to load config file: invalid tags \"invalid\", allowed values: request, response, add, remove, change, generalize, specialize, increase, decrease, set, body, parameters, properties, headers, security, components")
}

func TestViper_ValidTags(t *testing.T) {
	v := NewViperMock()

	cmd := cobra.Command{}

	v.UnmarshalExactMock = func(rawVal any, opts ...viper.DecoderConfigOption) error {
		rawVal.(*internal.Config).Tags = []string{"request"}
		return nil
	}

	require.Nil(t, internal.RunViper(&cmd, v))
}
