package parser

import "fmt"

type NodeType int

const (
	Raw NodeType = iota
	Directive
	Include
	Parameter
	Root
	Endline
	Indent
)

func (nt NodeType) String() string {
	s := "Undefined"
	switch nt {
	case Raw:
		s = "Raw"
	case Directive:
		s = "Directive"
	case Include:
		s = "Include"
	case Parameter:
		s = "Parameter"
	case Root:
		s = "Root"
	case Indent:
		s = "Indent"
	case Endline:
		s = "Endline"
	}
	return s
}

type Node interface {
	TokenLiteral() string
	GetChildren() []Node
	GetParameters() []Node
	GetType() NodeType
	GetIndex() int

	String() string
}

type node struct {
	index      int
	raw        string
	parameters []*node
	children   []*node
	nodeType   NodeType
}

func newNode(nodeType NodeType, raw string, children, parameters []*node) *node {
	return &node{
		nodeType:   nodeType,
		raw:        raw,
		parameters: parameters,
		children:   children,
	}
}

func (n *node) String() string {
	prms := ""
	for _, p := range n.parameters {
		prms += fmt.Sprintf("[%d: %s] ", p.index, p.raw)
	}
	s := fmt.Sprintf("[%s %s %s [children: %d]]", n.raw, n.nodeType, prms, len(n.children))
	return s
}

func (n *node) TokenLiteral() string {
	return n.raw
}

func (n *node) GetChildren() []Node {
	list := []Node{}
	for _, c := range n.children {
		list = append(list, c)
	}
	return list
}

func (n *node) GetType() NodeType {
	return n.nodeType
}

func (n *node) GetParameters() []Node {
	list := []Node{}
	for _, p := range n.parameters {
		list = append(list, p)
	}
	return list
}

func (n *node) GetIndex() int {
	return n.index
}
