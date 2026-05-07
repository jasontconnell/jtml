package parser

import "fmt"

type NodeType int

const (
	Raw NodeType = iota
	Directive
	Include
	Stream
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
	case Stream:
		s = "Stream"
	case Root:
		s = "Root"
	}
	return s
}

type Node interface {
	TokenLiteral() string
	GetChildren() []Node
	GetParameters() []Node
	GetType() NodeType
	GetDepth() int
	GetIndex() int
	GetEndline() bool

	String() string
}

type node struct {
	index      int
	raw        string
	parameters []*node
	children   []*node
	nodeType   NodeType
	depth      int
	endline    bool
}

func newNode(nodeType NodeType, raw string, parameters []*node, depth int, endline bool) *node {
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
		prms += fmt.Sprintf("[%d: %s] ", p.index, p.raw)
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

func (n *node) GetParameters() []Node {
	list := []Node{}
	for _, p := range n.parameters {
		list = append(list, p)
	}
	return list
}

func (n *node) GetDepth() int {
	return n.depth
}

func (n *node) GetIndex() int {
	return n.index
}

func (n *node) GetEndline() bool {
	return n.endline
}
