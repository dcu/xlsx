package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_run(t *testing.T) {
	os.Args = []string{"xlsx2csv", "../../test-data/test.xlsx"}
	run()
}

func Test_runError(t *testing.T) {
	c := require.New(t)

	exitFunc = func(code int) {
		c.Equal(0, code)
	}

	os.Args = []string{}
	run()
}

func Test_runNoFile(t *testing.T) {
	c := require.New(t)

	os.Args = []string{"xlsx2csv", "../../test-data/test.bar"}
	exitFunc = func(code int) {
		c.Equal(1, code)
	}

	run()
}

func Test_rowToStringArray(t *testing.T) {
	c := require.New(t)

	c.Equal([]string{"a", "b", "c", "d"}, rowToStringArray([][]byte{
		[]byte("a"),
		[]byte("b"),
		[]byte("c"),
		[]byte("d"),
	}))
}
