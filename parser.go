package xslx

import (
	"archive/zip"
	"io"
	"os"
)

// Parser is a parse in charge of handling the XLST file.
type Parser struct {
	sharedStrings []string
}

// Parse parses the given file
func (p *Parser) Parse(filePath string, cb func(sheet int, row []string)) error {
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
func (p *Parser) ParseReader(reader io.ReaderAt, size int64, cb func(sheet int, row []string)) error {
	zipReader, err := zip.NewReader(reader, size)
	if err != nil {
		return err
	}

	files := map[string]*zip.File{}
	for _, file := range zipReader.File {
		files[file.Name] = file
	}

	err = p.loadSharedStrings(files["xl/sharedStrings.xml"])
	if err != nil {
		return err
	}

	return p.loadSheet(files["xl/worksheets/sheet1.xml"], 1, cb) // TODO: extend to multi sheet
}
