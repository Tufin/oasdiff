package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

func printJSON(stdout io.Writer, output interface{}) error {
	if reflect.ValueOf(output).IsNil() {
		return nil
	}

	bytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "%s\n", bytes)
	return nil
}
