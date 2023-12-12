package load

import (
	"fmt"
	"net/url"
)

type SourceType int

const (
	SourceTypeStdin SourceType = iota
	SourceTypeURL
	SourceTypeFile
)

type Source struct {
	Path string
	Uri  *url.URL
	Type SourceType
}

func NewSource(path string) *Source {
	if path == "-" {
		return &Source{
			Path: "stdin",
			Type: SourceTypeStdin,
		}
	}

	if uri, err := getURL(path); err == nil {
		return &Source{
			Path: path,
			Type: SourceTypeURL,
			Uri:  uri,
		}
	}

	return &Source{
		Path: path,
		Type: SourceTypeFile,
	}
}

func (source *Source) String() string {
	return source.Path
}

func (source *Source) Out() string {
	if source.IsStdin() {
		return source.Path
	}
	return fmt.Sprintf("%q", source.Path)
}

func (source *Source) IsStdin() bool {
	return source.Type == SourceTypeStdin
}

func (source *Source) IsFile() bool {
	return source.Type == SourceTypeFile
}
