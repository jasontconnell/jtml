package parser

import (
	"fmt"
	"log"
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

	p.DebugPrint(root)
	return root
}

func (p *parser) recurseParse(tokens []lexer.Token, cur *node, idx int, depth int) int {
	start := idx
	for idx < len(tokens) {
		tk := tokens[idx]

		switch tk.Type {
		case lexer.Raw:
			n := newNode(Raw, tk.Value, nil, depth)
			cur.children = append(cur.children, n)
			idx++
		case lexer.Parameter:
			log.Println("came across parameter, skipping", tk.Type, tk.Value)
			idx++
		case lexer.Include, lexer.Directive:
			log.Println("got", tk.Type, tk.Value, "parent", cur.raw)
			nt := Include
			if tk.Type == lexer.Directive {
				nt = Directive
			}

			prms := p.getParameters(tokens, idx+1, tk.Level)
			idx += len(prms) + 1

			log.Println("got params", tk.Type, tk.Value, len(prms))

			n := newNode(nt, tk.Value, prms, depth)
			cur.children = append(cur.children, n)

			recurse := p.hasChildren(tokens, idx, tk.Level)
			log.Println("has children", tk.Type, tk.Value, recurse)
			if recurse {
				log.Println("recurse", n.raw, "parent", cur.raw)
				nc := p.recurseParse(tokens, n, idx, depth+1)
				idx += nc + len(n.children) // adjust idx since we parsed those already
			} else {
				// log.Println("no children", tk.Type, tk.Value, "next token", tokens[idx].Value, tokens[idx].Type)
			}
		}
	}
	return idx - start
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
