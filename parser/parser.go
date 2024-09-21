package parser

import (
	"cathon/ast"
	"cathon/lexer"
	"cathon/token"
)

func New(lexerP *lexer.Lexer) Parser {
	p := Parser{l: lexerP}
	p.NextToken()
	p.NextToken()

	return p
}

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
}

func (p *Parser) NextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (parserP *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for parserP.curToken.Type != token.EOF {
		stmt := parserP.ParseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		parserP.NextToken()
	}
	return program
}

func (parserP *Parser) ParseStatement() ast.Statement {
	switch parserP.curToken.Type {
	case token.LET:
		return parserP.ParseLetStatement()
	default:
		return nil
	}
}

func (parserP *Parser) ParseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: parserP.curToken}

	if !parserP.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: parserP.curToken, Value: parserP.curToken.Literal}

	for !parserP.curTokenIs(token.SEMICOLON) {
		parserP.NextToken()
	}

	return stmt
}

func (parserP *Parser) curTokenIs(tokenB token.TokenType) bool {
	return parserP.curToken.Type == tokenB
}
func (parserP *Parser) peekTokenIs(tokenB token.TokenType) bool {
	return parserP.peekToken.Type == tokenB
}
func (parserP *Parser) expectPeek(tokenB token.TokenType) bool {
	if parserP.peekTokenIs(tokenB) {
		parserP.NextToken()
		return true
	} else {
		return false
	}
}