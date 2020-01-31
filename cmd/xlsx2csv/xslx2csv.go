package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/dcu/xlsx"
)

var exitFunc = os.Exit

func main() {
	run()
}

func run() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: xlsx2csv <file.xlsx>")
		exitFunc(0)
		return
	}

	csvWriter := csv.NewWriter(os.Stdout)
	defer csvWriter.Flush()

	parser := xlsx.NewParser()
	err := parser.Parse(os.Args[1], func(sheet int, row [][]byte) error {
		return csvWriter.Write(rowToStringArray(row))
	})

	if err != nil {
		fmt.Println("Error parsing file:", err.Error())
		exitFunc(1)
		return
	}
}

func rowToStringArray(row [][]byte) []string {
	result := []string{}
	for _, cell := range row {
		result = append(result, string(cell))
	}

	return result
}
