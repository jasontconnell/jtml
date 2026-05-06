package jtml

import (
	"testing"

	"github.com/jasontconnell/jtml/lexer"
	"github.com/jasontconnell/jtml/parser"
	"github.com/jasontconnell/jtml/process"
)

var home = `
#jtml Main
 #head
  #css
  #js
 
 #body body
  #container 
    Hello, World!
`

func TestJTML(t *testing.T) {
	tokens := lexer.Lex(home)
	for _, tk := range tokens {
		t.Log(tk.Type, tk.Value, tk.Level, tk.LineNum)
	}
}

func TestParseRaw(t *testing.T) {
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
		t.Log(tmpl.Name)
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
