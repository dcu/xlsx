package xslx

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	numbersRx = regexp.MustCompile(`\d`)
)

func (p *Parser) loadSheet(f *zip.File, sheet int, cb func(sheet int, row []string)) error {
	reader, err := f.Open()
	if err != nil {
		return fmt.Errorf("opening shared strings file: %w", err)
	}
	defer func() {
		_ = reader.Close()
	}()

	decoder := xml.NewDecoder(reader)
	return p.loopRows(decoder, sheet, cb)
}

func (p *Parser) loopRows(decoder *xml.Decoder, sheet int, cb func(sheet int, row []string)) error {
	expectingString := false
	stringLocation := "inline"
	currentCell := ""

	var row []string
	totalColumns := 0

	count := 0
	for {
		// Read tokens from the XML document in a stream.
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "v" {
				expectingString = true
			} else if se.Name.Local == "c" {
				stringLocation = "inline"
				for _, attr := range se.Attr {
					if attr.Name.Local == "t" && attr.Value == "s" {
						stringLocation = "shared"
					} else if attr.Name.Local == "r" {
						currentCell = attr.Value
					}
				}
			} else if se.Name.Local == "row" {
				row = make([]string, totalColumns)
			} else if se.Name.Local == "dimension" {
				for _, attr := range se.Attr {
					if attr.Name.Local == "ref" {
						parts := strings.SplitN(attr.Value, ":", 2)
						totalColumns = columnToIndex(parts[1]) + 1
						break
					}
				}
			}
		case xml.CharData:
			if expectingString {
				//println(stringLocation, currentCell, string(se))
				if stringLocation == "shared" {
					pos, _ := strconv.Atoi(string(se))
					row[columnToIndex(currentCell)] = p.sharedStrings[pos]
				} else {
					row[columnToIndex(currentCell)] = string(se)
				}
			}
		case xml.EndElement:
			if se.Name.Local == "row" {
				cb(sheet, row)
				count++
			}
		}
	}

	return nil
}

func columnToIndex(columnName string) int {
	columnName = numbersRx.ReplaceAllString(columnName, "")
	sum := 0
	for i := 0; i < len(columnName); i++ {
		sum *= 26
		sum += int(columnName[i] - 'A' + 1)
	}

	return sum - 1
}
