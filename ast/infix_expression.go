package ast

import (
	"bytes"
	"interpreter/token"
)

var _ Expression = (*InfixExpression)(nil)

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (i InfixExpression) TokenLiteral() string { return i.Token.Literal }

func (i InfixExpression) expressionNode() {}

func (i InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")

	return out.String()
}
