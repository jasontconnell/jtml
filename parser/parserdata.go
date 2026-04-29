package parser

type NodeType int

const (
	Raw NodeType = iota
	Directive
	Include
	Root
)

type Node interface {
	TokenLiteral() string
	GetChildren() []Node
	GetParameters() []Parameter
	GetType() NodeType
	GetDepth() int
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
}

func newNode(nodeType NodeType, raw string, parameters []parameter, depth int) *node {
	return &node{
		nodeType:   nodeType,
		raw:        raw,
		parameters: parameters,
		children:   []*node{},
		depth:      depth,
	}
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
