package parser

import (
	"interpreter/ast"
	"interpreter/lexer"
	"strconv"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)
	if program == nil {
		t.Fatalf("program is nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements has wrong length. exptected=%d, got=%d", 3, len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "let" {
		t.Fatalf("statement.TokenLiteral not 'let'. got=%q", statement.TokenLiteral())
		return false
	}

	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("statement not *ast.LetStatement. got=%T", statement)
		return false
	}

	if letStatement.Name.Value != name {
		t.Errorf("letStatement.Name.Value not %q. got=%q", name, letStatement.Name.Value)
		return false
	}

	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("letStatement.Name.TokenLiteral not %q. got=%q", name, letStatement.Name.TokenLiteral())
		return false
	}

	return true
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements has wrong length. got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	identifier, ok := statement.Value.(*ast.Identifier)
	if !ok {
		t.Fatalf("statement.Value is not ast.Identifier. got=%T", statement.Value)
	}

	if identifier.Value != "foobar" {
		t.Errorf("identifier.Value not %q. got=%q", "foobar", identifier.Value)
	}
	if identifier.TokenLiteral() != "foobar" {
		t.Errorf("identifier.TokenLiteral not %q. got=%q", "foobar", identifier.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements has wrong length. got=%d", len(program.Statements))
	}
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := statement.Value.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("statement.Value is not ast.IntegerLiteral. got=%T", statement.Value)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %q. got=%q", "5", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integarValue int64
	}{
		{"!5;", "!", 5},
		{"!15;", "!", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements has wrong length. got=%d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		expression, ok := statement.Value.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("statement.Value is not ast.PrefixExpression. got=%T", statement.Value)
		}
		if expression.Operator != tt.operator {
			t.Errorf("expression.Operator not %q. got=%q", tt.operator, expression.Operator)
		}
		if !testIntegerLiteral(t, expression.RightExpression, tt.integarValue) {
			t.Errorf("IntegerLiteral not %d. got=%d", tt.integarValue, expression.RightExpression)
			return
		}
	}
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input         string
		expectedLeft  int64
		operator      string
		expectedRight int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements has wrong length. got=%d", len(program.Statements))
		}
		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		expression, ok := statement.Value.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("statement.Expression is not ast.InfixExpression. got=%T", statement)
		}

		if !testIntegerLiteral(t, expression.Left, tt.expectedLeft) {
			return
		}
		if expression.Operator != tt.operator {
			t.Errorf("expression.Operator not %q. got=%q", tt.operator, expression.Operator)
		}
		if !testIntegerLiteral(t, expression.Right, tt.expectedRight) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}
	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Statement not *ast.ReturnStatement. got=%T", statement)
			continue
		}
		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("returnStatement.TokenLiteral not 'return, got '%q'", returnStatement.TokenLiteral())
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser encountered %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error: %s", msg)
	}
	t.FailNow()
}

func testIntegerLiteral(t *testing.T, integerLiteral ast.Expression, expectedValue int64) bool {
	integer, ok := integerLiteral.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("IntegerLiteral not ast.IntegerLiteral. got=%T", integerLiteral)
		return false
	}

	if integer.Value != expectedValue {
		t.Errorf("IntegerLiteral not %d. got=%d", expectedValue, integer.Value)
		return false
	}

	if integer.TokenLiteral() != strconv.FormatInt(expectedValue, 10) {
		t.Errorf("IntegerLiteral not %q. got=%q", strconv.FormatInt(expectedValue, 10), integer.TokenLiteral())
		return false
	}
	return true
}
