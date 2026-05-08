package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jasontconnell/jtml/process"
)

func main() {
	src := flag.String("src", "", "the source directory")
	dest := flag.String("dest", "", "the destination directory")
	srcext := flag.String("srcext", "jtml", "the source filename extension")
	destext := flag.String("destext", "html", "the destination filename extension")
	flag.Parse()

	start := time.Now()

	if *src == "" || *dest == "" {
		flag.PrintDefaults()
		return
	}

	templates, err := process.ParseTemplates(*src, *srcext)
	if err != nil {
		log.Fatal(err)
	}

	templateResults := process.ProcessTemplates(templates)

	var errs error
	for _, res := range templateResults {
		filename := res.Template.Name + "." + *destext
		path := filepath.Join(*dest, filename)
		err = os.WriteFile(path, []byte(res.Contents), os.ModePerm)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("writing file %s %w", filename, err))
		}
	}

	if errs != nil {
		log.Fatal(errs)
	}

	log.Println("finished", time.Since(start))
}
