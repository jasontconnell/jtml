package process

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jasontconnell/jtml/data"
	"github.com/jasontconnell/jtml/lexer"
	"github.com/jasontconnell/jtml/parser"
)

func ParseTemplates(path string) ([]data.Template, error) {
	roots := []rootNode{}
	err := filepath.Walk(path, func(fpath string, f fs.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}
		_, fn := filepath.Split(fpath)
		ext := filepath.Ext(fn)

		if ext != ".txt" {
			return nil
		}

		b, err := os.ReadFile(fpath)
		if err != nil {
			return fmt.Errorf("reading file %s. %w", fpath, err)
		}

		tokens := lexer.Lex(string(b))
		p := parser.New()
		root := p.Parse(tokens)

		isPartial := strings.HasPrefix(fn, "_")

		roots = append(roots, rootNode{Node: root, Name: strings.TrimRight(strings.TrimLeft(fn, "_"), ext), IsPartial: isPartial})

		return nil
	})

	templates := []data.Template{}
	for _, r := range roots {
		t := toTemplate(r)
		templates = append(templates, t)
	}

	return templates, err
}

func ProcessTemplates(templates []data.Template) ([]data.TemplateResult, error) {
	tm := make(map[string]data.Template)
	for _, t := range templates {
		tm[t.Name] = t
	}

	var results []data.TemplateResult
	var errs error
	for _, t := range templates {
		if t.IsPartial {
			continue
		}

		output, err := processTemplate(t, tm, nil)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("processing template %s %w", t.Name, err))
			continue
		}

		res := data.TemplateResult{
			Template: t,
			Contents: output,
		}
		results = append(results, res)
	}

	return results, errs
}

func processTemplate(template data.Template, tm map[string]data.Template, parameters []data.Parameter) (string, error) {
	b := bytes.NewBufferString("")
	processNode(template, template.RootNode, tm, parameters, b)

	return b.String(), nil
}

func processNode(template data.Template, tn data.TemplateNode, tm map[string]data.Template, parameters []data.Parameter, buf *bytes.Buffer) {
	switch nt := tn.(type) {
	case data.Raw:
		val := nt.Value
		if len(parameters) > 0 {
			for j := len(parameters) - 1; j >= 0; j-- { // go backwards in case there's $10 and $1
				idx := j + 1
				val = strings.ReplaceAll(val, fmt.Sprintf("$%d", idx), parameters[j].Value)
			}
		}
		buf.WriteString(val)
	case data.Include:
		tmp, ok := tm[nt.Name]
		pre, post := getPrePost(tmp)

		if ok {
			buf.WriteString(pre)
			val, err := processTemplate(tmp, tm, nt.Parameters)
			if err != nil {
				log.Println("error processing template", tmp.Name, "in", template.Name, err)

			}
			buf.WriteString(val + "\n")
			processNodes(template, nt.Children, tm, nt.Parameters, buf)
			buf.WriteString(post + "\n")
		}
	case data.Root:
		processNodes(template, nt.Children, tm, parameters, buf)
	}
}

func processNodes(template data.Template, nodes []data.TemplateNode, tm map[string]data.Template, parameters []data.Parameter, buf *bytes.Buffer) {
	for _, c := range nodes {
		processNode(template, c, tm, parameters, buf)
	}
}

func paramValue(p []data.Parameter) string {
	s := ""
	for _, pp := range p {
		s += pp.Value + " "
	}
	return strings.TrimRight(s, " ")
}

func getPrePost(tmp data.Template) (string, string) {
	var pre, post string
	for _, d := range tmp.Directives() {
		if d.Name == "open" {
			pre = paramValue(d.Parameters)
		} else if d.Name == "close" {
			post = paramValue(d.Parameters)
		}
	}
	return pre, post
}
