package ast

import "interpreter/token"

var _ Expression = (*Identifier)(nil)

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) expressionNode()      {}
func (i *Identifier) String() string {
	return i.Value
}
