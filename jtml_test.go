package jtml

import (
	"testing"

	"github.com/jasontconnell/jtml/lexer"
	"github.com/jasontconnell/jtml/parser"
)

var home = `#jtml Main
 #head
  #css
  #js
 
 #body body
  #container 
   Hello, World!
  #footer
`

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

func TestParseFull(t *testing.T) {
	tokens := lexer.Lex(home)
	p := parser.New()
	root := p.Parse(tokens)
	p.DebugPrint(root)

	for _, tk := range tokens {
		t.Log(tk)
	}
}

func TestParseDirectives(t *testing.T) {
	tokens := lexer.Lex(directives)
	for _, tk := range tokens {
		t.Log(tk)
	}
	p := parser.New()
	root := p.Parse(tokens)
	p.DebugPrint(root)

	for _, tk := range tokens {
		t.Log(tk)
	}
}
