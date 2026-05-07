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
    Hello, World!`

var directives = `@open
{{- define "$1" -}}
<!DOCTYPE html>
<html lang="en">

@close
</html>
{{ end }}`

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
	p.DebugPrint(root)
}

func TestParseRawDirectives(t *testing.T) {
	tokens := lexer.Lex(directives)
	for _, tk := range tokens {
		t.Log(tk)
	}
	p := parser.New()
	root := p.Parse(tokens)
	p.DebugPrint(root)
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
