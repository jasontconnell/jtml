package data

type Template struct {
	Name      string
	RootNode  TemplateNode
	IsPartial bool
}

func (t Template) Directives() []Directive {
	ds := []Directive{}
	t.traverseNodes(t.RootNode, func(tn TemplateNode) {
		t, ok := tn.(Directive)
		if ok {
			ds = append(ds, t)
		}
	})
	return ds
}

func (t Template) traverseNodes(n TemplateNode, f func(tn TemplateNode)) {
	for _, tn := range n.children() {
		f(tn)
		if len(tn.children()) > 0 {
			t.traverseNodes(tn, f)
		}
	}
}

type TemplateNode interface {
	node()
	children() []TemplateNode
	name() string
}

type Include struct {
	Name       string
	Parameters []Parameter
	Children   []TemplateNode
}

type Directive struct {
	Name       string
	Parameters []Parameter
	Children   []TemplateNode
}

type Parameter struct {
	Index int
	Value string
}

type Raw struct {
	Value string
}

type Root struct {
	Children []TemplateNode
}

type TemplateResult struct {
	Template Template
	Contents string
}

func (n Include) node()   {}
func (n Directive) node() {}
func (n Raw) node()       {}
func (n Root) node()      {}

func (n Include) children() []TemplateNode   { return n.Children }
func (n Directive) children() []TemplateNode { return n.Children }
func (n Raw) children() []TemplateNode       { return nil }
func (n Root) children() []TemplateNode      { return n.Children }

func (n Include) name() string   { return n.Name }
func (n Directive) name() string { return n.Name }
func (n Raw) name() string       { return n.Value }
func (n Root) name() string      { return "_root_" }
