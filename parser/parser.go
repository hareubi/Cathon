package parser

import (
	"cathon/ast"
	"cathon/lexer"
	"cathon/token"
	"fmt"
)

func New(lexerP *lexer.Lexer) Parser {
	p := Parser{
		l:      lexerP,
		errors: []string{},
	}
	p.NextToken()
	p.NextToken()

	return p
}

type Parser struct {
	l *lexer.Lexer

	errors    []string
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
	case token.RETURN:
		return parserP.ParseReturnStatement()
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

	if !parserP.expectPeek(token.ASSIGN) {
		return nil
	}

	for !parserP.curTokenIs(token.SEMICOLON) {
		parserP.NextToken()
	}

	return stmt
}

func (parserP *Parser) ParseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: parserP.curToken}

	if !parserP.expectPeek(token.INT) {
		return nil
	}

	stmt.ReturnValue = &ast.Identifier{Token: parserP.curToken, Value: parserP.curToken.Literal}

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
		parserP.PeekError(tokenB)
		return false
	}
}

func (parserP *Parser) PeekError(tokenType token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", tokenType, parserP.peekToken.Type)
	parserP.errors = append(parserP.errors, msg)
}
func (parserP *Parser) Errors() []string {
	return parserP.errors
}
