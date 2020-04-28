package bencode

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"errors"
)

// ErrSyntax describes a bencode syntax error.
type ErrSyntax struct {
	msg string
	pos int64
}

func (e *ErrSyntax) Error() string { return fmt.Sprintf("bencode: %d: %s", e.pos, e.msg) }

var errValueEnd = errors.New("value end")

// Parser parses
type Parser struct {
	reader *bufio.Reader
	offset int64
	// tree   *Value
}

// NewParser returns a new parser
func NewParser(r io.Reader) *Parser {
	br := bufio.NewReader(r)
	return &Parser{reader: br}
}

// Parse parses.
func (p *Parser) Parse() (v Value, err error) {
	v, err = p.parseValue()
	return
}

func (p *Parser) parseValue() (Value, error) {
	bs, err := p.reader.Peek(1)
	if err != nil {
		return nil, err
	}
	b := bs[0]

	switch b {
	case 'i':
		return p.parseInt()
	case 'l':
		return p.parseList()
	case 'd':
		return p.parseDict()
	case 'e':
		return p.parseEndOfValue()
	default:
		if b >= '0' && b <= '9' {
			return p.parseString()
		}
	}
	return nil, &ErrSyntax{msg: "unexpected token"}
}

func (p *Parser) parseInt() (Int, error) {
	if err := p.skipDelimeter(); err != nil {
		return Int(0), err
	}

	// read until delimeter 'e'
	s, err := p.reader.ReadString('e')
	p.offset += int64(len(s))
	if err != nil {
		return Int(0), &ErrSyntax{pos: p.offset, msg: "cannot find the end delimeter of the integer"}
	}
	s = s[:len(s)-1] // trim delimeters 'e'

	// parse as integer
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return Int(0), &ErrSyntax{pos: p.offset, msg: fmt.Sprintf("cannot parse '%v' as integer", s)}
	}

	return Int(i), nil
}

func (p *Parser) parseString() (String, error) {
	// parse string length
	var length int64
	s, err := p.reader.ReadString(':')
	if err != nil {
		return String(""), &ErrSyntax{pos: p.offset, msg: "cannot find string length delimeter"}
	}
	s = s[:len(s)-1] // trim delimeter ':'
	length, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		return String(""), &ErrSyntax{pos: p.offset, msg: fmt.Sprintf("cannot parse string length '%v' as integer", s)}
	}
	p.offset += int64(len(s) + 1)

	// parse string value
	bs := make([]byte, length)
	n, err := io.ReadFull(p.reader, bs)
	if err != nil {
		return String(""), &ErrSyntax{pos: p.offset, msg: "string length is wrong"}
	}
	p.offset += int64(n)

	return String(bs), nil
}

func (p *Parser) parseList() (List, error) {
	list := List{}

	if err := p.skipDelimeter(); err != nil {
		return list, err
	}

ParseValuesLoop:
	for {
		item, err := p.parseValue()
		if err != nil {
			if err == errValueEnd {
				break ParseValuesLoop
			}
			return list, err
		}
		list = append(list, item)
	}

	return list, nil
}

func (p *Parser) parseDict() (*Dict, error) {
	dict := NewDict()

	if err := p.skipDelimeter(); err != nil {
		return dict, err
	}

ParseItemsLoop:
	for {
		// parse item key
		v, err := p.parseValue()
		if err != nil {
			if err == errValueEnd {
				break ParseItemsLoop
			}
			return dict, err
		}
		key, ok := v.(String)
		if !ok {
			return dict, &ErrSyntax{pos: p.offset, msg: "dict key is not a string"}
		}

		// parse item value
		value, err := p.parseValue()
		if err != nil {
			return dict, err
		}

		dict.Set(key, value)
	}

	return dict, nil
}

func (p *Parser) parseEndOfValue() (Value, error) {
	if b, err := p.reader.ReadByte(); err != nil {
		return nil, &ErrSyntax{pos: p.offset, msg: "unexpected end of data"}
	} else if b != 'e' {
		return nil, &ErrSyntax{pos: p.offset, msg: "unexpected token"}
	}
	p.offset++
	return nil, errValueEnd
}

func (p *Parser) skipDelimeter() error {
	b, err := p.reader.ReadByte()
	if err != nil {
		return &ErrSyntax{pos: p.offset, msg: "unexpected end of data"}
	}
	if b != 'i' && b != 'l' && b != 'd' {
		return &ErrSyntax{pos: p.offset, msg: "unexpected token"}
	}
	p.offset++
	if b == 'e' {
		return errValueEnd
	}
	return nil
}
