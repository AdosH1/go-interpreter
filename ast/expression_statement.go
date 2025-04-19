package ast

import "interpreter/token"

type ExpressionStatement struct {
	Token token.Token
	Value Expression
}

func (e *ExpressionStatement) TokenLiteral() string { return e.Token.Literal }
func (e *ExpressionStatement) statementNode()       {}
func (e *ExpressionStatement) String() string {
	if e.Value != nil {
		return e.Value.String()
	}
	return ""
}
