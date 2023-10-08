package load

import "fmt"

type Source struct {
	Path  string
	Stdin bool
}

func GetSource(path string) Source {
	stdin := path == "-"
	if stdin {
		path = "stdin"
	}

	return Source{
		Path:  path,
		Stdin: stdin,
	}
}

func (source Source) Out() string {
	if source.Stdin {
		return source.Path
	}
	return fmt.Sprintf("%q", source.Path)
}
