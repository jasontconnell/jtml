package parser

import (
	"fmt"
	"strings"

	"github.com/jasontconnell/jtml/lexer"
)

type parser struct {
}

func New() *parser {
	return new(parser)
}

func (p *parser) DebugPrint(r Node) {
	fmt.Println(strings.Repeat(" ", r.GetDepth()), r.TokenLiteral())
	c := r.GetChildren()
	if len(c) > 0 {
		for _, child := range c {
			p.DebugPrint(child)
		}
	}
}

func (p *parser) Parse(tokens []lexer.Token) Node {
	root := newNode(Root, "", nil, 0)
	p.recurseParse(tokens, root, 0, 0)

	return root
}

func (p *parser) recurseParse(tokens []lexer.Token, cur *node, idx int, depth int) int {
	start := idx
	for idx < len(tokens) {
		tk := tokens[idx]

		switch tk.Type {
		case lexer.Raw:
			st, ridx := p.consumeRawTokens(tokens, idx)
			val := strings.Join(st, " ")

			n := newNode(Raw, val, nil, depth)
			cur.children = append(cur.children, n)
			idx = ridx + 1
		case lexer.Parameter:
			idx++
		case lexer.Include:
			prms := p.getParameters(tokens, idx+1, tk.Level)
			idx += len(prms) + 1

			n := newNode(Include, tk.Value, prms, depth)
			cur.children = append(cur.children, n)

			recurse := p.hasChildren(tokens, idx, tk.Level)
			if recurse {
				nc := p.recurseParse(tokens, n, idx, depth+1)
				idx += nc + len(n.children) // adjust idx since we parsed those already
			}
		case lexer.Directive:
			ss, didx := p.consumeRawTokens(tokens, idx+1)
			n := newNode(Directive, tk.Value, []parameter{{index: 0, value: strings.Join(ss, " ")}}, tk.Level)
			cur.children = append(cur.children, n)
			idx = didx
		}
	}
	return idx - start
}

func (p *parser) consumeRawTokens(tokens []lexer.Token, idx int) ([]string, int) {
	st := []string{}
	raws := tokens[idx].Type == lexer.Raw
	for raws && idx < len(tokens) {
		stk := tokens[idx]
		st = append(st, stk.Value)
		raws = false
		if idx+1 < len(tokens) {
			raws = tokens[idx+1].Type == lexer.Raw
			if raws {
				idx++
			}
		}
	}

	return st, idx
}

func (p *parser) getParameters(tokens []lexer.Token, idx, level int) []parameter {
	prms := []parameter{}
	for i := idx; i < len(tokens); i++ {
		tk := tokens[i]
		if tk.Level != level || tk.Type != lexer.Parameter {
			break
		}

		p := parameter{
			index: len(prms),
			value: tk.Value,
		}
		prms = append(prms, p)

	}
	return prms
}

func (p *parser) hasChildren(tokens []lexer.Token, startIndex int, level int) bool {
	if startIndex > len(tokens) {
		return false
	}
	hasSubs := false
	for i := startIndex; i < len(tokens); i++ {
		tk := tokens[i]

		if tk.Level == level+1 {
			hasSubs = true
			break
		} else if tk.Level <= level {
			hasSubs = false
			break
		}
	}
	return hasSubs
}
