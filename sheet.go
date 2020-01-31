package xlsx

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type sheetParser struct {
	parser          *Parser
	sheet           int
	expectingString bool
	stringLocation  string
	currentCell     string
	totalColumns    int
	currentRow      [][]byte
}

func newSheetParser(sheet int, parser *Parser) *sheetParser {
	return &sheetParser{
		parser: parser,
		sheet:  sheet,
	}
}

func (sp *sheetParser) loadSheet(f *zip.File, cb func(sheet int, row [][]byte)) error {
	reader, err := f.Open()
	if err != nil {
		return fmt.Errorf("opening shared strings file: %w", err)
	}
	defer func() {
		_ = reader.Close()
	}()

	decoder := xml.NewDecoder(reader)
	return sp.loopRows(decoder, cb)
}

func (sp *sheetParser) loopRows(decoder *xml.Decoder, cb func(sheet int, row [][]byte)) error {
	for {
		// Read tokens from the XML document in a stream.
		t, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}

		if err := sp.handleToken(t, cb); err != nil {
			return err
		}
	}

	return nil
}

func (sp *sheetParser) handleToken(t xml.Token, cb func(sheet int, row [][]byte)) error {
	// Inspect the type of the token just read.
	switch se := t.(type) {
	case xml.StartElement:
		sp.handleStartElement(&se)
	case xml.CharData:
		sp.handleCharData(se)
	case xml.EndElement:
		sp.handleEndElement(&se, cb)
	}

	return nil
}

func (sp *sheetParser) handleEndElement(se *xml.EndElement, cb func(sheet int, row [][]byte)) {
	if se.Name.Local == "row" {
		cb(sp.sheet, sp.currentRow)
	}
}

func (sp *sheetParser) handleCharData(se xml.CharData) {
	if !sp.expectingString {
		return
	}

	if sp.stringLocation == "shared" {
		pos, _ := strconv.Atoi(string(se))
		sp.currentRow[columnToIndex(sp.currentCell)] = []byte(sp.parser.sharedStrings[pos])
	} else {
		buf := make([]byte, len(se))
		copy(buf, se)

		sp.currentRow[columnToIndex(sp.currentCell)] = buf
	}
}

func (sp *sheetParser) handleStartElement(se *xml.StartElement) {
	switch se.Name.Local {
	case "v":
		sp.expectingString = true
	case "c":
		sp.stringLocation = "inline"
		sp.parseCellAttributes(se)
	case "row":
		sp.currentRow = make([][]byte, sp.totalColumns)
	case "dimension":
		for _, attr := range se.Attr {
			if attr.Name.Local == "ref" {
				parts := strings.SplitN(attr.Value, ":", 2)
				sp.totalColumns = columnToIndex(parts[1]) + 1
				break
			}
		}
	}
}

func (sp *sheetParser) parseCellAttributes(cell *xml.StartElement) {
	for _, attr := range cell.Attr {
		if attr.Name.Local == "t" && attr.Value == "s" {
			sp.stringLocation = "shared"
		} else if attr.Name.Local == "r" {
			sp.currentCell = attr.Value
		}
	}
}

func columnToIndex(columnName string) int {
	sum := 0
	for i := 0; i < len(columnName); i++ {
		ch := columnName[i]
		if ch < 'A' || ch > 'Z' {
			break
		}

		sum *= 26
		sum += int(columnName[i] - 'A' + 1)
	}

	return sum - 1
}
