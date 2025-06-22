package ast

import "interpreter/token"

var _ Expression = (*PrefixExpression)(nil)

type PrefixExpression struct {
	Token           token.Token
	Operator        string
	RightExpression Expression
}

func (p *PrefixExpression) expressionNode()      {}
func (p *PrefixExpression) TokenLiteral() string { return p.Token.Literal }
func (p *PrefixExpression) String() string {
	return "(" + p.Operator + p.RightExpression.String() + ")"
}
