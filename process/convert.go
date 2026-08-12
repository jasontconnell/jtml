package process

import (
	"github.com/jasontconnell/jtml/data"
	"github.com/jasontconnell/jtml/parser"
)

type rootNode struct {
	Node      parser.Node
	Name      string
	WriteFile bool
}

func toTemplate(r rootNode) data.Template {
	t := data.Template{Name: r.Name, RootNode: convertNode(r.Node), WriteFile: r.WriteFile}
	return t
}

func convertNode(n parser.Node) data.TemplateNode {
	var tn data.TemplateNode
	switch n.GetType() {
	case parser.Directive:
		tn = data.Directive{
			Name:     n.TokenLiteral(),
			Children: convertNodes(n.GetChildren()),
		}
	case parser.Include:
		tn = data.Include{
			Name:       n.TokenLiteral(),
			Parameters: convertParameters(n.GetParameters()),
			Children:   convertNodes(n.GetChildren()),
		}
	case parser.Root:
		tn = data.Root{
			Children: convertNodes(n.GetChildren()),
		}
	case parser.Raw:
		tn = data.Raw{
			Value:    n.TokenLiteral(),
			Children: convertNodes(n.GetChildren()),
		}
	case parser.Indent:
		tn = data.Indent{}
	case parser.Endline:
		tn = data.Endline{Newline: n.String()}
	}
	return tn
}

func convertParameters(plist []parser.Node) []data.Parameter {
	dplist := []data.Parameter{}
	for _, p := range plist {
		dp := data.Parameter{
			Index: p.GetIndex(),
			Value: p.TokenLiteral(),
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
