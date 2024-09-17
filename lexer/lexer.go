package lexer

import (
	"cathon/token"
	"unicode"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           rune
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.ReadChar()
	return l
}

func (l *Lexer) ReadChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = []rune(l.input)[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func NewToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.SkipWhiteSpace()

	switch l.ch {
	case '=':
		tok = NewToken(token.TokenType(token.ASSIGN), '=')
	case ';':
		tok = NewToken(token.TokenType(token.SEMICOLON), ';')
	case '(':
		tok = NewToken(token.TokenType(token.LPAREN), '(')
	case ')':
		tok = NewToken(token.TokenType(token.RPAREN), ')')
	case ',':
		tok = NewToken(token.TokenType(token.COMMA), ',')
	case '+':
		tok = NewToken(token.TokenType(token.PLUS), '+')
	case '{':
		tok = NewToken(token.TokenType(token.LBRACE), '{')
	case '}':
		tok = NewToken(token.TokenType(token.RBRACE), '}')
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if unicode.IsLetter(l.ch) {
			tok.Literal = l.ReadIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if unicode.IsDigit(l.ch) {
			tok.Literal = l.ReadNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = NewToken(token.ILLEGAL, l.ch)
		}
	}
	l.ReadChar()
	return tok
}

func (l *Lexer) ReadIdentifier() string {
	position := l.position
	for unicode.IsLetter(l.ch) {
		l.ReadChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) ReadNumber() string {
	position := l.position
	for unicode.IsDigit(l.ch) {
		l.ReadChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) SkipWhiteSpace() {
	for unicode.IsSpace(l.ch) {
		l.ReadChar()
	}
}
