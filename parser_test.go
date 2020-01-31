package xlsx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	c := require.New(t)

	rows := [][][]byte{}

	parser := NewParser()
	err := parser.Parse("test-data/test.xlsx", func(sheet int, row [][]byte) error {
		rows = append(rows, row)
		return nil
	})

	c.Equal(13, len(parser.sharedStrings))
	c.Equal(23, len(rows))

	c.Equal("title1", string(rows[1][0]))
	c.Equal("4", string(rows[3][3]))
	c.NoError(err)
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parser := NewParser()
		err := parser.Parse("test-data/test.xlsx", func(sheet int, row [][]byte) error { return nil })
		if err != nil {
			b.Fatal(err)
		}

	}
}

func Benchmark_columnToIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		index := columnToIndex("AAAAAAAA12345")
		if index != 8353082582 {
			b.Fatal("wrong index", index)
		}
	}
}
