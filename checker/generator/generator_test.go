package generator_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker/generator"
)

func TestGenerator(t *testing.T) {
	out, err := os.Create("messages.yaml")
	require.NoError(t, err)
	defer out.Close()

	require.NoError(t, generator.Generate(generator.GetAll, out))
}

func TestTreeGeneratoFiler(t *testing.T) {
	file, err := os.Create("messages.yaml")
	require.NoError(t, err)
	defer file.Close()
	require.NoError(t, generator.Generate(generator.GetTree("tree.yaml"), file))
}

func TestTreeGenerator(t *testing.T) {
	var out bytes.Buffer
	require.NoError(t, generator.Generate(generator.GetTree("tree.yaml"), &out))
	count, err := lineCounter(&out)
	require.NoError(t, err)
	require.Equal(t, 260, count)
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
