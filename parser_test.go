package xslx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	c := require.New(t)

	parser := &Parser{}
	err := parser.Parse("test-data/test.xlsx", func(sheet int, row [][]byte) {
		//fmt.Printf("%d %#v\n", sheet, row)
	})
	c.NoError(err)
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parser := &Parser{}
		err := parser.Parse("test-data/test.xlsx", func(sheet int, row [][]byte) {})
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
