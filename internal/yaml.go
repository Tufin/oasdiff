package internal

import (
	"fmt"
	"io"
	"reflect"

	"gopkg.in/yaml.v3"
)

func printYAML(stdout io.Writer, output interface{}) error {
	if reflect.ValueOf(output).IsNil() {
		return nil
	}

	bytes, err := yaml.Marshal(output)
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "%s", bytes)
	return nil
}
