package process

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/jasontconnell/jtml/data"
	"github.com/jasontconnell/jtml/lexer"
	"github.com/jasontconnell/jtml/parser"
)

const (
	whitespace string = " \t"
	trimset    string = " \r\n\t"
	crlf       string = "\r\n"
	newline    string = "\n"
)

func ParseTemplates(path string, srcext string) ([]data.Template, error) {
	templates := []data.Template{}
	err := filepath.Walk(path, func(fpath string, f fs.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}

		dir, fn := filepath.Split(fpath)
		ext := filepath.Ext(fn)

		if ext != "."+srcext {
			return nil
		}

		b, err := os.ReadFile(fpath)
		if err != nil {
			return fmt.Errorf("reading file %s. %w", fpath, err)
		}

		tokens := lexer.Lex(string(b))
		p := parser.New()
		root := p.Parse(tokens)

		d := strings.TrimPrefix(dir, path)
		cleandir := strings.TrimPrefix(strings.TrimPrefix(d, "/"), "\\")

		rawname := strings.TrimSuffix(fn, "."+srcext)

		name := strings.Replace(cleandir, "\\", "/", -1) + strings.TrimLeft(rawname, "_")

		isPartial := strings.HasPrefix(fn, "_")

		t := toTemplate(rootNode{Node: root, Name: name, WriteFile: !isPartial})

		templates = append(templates, t)

		return nil
	})

	return templates, err
}

func ProcessTemplates(templates []data.Template) []data.TemplateResult {
	tm := make(map[string]data.Template)
	for _, t := range templates {
		tm[t.Name] = t
	}

	var results []data.TemplateResult
	for _, t := range templates {
		if !t.WriteFile {
			continue
		}

		output := processTemplate(t, tm, nil, 0)

		res := data.TemplateResult{
			Template: t,
			Contents: output,
		}
		results = append(results, res)
	}

	return results
}

func processTemplate(template data.Template, tm map[string]data.Template, parameters []data.Parameter, depth int) string {
	b := bytes.NewBufferString("")
	processNode(template, template.RootNode, tm, parameters, depth, b)

	return b.String()
}

func processNode(template data.Template, tn data.TemplateNode, tm map[string]data.Template, parameters []data.Parameter, depth int, buf *bytes.Buffer) {
	switch nt := tn.(type) {
	case data.Raw:
		val := replaceParams(nt.Value, parameters)
		buf.WriteString(adjustDepth(val, depth))
	case data.Stream:
		val := processStream(nt.Stream, parameters)
		buf.WriteString(adjustDepth(val, depth))
	case data.Include:
		tmp, ok := tm[nt.Name]
		pre, post := getPrePost(tmp)

		pre = replaceParams(adjustDepth(pre, depth), nt.Parameters)
		post = replaceParams(adjustDepth(post, depth), nt.Parameters)

		if ok {
			if pre != "" {
				buf.WriteString(pre)
			}

			val := processTemplate(tmp, tm, nt.Parameters, depth)
			buf.WriteString(val)
			processNodes(template, nt.Children, tm, parameters, depth+1, buf)

			if post != "" {
				buf.WriteString(post)
			}
		} else {
			buf.WriteString(adjustDepth("<!-- template "+nt.Name+" doesn't exist -->", depth))
		}
	case data.Root:
		processNodes(template, nt.Children, tm, parameters, depth, buf)
	}
}

func processNodes(template data.Template, nodes []data.TemplateNode, tm map[string]data.Template, parameters []data.Parameter, depth int, buf *bytes.Buffer) {
	for _, c := range nodes {
		processNode(template, c, tm, parameters, depth, buf)
	}
}

func processStream(nodes []data.TemplateNode, params []data.Parameter) string {
	var s string
	linestart := true
	for _, n := range nodes {
		if st, ok := n.(data.Raw); ok {
			var pre string
			if linestart && st.Depth > 0 {
				pre = strings.Repeat(" ", st.Depth)
			}
			s += pre + st.Value
			if st.Endline {
				linestart = true
				s += "\r\n"
			} else {
				s += " "
				linestart = false
			}
		}
	}
	s = replaceParams(s, params)
	return s
}

func paramValue(plist []data.TemplateNode) string {
	s := ""
	for _, p := range plist {
		if nt, ok := p.(data.Raw); ok {
			s += strings.Trim(nt.Value, trimset)
			if nt.Endline {
				s += "\n"
			} else {
				s += " "
			}
		}
	}
	return strings.TrimRight(s, whitespace)
}

func replaceParams(val string, plist []data.Parameter) string {
	for j := len(plist) - 1; j >= 0; j-- { // go backwards in case there's $10 and $1
		idx := j + 1
		val = strings.ReplaceAll(val, fmt.Sprintf("$%d", idx), plist[j].Value)
	}
	return val
}

func adjustDepth(val string, depth int) string {
	if len(val) == 0 {
		return ""
	}
	s := strings.Split(val, newline)
	for i := 0; i < len(s); i++ {
		s[i] = strings.TrimRight(s[i], trimset)
		if len(s[i]) > 0 {
			s[i] = strings.Repeat(" ", depth) + s[i] + crlf
		}
	}
	return strings.Join(s, "")
}

func getPrePost(tmp data.Template) (string, string) {
	var pre, post string
	for _, d := range tmp.Directives() {
		switch d.Name {
		case "open":
			pre = paramValue(d.Children)
		case "close":
			post = paramValue(d.Children)
		}
	}
	return pre, post
}
