package parser

import "fmt"

type NodeType int

const (
	Raw NodeType = iota
	Directive
	Include
	Root
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
	case Root:
		s = "Root"
	}
	return s
}

type Node interface {
	TokenLiteral() string
	GetChildren() []Node
	GetParameters() []Parameter
	GetType() NodeType
	GetDepth() int

	String() string
}

type Parameter interface {
	GetIndex() int
	GetValue() string
}

type parameter struct {
	index int
	value string
}

type node struct {
	raw        string
	parameters []parameter
	children   []*node
	nodeType   NodeType
	depth      int
	endline    bool
}

func newNode(nodeType NodeType, raw string, parameters []parameter, depth int, endline bool) *node {
	return &node{
		nodeType:   nodeType,
		raw:        raw,
		parameters: parameters,
		children:   []*node{},
		depth:      depth,
		endline:    endline,
	}
}

func (n *node) String() string {
	prms := ""
	for _, p := range n.parameters {
		prms += fmt.Sprintf("[%d: %s] ", p.index, p.value)
	}
	s := fmt.Sprintf("%s %s %s (%d) [children: %d]", n.raw, n.nodeType, prms, n.depth, len(n.children))
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

func (n *node) GetParameters() []Parameter {
	list := []Parameter{}
	for _, p := range n.parameters {
		list = append(list, p)
	}
	return list
}

func (n *node) GetDepth() int {
	return n.depth
}

func (p parameter) GetIndex() int {
	return p.index
}

func (p parameter) GetValue() string {
	return p.value
}
