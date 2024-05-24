package structs

import (
	"bufio"
	"errors"
	"net"
)

type Command struct {
	args []string
	conn net.Conn
}

func NewCommand(conn net.Conn) *Command {
	return &Command{
		args: []string{},
		conn: conn,
	}
}

// To read raw tcp connection and parse commands
type Parser struct {
	conn   net.Conn
	Reader *bufio.Reader
	line   []byte
	pos    int
}

// Creats a parse which reads from the given connection
func NewParser(conn net.Conn) *Parser {
	return &Parser{
		conn:   conn,
		Reader: bufio.NewReader(conn),
		line:   make([]byte, 0),
		pos:    0,
	}
}

func (p *Parser) current() byte {
	if p.atEnd() {
		return '\r'
	}
	return p.line[p.pos]
}

func (p *Parser) advance() {
	p.pos++
}

func (p *Parser) atEnd() bool {
	return p.pos >= len(p.line)
}

func (p *Parser) ReadLine() ([]byte, error) {
	line, err := p.Reader.ReadBytes('\r')
	if err != nil {
		return line, err
	}
	_, err = p.Reader.ReadByte()
	if err != nil {
		return nil, err
	}
	return line[:len(line)-1], nil
}

// To read string arguments from current line
func (p *Parser) readString() (s []byte, err error) {
	for p.current() != '"' && !p.atEnd() {
		curr := p.current()
		p.advance()
		next := p.current()
		// To handle escaped quote inside the string
		// eg. - "quoted \"text\" here"
		if curr == '\\' && next == '"' {
			s = append(s, '"')
			// Need to adance the pointer twice, for backslash and quote
			p.advance()
		} else {
			s = append(s, curr)
		}
	}

	if p.current() != '"' {
		return nil, errors.New("unbalanced quote in request")
	}
	// advancing from the closing quote
	p.advance()
	return
}
