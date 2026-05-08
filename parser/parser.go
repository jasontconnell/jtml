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
	root := newNode(Root, "", nil, -1, false)

	stack := collections.NewStack[*node]()
	stack.Push(root)

	p.parse(tokens, stack)

	return root
}

func (p *parser) parse(tokens []lexer.Token, stack collections.Stack[*node]) {
	i := 0
	for i < len(tokens) {
		tk := tokens[i]
		n, num := p.nodeFromToken(tokens, tk, i, tk.Level)

		cur, ok := stack.Peek()
		if !ok {
			log.Fatal("peek on empty stack")
			return
		}
		cur.children = append(cur.children, n)

		next, nextIndex, hasNext := p.nextTokenNode(tokens, i+num)

		if hasNext {
			i = nextIndex
			log.Println(next.Level, tk.Level)
			if next.Level > tk.Level {
				log.Println("pushing", n)
				stack.Push(n)
			} else {
				log.Println("maybe pop?", cur, next)
				for cur.depth > next.Level {
					log.Println("popping")
					cur = stack.Pop()
					log.Println(cur)
				}
			}
		} else {
			i += num
		}
	}
}

func (p *parser) nodeFromToken(tokens []lexer.Token, tk lexer.Token, idx, depth int) (*node, int) {
	var n *node
	consumed := 0
	switch tk.Type {
	case lexer.Raw:
		rawval, num := p.consumeRawTokens(tokens, idx)
		n = newNode(Stream, "", rawval, depth, tk.Endline)
		consumed = num
	case lexer.Include:
		prms := p.getParameters(tokens, idx+1, tk.Level)
		n = newNode(Include, tk.Value, prms, tk.Level, tk.Endline)
		consumed = len(prms) + 1
	case lexer.Directive:
		rawval, num := p.consumeRawTokens(tokens, idx+1)
		n = newNode(Directive, tk.Value, rawval, tk.Level, tk.Endline)
		consumed = num + 1
	}
	return n, consumed
}

func (p *parser) nextTokenNode(tokens []lexer.Token, start int) (lexer.Token, int, bool) {
	var tk lexer.Token
	var found bool
	var index int
	for i := start; i < len(tokens) && !found; i++ {
		switch tokens[i].Type {
		case lexer.Parameter:
			continue
		case lexer.Raw:
			tk = tokens[i]
			found = true
			index = i
		case lexer.Include, lexer.Directive:
			tk = tokens[i]
			found = true
			index = i
		}
	}

	return tk, index, found
}

func (p *parser) consumeRawTokens(tokens []lexer.Token, start int) ([]*node, int) {
	st := []lexer.Token{}
	raws := tokens[start].Type == lexer.Raw
	for raws && start+len(st) < len(tokens) {
		stk := tokens[start+len(st)]
		st = append(st, stk)
		raws = false
		if start+len(st) < len(tokens) {
			raws = tokens[start+len(st)].Type == lexer.Raw
		}
	}

	nodes := []*node{}

	for _, tk := range st {
		n := newNode(Raw, tk.Value, nil, tk.Level, tk.Endline)
		nodes = append(nodes, n)
	}

	return nodes, len(st)
}

func (p *parser) getParameters(tokens []lexer.Token, start, level int) []*node {
	prms := []*node{}
	for i := start; i < len(tokens); i++ {
		tk := tokens[i]
		if tk.Level != level || tk.Type != lexer.Parameter {
			break
		}

		p := &node{
			index: len(prms),
			raw:   tk.Value,
		}
		prms = append(prms, p)
	}
	return prms
}
