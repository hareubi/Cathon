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

func (l *Lexer) PeekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return []rune(l.input)[l.readPosition]
	}
}

func NewToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.SkipWhiteSpace()

	switch l.ch {
	case '=':
		if l.PeekChar() == '=' {
			ch := l.ch
			l.ReadChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = NewToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = NewToken(token.PLUS, l.ch)
	case '-':
		tok = NewToken(token.MINUS, l.ch)
	case '!':
		if l.PeekChar() == '=' {
			ch := l.ch
			l.ReadChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOTEQ, Literal: literal}
		} else {
			tok = NewToken(token.BANG, l.ch)
		}
	case '/':
		tok = NewToken(token.SLASH, l.ch)
	case '*':
		tok = NewToken(token.ASTERISK, l.ch)
	case '<':
		tok = NewToken(token.LT, l.ch)
	case '>':
		tok = NewToken(token.GT, l.ch)
	case ';':
		tok = NewToken(token.SEMICOLON, l.ch)
	case ':':
		tok = NewToken(token.COLON, l.ch)
	case ',':
		tok = NewToken(token.COMMA, l.ch)
	case '{':
		tok = NewToken(token.LBRACE, l.ch)
	case '}':
		tok = NewToken(token.RBRACE, l.ch)
	case '(':
		tok = NewToken(token.LPAREN, l.ch)
	case ')':
		tok = NewToken(token.RPAREN, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.ReadString()
	case '[':
		tok = NewToken(token.LBRACKET, l.ch)
	case ']':
		tok = NewToken(token.RBRACKET, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if unicode.IsLetter(l.ch) {
			tok.Literal = l.ReadIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if unicode.IsDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.ReadNumber()
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

func (l *Lexer) ReadString() string {
	position := l.position + 1
	for {
		l.ReadChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}