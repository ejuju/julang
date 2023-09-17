package julang

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

type Position struct{ Line, Column int }

func (p Position) String() string { return strconv.Itoa(p.Line) + ":" + strconv.Itoa(p.Column) }

type TokenType uint8

const (
	_ TokenType = iota
	TokenTypeWord
	TokenTypeLiteralText
	TokenTypeQuotation
)

var TokenLabels = [...]string{
	TokenTypeWord:        "word",
	TokenTypeLiteralText: "text",
	TokenTypeQuotation:   "quotation",
}

func (tt TokenType) String() string { return TokenLabels[tt] }

type Token struct {
	Position Position
	Type     TokenType
	Value    string
}

func (t Token) String() string { return fmt.Sprintf("%s (%s) %q", t.Type, t.Position, t.Value) }

type SyntaxError struct {
	Position Position
	Message  string
}

func (err SyntaxError) Error() string { return fmt.Sprintf("%s (%s)", err.Message, err.Position) }

type TokenStream struct {
	r io.Reader
	p Position
}

func NewTokenStream(r io.Reader) *TokenStream { return &TokenStream{r: r, p: Position{1, 1}} }

func (ts *TokenStream) read() (byte, error) {
	buf := make([]byte, 1)
	_, err := io.ReadFull(ts.r, buf)
	if err != nil {
		return 0, err
	}
	c := buf[0]
	if c == '\n' {
		ts.p.Line++
		ts.p.Column = 1
	} else {
		ts.p.Column++
	}
	return c, nil
}

func (ts *TokenStream) Next() (Token, error) {
	for {
		c, err := ts.read()
		if err != nil {
			return Token{}, err
		}
		switch {
		case isSpace(c):
			// Skip whitespace
			continue
		case c == '[':
			// Tokenize quotation
			depth := 1
			start := ts.p
			var v []byte
			for {
				c, err := ts.read()
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					return Token{}, err
				}
				if c == '[' {
					depth++
				} else if c == ']' {
					depth--
				}
				if depth == 0 {
					break
				}
				v = append(v, c)
			}
			if depth > 0 {
				return Token{}, SyntaxError{Position: ts.p, Message: "missing closing ']'"}
			}
			return Token{Position: start, Type: TokenTypeQuotation, Value: string(v)}, nil
		case c == '"', c == '\'':
			// Tokenize literal text
			start := ts.p
			quote := c
			var v []byte
			for {
				c, err := ts.read()
				if errors.Is(err, io.EOF) || c == quote {
					break
				} else if err != nil {
					return Token{}, err
				}
				v = append(v, c)
			}
			return Token{Position: start, Type: TokenTypeLiteralText, Value: string(v)}, nil
		case c >= 33 && c <= 126 && c != '[' && c != ']':
			// Tokenize word
			start := ts.p
			v := []byte{c}
			for {
				c, err := ts.read()
				if errors.Is(err, io.EOF) || isSpace(c) {
					break
				} else if err != nil {
					return Token{}, err
				}
				v = append(v, c)
			}
			return Token{Position: start, Type: TokenTypeWord, Value: string(v)}, nil
		default:
			return Token{}, SyntaxError{Position: ts.p, Message: fmt.Sprintf("unexpected character %q", c)}
		}
	}
}

func isSpace(c byte) bool { return c == ' ' || c == '\n' || c == '\t' }
