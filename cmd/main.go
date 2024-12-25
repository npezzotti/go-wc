package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/npezzotti/wc"
)

const defaultTemplate = "{{range .Files}}{{printf \"%5d\" .Lines}}\t{{printf \"%5d\" .Words}}\t{{printf \"%5d\" .Bytes}} {{.Name}}\n{{end}}" +
	"{{if gt (len .Files) 1}}{{printf \"%5d\" .TotalLines}}\t{{printf \"%5d\" .TotalWords}}\t{{printf \"%5d\" .TotalBytes}} total\n{{end}}"

var (
	progName   = filepath.Base(flag.CommandLine.Name())
	jsonOutput = flag.Bool("json", false, "output in JSON")
	goTemplate = flag.String("template", "", "output in go-template")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix(fmt.Sprintf("%s: ", progName))

	flag.Parse()
	os.Exit(run(*jsonOutput, *goTemplate))
}

func run(jsonOutput bool, goTemplate string) (status int) {
	var fmtr wc.WCFormatter
	var err error

	if jsonOutput {
		fmtr = wc.NewJsonFormatter(os.Stdout)
	} else {
		tmplString := defaultTemplate
		if goTemplate != "" {
			tmplString = goTemplate
		}

		fmtr, err = wc.NewTemplateFormatter(
			tmplString,
			os.Stdout,
		)
		if err != nil {
			log.Printf("failed to initalize TemplateFormatter: %s", err.Error())
			status = 1
			return
		}
	}

	var wordCount wc.WordCount
	if len(flag.Args()) > 0 {
		for _, fileName := range flag.Args() {
			f, err := os.Open(fileName)
			if err != nil {
				log.Print(err)
				status = 1
				continue
			}
			defer f.Close()

			if err := wordCount.AddFile(f, wc.File{Name: f.Name()}); err != nil {
				log.Print(err)
				status = 1
			}
		}
	} else {
		wordCount.AddFile(os.Stdin, wc.File{Name: ""})
	}

	if wordCount.Files != nil {
		if err := fmtr.Write(wordCount); err != nil {
			log.Print(err)
			status = 1
		}
	}

	return
}
