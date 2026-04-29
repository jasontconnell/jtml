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
		log.Println(fpath)
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
		log.Println("===============parsing template", fn)
		p := parser.New()
		root := p.Parse(tokens)

		roots = append(roots, rootNode{Node: root, Name: strings.TrimRight(fn, ext)})

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
		if t.IsPartial() {
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

	var pre, post string
	for _, d := range template.Directives() {
		if d.Name == "open" {
			pre = d.Parameters[0].Value
		} else if d.Name == "close" {
			post = d.Parameters[0].Value
		}
	}

	log.Println("pre and post", pre, post)

	for _, n := range template.Nodes {
		switch nt := n.(type) {
		case data.Raw:
			val := nt.Value
			if len(parameters) > 0 {
				for j := len(parameters) - 1; j >= 0; j-- { // go backwards in case there's $10 and $1
					idx := j + 1
					val = strings.ReplaceAll(val, fmt.Sprintf("$%d", idx), parameters[j].Value)
				}
			}
			b.WriteString(val + " ")
		case data.Include:
			tmp, ok := tm[nt.Name]

			if ok {
				val, err := processTemplate(tmp, tm, nt.Parameters)
				if err != nil {
					log.Println("error processing template", tmp.Name, "in", template.Name, err)
					continue
				}
				b.WriteString(pre + val + post + "\n")
			} else {
				log.Println("template not found", nt.Name)
			}
		}
	}
	return b.String(), nil
}
