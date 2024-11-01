package parser

import (
	"cathon/ast"
	"cathon/lexer"
	"cathon/token"
	"fmt"
	"strconv"
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
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(left ast.Expression) ast.Expression
)
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOTEQ:    EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN: CALL,
}

func New(lexerP *lexer.Lexer) *Parser {
	p := &Parser{
		l:      lexerP,
		errors: []string{},
	}
	p.NextToken()
	p.NextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.ParseIdentifier)
	p.registerPrefix(token.INT, p.ParseIntegerLiteral)
	p.registerPrefix(token.BANG, p.ParsePrefixExpression)
	p.registerPrefix(token.MINUS, p.ParsePrefixExpression)
	p.registerPrefix(token.TRUE, p.ParseBool)
	p.registerPrefix(token.FALSE, p.ParseBool)
	p.registerPrefix(token.IF, p.ParseIfExpression)
	p.registerPrefix(token.LPAREN, p.ParseGroupedExpression)
	p.registerPrefix(token.FUNCTION, p.ParseFunctionExpression)
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.ParseInfixExpression)
	p.registerInfix(token.MINUS, p.ParseInfixExpression)
	p.registerInfix(token.ASTERISK, p.ParseInfixExpression)
	p.registerInfix(token.SLASH, p.ParseInfixExpression)
	p.registerInfix(token.EQ, p.ParseInfixExpression)
	p.registerInfix(token.NOTEQ, p.ParseInfixExpression)
	p.registerInfix(token.LT, p.ParseInfixExpression)
	p.registerInfix(token.GT, p.ParseInfixExpression)
	p.registerInfix(token.LPAREN, p.ParseCallExpression)

	return p
}
func (p *Parser) NextToken() {
	defer untrace(trace("NextToken"))
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
	fmt.Printf("%s %p\n", p.curToken.Literal, p)
}
func (parserP *Parser) ParseIdentifier() ast.Expression {
	defer untrace(trace("ParseIdentifier"))
	return &ast.Identifier{Token: parserP.curToken, Value: parserP.curToken.Literal}
}
func (parserP *Parser) ParseIntegerLiteral() ast.Expression {
	defer untrace(trace("ParseIntegerLiteral"))
	lit := &ast.IntegerLiteral{Token: parserP.curToken}

	value, err := strconv.ParseInt(parserP.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", parserP.curToken.Literal)
		parserP.errors = append(parserP.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}
func (parserP *Parser) ParseBool() ast.Expression {
	defer untrace(trace("ParseBool: " + parserP.curToken.Literal))
	fmt.Printf("%p\n", parserP)
	return &ast.Boolean{Token: parserP.curToken, Value: parserP.CurTokenIs(token.TRUE)}
}
func (parserP *Parser) ParseGroupedExpression() ast.Expression {
	defer untrace(trace("ParseGroupedExpression"))
	parserP.NextToken()
	exp := parserP.ParseExpression(LOWEST)
	if !parserP.ExpectPeek(token.RPAREN) {
		return nil
	}

	return exp
}
func (parserP *Parser) ParseIfExpression() ast.Expression {
	defer untrace(trace("ParseIfExpression"))

	exp := &ast.IfExpression{Token: parserP.curToken}

	if !parserP.ExpectPeek(token.LPAREN) {
		return nil
	}

	parserP.NextToken()
	exp.Condition = parserP.ParseExpression(LOWEST)

	if !parserP.ExpectPeek(token.RPAREN) || !parserP.ExpectPeek(token.LBRACE) {
		return nil
	}
	exp.Consequence = parserP.ParseBlockStatement()

	if parserP.PeekTokenIs(token.ELSE) {
		parserP.NextToken()
		if !parserP.ExpectPeek(token.LBRACE) {
			return nil
		}

		exp.Alternative = parserP.ParseBlockStatement()
	}

	return exp
}
func (parserP *Parser) ParseFunctionExpression() ast.Expression {
	exp := &ast.FunctionLiteral{Token: parserP.curToken}

	if !parserP.ExpectPeek(token.LPAREN){
		return nil
	}
	exp.Parameters = parserP.ParseFunctionParameters()

	if !parserP.ExpectPeek(token.LBRACE){
		return nil
	}
	exp.Body = parserP.ParseBlockStatement()

	return exp
}
func (parserP*Parser) ParseFunctionParameters() []*ast.Identifier{
	identifiers := []*ast.Identifier{}

	if parserP.PeekTokenIs(token.RPAREN) {
	parserP.NextToken()
	return identifiers
	}
	parserP.NextToken()

	ident := &ast.Identifier{Token: parserP.curToken, Value: parserP.curToken.Literal}
	identifiers = append(identifiers, ident)

	for parserP.PeekTokenIs(token.COMMA) {
		parserP.NextToken()
		parserP.NextToken()
		ident:= &ast.Identifier{Token: parserP.curToken, Value: parserP.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !parserP.ExpectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}
func (parserP *Parser) ParseBlockStatement() *ast.BlockStatement {
	defer untrace(trace("ParseBlockStatement"))
	block := &ast.BlockStatement{Token: parserP.curToken}
	block.Statements = []ast.Statement{}

	parserP.NextToken()

	for !parserP.CurTokenIs(token.RBRACE) && !parserP.CurTokenIs(token.EOF) {
		stmt := parserP.ParseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		parserP.NextToken()
	}
	return block
}
func (parserP *Parser) ParsePrefixExpression() ast.Expression {
	defer untrace(trace("ParsePrefixExpression"))
	expression := &ast.PrefixExpression{
		Token:    parserP.curToken,
		Operator: parserP.curToken.Literal,
	}

	parserP.NextToken()

	expression.Right = parserP.ParseExpression(PREFIX)

	return expression
}
func (parserP *Parser) ParseInfixExpression(left ast.Expression) ast.Expression {
	defer untrace(trace("ParseInfixExpression"))
	expression := &ast.InfixExpression{
		Token:    parserP.curToken,
		Operator: parserP.curToken.Literal,
		Left:     left,
	}
	precedence := parserP.CurPrecedence()
	parserP.NextToken()
	expression.Right = parserP.ParseExpression(precedence)

	return expression
}
func (parserP *Parser) PeekPrecedence() int {
	if p, ok := precedences[parserP.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}
func (parserP *Parser) CurPrecedence() int {
	if p, ok := precedences[parserP.curToken.Type]; ok {
		return p
	}
	return LOWEST
}
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
	defer untrace(trace("ParseStatement"))
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
	defer untrace(trace("ParseLetStatement"))
	stmt := &ast.LetStatement{Token: parserP.curToken}

	if !parserP.ExpectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: parserP.curToken, Value: parserP.curToken.Literal}

	if !parserP.ExpectPeek(token.ASSIGN) {
		return nil
	}
	parserP.NextToken()
	stmt.Value = parserP.ParseExpression(LOWEST)
	for !parserP.CurTokenIs(token.SEMICOLON) {
		parserP.NextToken()
	}

	return stmt
}
func (parserP *Parser) ParseReturnStatement() *ast.ReturnStatement {
	defer untrace(trace("ParseReturnStatement"))
	stmt := &ast.ReturnStatement{Token: parserP.curToken}

	parserP.NextToken()

	stmt.ReturnValue = parserP.ParseExpression(LOWEST)

	for !parserP.CurTokenIs(token.SEMICOLON) {
		parserP.NextToken()
	}

	return stmt
}
func (parserP *Parser) ParseExpressionStatement() *ast.ExpressionStatement {
	defer untrace(trace("ParseExpressionStatement"))
	stmt := &ast.ExpressionStatement{Token: parserP.curToken}
	stmt.Expression = parserP.ParseExpression(LOWEST)

	if parserP.PeekTokenIs(token.SEMICOLON) {
		parserP.NextToken()
	}

	return stmt
}
func (parserP *Parser) ParseExpression(precedence int) ast.Expression {
	defer untrace(trace("ParseExpression"))
	prefix := parserP.prefixParseFns[parserP.curToken.Type]
	if prefix == nil {
		parserP.RegisterParsePrefixError(parserP.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !parserP.PeekTokenIs(token.SEMICOLON) && precedence < parserP.PeekPrecedence() {
		infix := parserP.infixParseFns[parserP.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		parserP.NextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}
func (parserP *Parser) CurTokenIs(tokenType token.TokenType) bool {
	return parserP.curToken.Type == tokenType
}
func (parserP *Parser) PeekTokenIs(tokenType token.TokenType) bool {
	return parserP.peekToken.Type == tokenType
}
func (parserP *Parser) ExpectPeek(tokenType token.TokenType) bool {
	if parserP.PeekTokenIs(tokenType) {
		parserP.NextToken()
		return true
	} else {
		parserP.PeekError(tokenType)
		return false
	}
}
func (parserP *Parser) PeekError(tokenType token.TokenType) {
		fmt.Printf("wth its not working")
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", tokenType, parserP.peekToken.Type)
	parserP.errors = append(parserP.errors, msg)
}
func (parserP *Parser) Errors() []string {
	return parserP.errors
}
func (parserP *Parser) RegisterParsePrefixError(tokenType token.TokenType) {
	fmt.Printf("wth its not working")
	msg := fmt.Sprintf("no parse prefix function for %s", tokenType)
	parserP.errors = append(parserP.errors, msg)
}
func (parserP *Parser) ParseCallExpression(function ast.Expression) ast.Expression {
	defer untrace(trace("ParseCallExpression"))
	exp := &ast.CallExpression{Token: parserP.curToken, Function: function}
	exp.Arguments = parserP.ParseCallArguments()
	return exp
}
func (parserP *Parser) ParseCallArguments()[]ast.Expression {
	defer untrace(trace("ParseCallArguments"))
	args := []ast.Expression{}

	if parserP.PeekTokenIs(token.RPAREN) {
		parserP.NextToken()
		return args
	}

	parserP.NextToken()
	args = append(args, parserP.ParseExpression(LOWEST))

	for parserP.PeekTokenIs(token.COMMA) {
		parserP.NextToken()
		parserP.NextToken()
		args = append(args, parserP.ParseExpression(LOWEST))
	}

	if !parserP.ExpectPeek(token.RPAREN) {
		return nil
	}

	return args
}