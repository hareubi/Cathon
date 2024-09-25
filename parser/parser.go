package parser

import (
	"cathon/ast"
	"cathon/lexer"
	"cathon/token"
	"fmt"
)

const (
	_ int = iota
	LOWEST
	EQUALS      //==
	LESSGREATER // <, >
	SUM         //+
	PRODUCT     //*
	PREFIX      //-x, !x
	CALL        //foo()
)

type Parser struct {
	l *lexer.Lexer

	errors    []string
	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(lexerP *lexer.Lexer) Parser {
	p := Parser{
		l:      lexerP,
		errors: []string{},
	}
	p.NextToken()
	p.NextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.ParseIdentifier)

	return p
}
func (p *Parser) NextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}
func (parserP *Parser) ParseIdentifier() ast.Expression {
	return &ast.Identifier{Token: parserP.curToken, Value: parserP.curToken.Literal}
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func() ast.Expression
)

func (parserP *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	parserP.prefixParseFns[tokenType] = fn
}
func (parserP *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	parserP.infixParseFns[tokenType] = fn
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
		return parserP.ParseExpressionStatement()
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
func (parserP *Parser) ParseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: parserP.curToken}
	stmt.Expression = parserP.ParseExpression(LOWEST)

	if parserP.peekTokenIs(token.SEMICOLON) {
		parserP.NextToken()
	}

	return stmt
}
func (parserP *Parser) ParseExpression(precedence int) ast.Expression {
	prefix := parserP.prefixParseFns[parserP.curToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()

	return leftExp
}

func (parserP *Parser) curTokenIs(tokenType token.TokenType) bool {
	return parserP.curToken.Type == tokenType
}
func (parserP *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return parserP.peekToken.Type == tokenType
}
func (parserP *Parser) expectPeek(tokenType token.TokenType) bool {
	if parserP.peekTokenIs(tokenType) {
		parserP.NextToken()
		return true
	} else {
		parserP.PeekError(tokenType)
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
