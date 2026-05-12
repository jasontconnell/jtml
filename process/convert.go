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
			Depth:    n.GetDepth(),
		}
	case parser.Include:
		tn = data.Include{
			Name:       n.TokenLiteral(),
			Parameters: convertParameters(n.GetParameters()),
			Children:   convertNodes(n.GetChildren()),
			Depth:      n.GetDepth(),
		}
	case parser.Root:
		tn = data.Root{
			Children: convertNodes(n.GetChildren()),
			Depth:    n.GetDepth(),
		}
	case parser.Raw:
		tn = data.Raw{
			Value:   n.TokenLiteral(),
			Depth:   n.GetDepth(),
			Endline: n.GetEndline(),
		}
	case parser.Stream:
		tn = data.Stream{
			Stream: convertNodes(n.GetParameters()),
			Depth:  n.GetDepth(),
		}
	}
	return tn
}

func convertParameters(plist []parser.Node) []data.Parameter {
	dplist := []data.Parameter{}
	for _, p := range plist {
		dp := data.Parameter{
			Index:   p.GetIndex(),
			Value:   p.TokenLiteral(),
			Endline: p.GetEndline(),
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
