package internal

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/tufin/oasdiff/diff"
	"golang.org/x/exp/slices"
)

type StringSlice struct {
	value   *[]string
	changed bool
}

func newStringSliceValue(val []string, p *[]string) *StringSlice {
	ssv := new(StringSlice)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func validate(records []string) error {
	for _, record := range records {
		if !slices.Contains(diff.ExcludeDiffOptions, record) {
			return fmt.Errorf("invalid value: %s", record)
		}
	}
	return nil
}

func readAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)

	records, err := csvReader.Read()
	if err != nil {
		return []string{}, err
	}

	if err := validate(records); err != nil {
		return []string{}, err
	}

	return records, nil
}

func writeAsCSV(vals []string) (string, error) {
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	err := w.Write(vals)
	if err != nil {
		return "", err
	}
	w.Flush()
	return strings.TrimSuffix(b.String(), "\n"), nil
}

func (s *StringSlice) Set(val string) error {
	v, err := readAsCSV(val)
	if err != nil {
		return err
	}
	if !s.changed {
		*s.value = v
	} else {
		*s.value = append(*s.value, v...)
	}
	s.changed = true
	return nil
}

func (s *StringSlice) Type() string {
	return "stringSlice"
}

func (s *StringSlice) String() string {
	str, _ := writeAsCSV(*s.value)
	return "[" + str + "]"
}
