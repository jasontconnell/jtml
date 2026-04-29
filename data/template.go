package data

type Template struct {
	Name  string
	Nodes []TemplateNode
}

func (t Template) IsPartial() bool {
	val := true
	for _, tn := range t.Nodes {
		d, ok := tn.(Directive)
		if ok && d.Name == "jtml" {
			val = false
			break
		}
	}
	return val
}

func (t Template) Directives() []Directive {
	ds := []Directive{}
	for _, n := range t.Nodes {
		t, ok := n.(Directive)
		if ok {
			ds = append(ds, t)
		}
	}
	return ds
}

type TemplateNode interface {
	node()
	children() []TemplateNode
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

type TemplateResult struct {
	Template Template
	Contents string
}

func (n Include) node()   {}
func (n Directive) node() {}
func (n Raw) node()       {}

func (n Include) children() []TemplateNode   { return n.Children }
func (n Directive) children() []TemplateNode { return n.Children }
func (n Raw) children() []TemplateNode       { return nil }
