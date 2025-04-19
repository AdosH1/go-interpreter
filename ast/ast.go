package ast

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Node interface {
	TokenLiteral() string
	String() string
}
