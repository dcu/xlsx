# xslx

[![Build Status](https://travis-ci.com/dcu/xlsx.svg?branch=master)](https://travis-ci.org/dcu/xlsx)

```go
    import "github.com/dcu/xlsx"
```

## Usage

#### type Parser

```go
type Parser struct {
}
```

Parser is a parse in charge of handling the XLST file.

#### func  NewParser

```go
func NewParser() *Parser
```
NewParser creates a new parser

#### func (*Parser) Parse

```go
func (p *Parser) Parse(filePath string, cb func(sheet int, row [][]byte)) error
```
Parse parses the given file

#### func (*Parser) ParseReader

```go
func (p *Parser) ParseReader(reader io.ReaderAt, size int64, cb func(sheet int, row [][]byte)) error
```
ParseReader parses the data from the given reader
