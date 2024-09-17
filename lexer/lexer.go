package lexer

import "cathon/token"

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

	switch l.ch {
	case '=':
		tok = NewToken(token.TokenType(token.ASSIGN), '=')
	case ';':
		tok = NewToken(token.TokenType(token.ASSIGN), ';')
	case '(':
		tok = NewToken(token.TokenType(token.ASSIGN), '(')
	case ')':
		tok = NewToken(token.TokenType(token.ASSIGN), ')')
	case ',':
		tok = NewToken(token.TokenType(token.ASSIGN), ',')
	case '+':
		tok = NewToken(token.TokenType(token.ASSIGN), '+')
	case '{':
		tok = NewToken(token.TokenType(token.ASSIGN), '{')
	case '}':
		tok = NewToken(token.TokenType(token.ASSIGN), '}')
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	}
	l.ReadChar()
	return tok
}
