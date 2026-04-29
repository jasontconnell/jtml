package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/jasontconnell/jtml/process"
)

func main() {
	src := flag.String("src", "", "the source directory")
	dest := flag.String("dest", "", "the destination directory")
	flag.Parse()

	if *src == "" || *dest == "" {
		flag.PrintDefaults()
		return
	}

	templates, err := process.ParseTemplates(*src)
	if err != nil {
		log.Fatal(err)
	}
	templateResults, err := process.ProcessTemplates(templates)
	if err != nil {
		log.Fatal(err)
	}

	for _, res := range templateResults {
		filename := res.Template.Name + ".html"
		path := filepath.Join(*dest, filename)
		err = os.WriteFile(path, []byte(res.Contents), os.ModePerm)
		if err != nil {
			log.Println("writing file", filename, err)
		}
	}
}
