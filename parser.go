package xlsx

import (
	"archive/zip"
	"io"
	"os"
)

// Callback is a function that gets the sheet number and the parsed row
type Callback func(sheet int, row [][]byte) error

// Parser is a parse in charge of handling the XLST file.
type Parser struct {
	sharedStrings [][]byte
}

// NewParser creates a new parser
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses the given file
func (p *Parser) Parse(filePath string, cb Callback) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	st, err := file.Stat()
	if err != nil {
		return err
	}

	return p.ParseReader(file, st.Size(), cb)
}

// ParseReader parses the data from the given reader
func (p *Parser) ParseReader(reader io.ReaderAt, size int64, cb Callback) error {
	zipReader, err := zip.NewReader(reader, size)
	if err != nil {
		return err
	}

	files := map[string]*zip.File{}
	for _, file := range zipReader.File {
		files[file.Name] = file
	}

	ssp := newSharedStringParser(p)
	err = ssp.loadSharedStrings(files["xl/sharedStrings.xml"])
	if err != nil {
		return err
	}

	sp := newSheetParser(1, p)
	return sp.loadSheet(files["xl/worksheets/sheet1.xml"], cb) // TODO: extend to support multi sheet docs
}
