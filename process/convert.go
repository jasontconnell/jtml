package process

import (
	"github.com/jasontconnell/jtml/data"
	"github.com/jasontconnell/jtml/parser"
)

type rootNode struct {
	Node parser.Node
	Name string
}

func toTemplate(r rootNode) data.Template {
	t := data.Template{Name: r.Name}
	return t
}

func convertNode(n parser.Node) data.TemplateNode {
	var tn data.TemplateNode
	switch n.GetType() {
	case parser.Directive:
		tn = data.Directive{
			Name:       n.TokenLiteral(),
			Parameters: convertParameters(n.GetParameters()),
			Children:   convertNodes(n.GetChildren()),
		}
	}
	return tn
}

func convertParameters(plist []parser.Parameter) []data.Parameter {
	dplist := []data.Parameter{}
	for _, p := range plist {
		dp := data.Parameter{
			Index: p.GetIndex(),
			Value: p.GetValue(),
		}
		dplist = append(dplist, dp)
	}
	return dplist
}

func convertNodes(ns []parser.Node) []data.TemplateNode {
	tns := []data.TemplateNode{}
	for _, n := range ns {
		tn := convertNode(n)
		tns = append(tns, tn)
	}
	return tns
}
