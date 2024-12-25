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
			log.Fatalf("failed to initalize TemplateFormatter: %s", err.Error())
		}
	}

	wordCount := wc.WordCount{}
	if len(flag.Args()) > 0 {
		for _, fileName := range flag.Args() {
			if err := validateFile(fileName); err != nil {
				log.Print(err)
				continue
			}

			f, err := os.Open(fileName)
			if err != nil {
				log.Printf("unable to open %s: %s\n", fileName, err.Error())
				continue
			}
			defer f.Close()

			wordCount.AddFile(f, wc.File{Name: f.Name()})
		}
	} else {
		wordCount.AddFile(os.Stdin, wc.File{Name: ""})
	}

	if wordCount.Files != nil {
		fmtr.Write(wordCount)
	}
}

func validateFile(fileName string) error {
	fileInfo, err := os.Lstat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s: file does not exist", fileName)
		} else {
			return fmt.Errorf("lstat %s: %w", fileName, err)
		}
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("%s is a directory", fileName)
	}

	return nil
}
