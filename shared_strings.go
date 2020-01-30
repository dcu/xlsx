package xslx

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
)

func (p *Parser) loadSharedStrings(f *zip.File) error {
	reader, err := f.Open()
	if err != nil {
		return fmt.Errorf("opening shared strings file: %w", err)
	}
	defer func() {
		_ = reader.Close()
	}()

	decoder := xml.NewDecoder(reader)

	return p.loopSharedStrings(decoder)
}

func (p *Parser) loopSharedStrings(decoder *xml.Decoder) error {
	expectingString := false
	currentString := ""

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
			if se.Name.Local == "sst" {
				for _, attr := range se.Attr {
					if attr.Name.Local == "uniqueCount" {
						size, _ := strconv.Atoi(attr.Value)
						p.sharedStrings = make([]string, 0, size)
					}
				}

			} else if se.Name.Local == "t" {
				expectingString = true
				currentString = ""
			}
		case xml.EndElement:
			if se.Name.Local == "t" {
				p.sharedStrings = append(p.sharedStrings, currentString)
			}
		case xml.CharData:
			if expectingString {
				currentString = string(se)
				expectingString = false
			}
		}
	}

	return nil
}
