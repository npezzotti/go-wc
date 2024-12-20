package main

import (
	"bufio"
	"flag"
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
	flag.Parse()

	tmplString := defaultTemplate
	if *goTemplate != "" {
		tmplString = *goTemplate
	}

	var fmtr wc.WCFormatter
	var err error
	if *jsonOutput {
		fmtr = wc.NewJsonFormatter(os.Stdout)
	} else {
		fmtr, err = wc.NewTemplateFormatter(
			tmplString,
			os.Stdout,
		)
		if err != nil {
			log.Fatalf("failed to initalize template formatter: %s", err.Error())
		}
	}

	wordCount := wc.WordCount{}
	if len(flag.Args()) > 0 {
		for _, fileName := range flag.Args() {
			fileInfo, err := os.Lstat(fileName)
			if err != nil {
				if os.IsNotExist(err) {
					log.Printf("%s: %s: file does not exist\n", progName, fileName)
				} else {
					log.Printf("%s: lstat %s: %s\n", progName, fileName, err.Error())
				}
				continue
			}
			
			if fileInfo.IsDir() {
				log.Printf("%s is a directory skipping", fileInfo.Name())
				continue
			}

			f, err := os.Open(fileName)
			if err != nil {
				log.Printf("%s: error opening file: %s\n", progName, err.Error())
				continue
			}
			defer f.Close()

			reader := bufio.NewReader(f)
			wordCount.AddFile(reader, wc.File{Name: f.Name()})
		}
	} else {
		reader := bufio.NewReader(os.Stdin)
		wordCount.AddFile(reader, wc.File{Name: ""})
	}

	if wordCount.Files != nil {
		fmtr.Write(wordCount)
	}
}
