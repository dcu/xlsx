package xlsx

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type sharedStringParser struct {
	*Parser
	expectingString bool
	currentString   []byte
}

func newSharedStringParser(p *Parser) *sharedStringParser {
	return &sharedStringParser{
		Parser: p,
	}
}

func (ssp *sharedStringParser) loadSharedStrings(f *zip.File) error {
	reader, err := f.Open()
	if err != nil {
		return fmt.Errorf("opening shared strings file: %w", err)
	}
	defer func() {
		_ = reader.Close()
	}()

	decoder := xml.NewDecoder(reader)

	return ssp.loopSharedStrings(decoder)
}

func (ssp *sharedStringParser) loopSharedStrings(decoder *xml.Decoder) error {
	for {
		// Read tokens from the XML document in a stream.
		t, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}

		if err := ssp.handleToken(t); err != nil {
			return err
		}
	}

	return nil
}

func (ssp *sharedStringParser) handleToken(token xml.Token) error {
	switch se := token.(type) {
	case xml.StartElement:
		switch se.Name.Local {
		case "t":
			ssp.expectingString = true
			ssp.currentString = []byte{}
		case "sst":
			for _, attr := range se.Attr {
				if attr.Name.Local == "uniqueCount" {
					size, _ := strconv.Atoi(attr.Value)
					ssp.sharedStrings = make([][]byte, 0, size)
					break
				}
			}
		}
	case xml.EndElement:
		if se.Name.Local == "t" {
			ssp.sharedStrings = append(ssp.sharedStrings, ssp.currentString)
		}
	case xml.CharData:
		if ssp.expectingString {
			buf := make([]byte, len(se))
			copy(buf, se)

			ssp.currentString = buf
			ssp.expectingString = false
		}
	}

	return nil
}
