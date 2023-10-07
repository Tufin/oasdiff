package load

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
