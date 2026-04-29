package jtml

import (
	"testing"

	"github.com/jasontconnell/jtml/lexer"
	"github.com/jasontconnell/jtml/parser"
	"github.com/jasontconnell/jtml/process"
)

var home = `@jtml

 #head
 
 #body
  #menu
  #h1 Welcome!
  #posts
`

func TestJTML(t *testing.T) {
	tokens := lexer.Lex(home)
	for _, tk := range tokens {
		t.Log(tk.Type, tk.Value, tk.Level, tk.LineNum)
	}
}

func TestParse(t *testing.T) {
	tokens := lexer.Lex(home)
	p := parser.New()
	root := p.Parse(tokens)
	for _, n := range root.GetChildren() {
		t.Logf("%T %s", n, n.TokenLiteral())
	}
}

func TestParseTemplates(t *testing.T) {
	tlist, err := process.ParseTemplates("./tmpl")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	for _, tmpl := range tlist {
		t.Log(tmpl.Name, len(tmpl.Nodes))
		for _, n := range tmpl.Nodes {
			t.Logf("%T", n)
		}
	}
}

func TestGenerate(t *testing.T) {
	tlist, err := process.ParseTemplates("./tmpl")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	_, err = process.ProcessTemplates(tlist)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
