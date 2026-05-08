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
	if n == nil {
		return
	}

	for _, tn := range n.children() {
		f(tn)
		if tn != nil && len(tn.children()) > 0 {
			t.traverseNodes(tn, f)
		}
	}
}

type TemplateNode interface {
	node()
	children() []TemplateNode
	name() string
	depth() int
}

type Include struct {
	Name       string
	Parameters []Parameter
	Children   []TemplateNode
	Depth      int
}

type Directive struct {
	Name       string
	Parameters []Parameter
	Children   []TemplateNode
	Depth      int
}

type Parameter struct {
	Index   int
	Value   string
	Endline bool
}

type Stream struct {
	Stream []TemplateNode
	Depth  int
}

type Raw struct {
	Value   string
	Depth   int
	Endline bool
}

type Root struct {
	Children []TemplateNode
	Depth    int
}

type TemplateResult struct {
	Template Template
	Contents string
}

func (n Include) node()   {}
func (n Directive) node() {}
func (n Raw) node()       {}
func (n Root) node()      {}
func (n Stream) node()    {}

func (n Include) children() []TemplateNode   { return n.Children }
func (n Directive) children() []TemplateNode { return n.Children }
func (n Raw) children() []TemplateNode       { return nil }
func (n Root) children() []TemplateNode      { return n.Children }
func (n Stream) children() []TemplateNode    { return n.Stream }

func (n Include) name() string   { return n.Name }
func (n Directive) name() string { return n.Name }
func (n Raw) name() string       { return n.Value }
func (n Root) name() string      { return "_root_" }
func (n Stream) name() string    { return "_stream_" }

func (n Include) depth() int   { return n.Depth }
func (n Directive) depth() int { return n.Depth }
func (n Raw) depth() int       { return n.Depth }
func (n Root) depth() int      { return -1 }
func (n Stream) depth() int    { return n.Depth }
