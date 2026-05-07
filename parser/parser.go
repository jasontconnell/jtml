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
	depth := r.GetDepth()
	if depth <= 0 {
		depth = 0
	}

	fmt.Println("---", strings.Repeat(" ", depth), r)

	c := r.GetChildren()
	if len(c) > 0 {
		for _, child := range c {
			p.DebugPrint(child)
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
	for i := 0; i < len(tokens); i++ {
		tk := tokens[i]

		n, num := p.nodeFromToken(tokens, tk, i, tk.Level)

		cur, ok := stack.Peek()
		if !ok {
			log.Fatal("peek on empty stack")
			return
		}
		log.Println(tk, "adding child to node", n, cur)
		cur.children = append(cur.children, n)

		next, hasNext := p.nextTokenNode(tokens, i+1)

		if hasNext && next.Level > tk.Level {
			stack.Push(n)
		} else if hasNext {
			for j := 0; j < tk.Level-next.Level; j++ {
				stack.Pop()
			}
		}

		i += num
	}
}

func (p *parser) nodeFromToken(tokens []lexer.Token, tk lexer.Token, idx, depth int) (*node, int) {
	var n *node
	consumed := 0
	switch tk.Type {
	case lexer.Raw:
		rawval, num := p.consumeRawTokens(tokens, idx)
		n = newNode(Raw, rawval, nil, depth, tk.IsEndline)
		consumed = num
	case lexer.Include:
		prms := p.getParameters(tokens, idx+1, tk.Level)
		n = newNode(Include, tk.Value, prms, tk.Level, tk.IsEndline)
		consumed = len(prms)
	case lexer.Directive:
		rawval, num := p.consumeRawTokens(tokens, idx+1)
		n = newNode(Directive, tk.Value, []parameter{{index: 0, value: rawval}}, tk.Level, tk.IsEndline)
		consumed = num
	}
	return n, consumed
}

func (p *parser) nextTokenNode(tokens []lexer.Token, start int) (lexer.Token, bool) {
	var tk lexer.Token
	var found bool
	for i := start; i < len(tokens) && !found; i++ {
		switch tokens[i].Type {
		case lexer.Parameter:
			continue
		case lexer.Include, lexer.Directive, lexer.Raw:
			tk = tokens[i]
			found = true
		}
	}

	return tk, found
}

func (p *parser) consumeRawTokens(tokens []lexer.Token, start int) (string, int) {
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

	rawval := ""
	for _, s := range st {
		rawval += s.Value
		if s.IsEndline {
			rawval += "\n"
		} else {
			rawval += " "
		}
	}
	rawval = strings.TrimRight(rawval, " ")

	return rawval, len(st)
}

func (p *parser) getParameters(tokens []lexer.Token, start, level int) []parameter {
	prms := []parameter{}
	for i := start; i < len(tokens); i++ {
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
