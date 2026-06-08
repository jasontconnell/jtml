package parser

import (
	"fmt"
	"log"
	"strings"

	"github.com/jasontconnell/collections"
	"github.com/jasontconnell/jtml/lexer"
)

type parser struct {
}

func New() *parser {
	return new(parser)
}

func (p *parser) DebugPrint(r Node) {
	printNode(r, 0)
}

func printNode(r Node, d int) {
	fmt.Println("---", strings.Repeat(" ", d), r)
	c := r.GetChildren()
	if len(c) > 0 {
		for _, child := range c {
			printNode(child, d+1)
		}
	}
}

func (p *parser) Parse(tokens []lexer.Token) Node {
	root := newNode(Root, "", nil, nil)

	stack := collections.NewStack[*node]()
	stack.Push(root)

	p.parse(tokens, stack)

	return root
}

func (p *parser) parse(tokens []lexer.Token, stack collections.Stack[*node]) {
	i := 0
	for i < len(tokens) {
		if tokens[i].Type == lexer.Endline {
			i++
			continue
		}
		indents := p.consumeIndents(tokens, i)
		i += indents
		tk := tokens[i]
		n, num := p.nodeFromToken(tokens, tk, i)
		i += num

		cur, ok := stack.Peek()
		if !ok {
			log.Fatal("peek on empty stack")
			return
		}
		cur.children = append(cur.children, n)

		nextLine, nextLineLevel, hasNextLine := p.nextLine(tokens, i)
		if hasNextLine {
			i = nextLine
			if nextLineLevel > indents {
				stack.Push(n)
			} else if indents > nextLineLevel {
				lv := indents
				for lv > nextLineLevel {
					stack.Pop()
					lv--
				}
			}
		} else {
			break
		}
	}

	for stack.Any() {
		stack.Pop()
	}
}

func (p *parser) nodeFromToken(tokens []lexer.Token, tk lexer.Token, idx int) (*node, int) {
	var n *node
	consumed := 0
	switch tk.Type {
	case lexer.Raw:
		rawval, num := p.consumeTokens(tokens, idx, func(tk lexer.Token) bool {
			return tk.Type == lexer.Raw
		}, nil)
		n = rawval
		consumed = num
	case lexer.Include:
		prms := p.getParameters(tokens, idx+1)
		n = newNode(Include, tk.Value, nil, prms)
		consumed = len(prms) + 1
	case lexer.Directive:
		rawval, num := p.consumeTokens(tokens, idx+1, func(tk lexer.Token) bool {
			return tk.Type != lexer.Directive
		}, func(tk lexer.Token) bool {
			return tk.Type == lexer.Endline
		})
		n = newNode(Directive, tk.Value, []*node{rawval}, nil)
		consumed = num
	}
	return n, consumed
}

func (p *parser) nextLine(tokens []lexer.Token, start int) (int, int, bool) {
	var idx int
	var level int
	var found bool
	for i := start; i < len(tokens); i++ {
		tk := tokens[i]
		if tk.Type == lexer.Endline && len(tokens) > i+1 {
			found = true
			idx = i + 1
			break
		}
	}
	if found {
		for i := idx; i < len(tokens); i++ {
			if tokens[i].Type != lexer.Indent {
				break
			}
			level++
		}
	}
	return idx, level, found
}

func (p *parser) consumeIndents(tokens []lexer.Token, start int) int {
	j := 0
	for i := start; i < len(tokens); i++ {
		if tokens[i].Type != lexer.Indent {
			break
		}
		j++
	}
	return j
}

func (p *parser) consumeTokens(tokens []lexer.Token, start int, check func(tk lexer.Token) bool, skip func(tk lexer.Token) bool) (*node, int) {
	st := []lexer.Token{}
	consume := check(tokens[start])
	consumed := 0
	for consume && start+consumed < len(tokens) {
		stk := tokens[start+consumed]
		if skip == nil || !skip(stk) {
			st = append(st, stk)
		}

		consume = false
		if start+consumed+1 < len(tokens) {
			consume = check(tokens[start+consumed+1])
		}
		consumed++
	}

	nodes := []*node{}
	for _, tk := range st {
		val := tk.Value
		if tk.Type == lexer.Endline {
			val = "\r\n"
		}
		n := newNode(Raw, val, nil, nil)
		nodes = append(nodes, n)
	}

	n := newNode(Raw, p.joinRaws(nodes), nil, nil)
	return n, consumed
}

func (p *parser) joinRaws(raws []*node) string {
	str := ""
	for _, r := range raws {
		str += r.raw + " "
	}
	return strings.TrimRight(str, " ")
}

func (p *parser) getParameters(tokens []lexer.Token, start int) []*node {
	prms := []*node{}
	for i := start; i < len(tokens); i++ {
		tk := tokens[i]
		if tk.Type != lexer.Parameter {
			break
		}

		p := &node{
			nodeType: Parameter,
			index:    len(prms),
			raw:      tk.Value,
		}
		prms = append(prms, p)
	}
	return prms
}
