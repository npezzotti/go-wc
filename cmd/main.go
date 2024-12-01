package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/npezzotti/wc"
)

var (
	json     = flag.Bool("json", false, "output in JSON")
	template = flag.String("template", "", "output in go-template")
)

func main() {
	flag.Parse()

	input, err := setupInput(flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer input.Close()

	var formatter wc.Formatter
	if *json {
		formatter = wc.JsonFormatter{}
	} else if *template != "" {
		formatter = wc.NewTemplateFormatter(sanitizeTemplate(*template))
	}

	wordCount := wc.NewWordCount(input, os.Stdout, formatter)

	if err := wordCount.Count(); err != nil {
		fmt.Fprintf(os.Stderr, "failed generating counts:%s\n", err.Error())
		os.Exit(1)
	}

	if _, err := wordCount.Report(); err != nil {
		fmt.Fprintf(os.Stderr, "failed creating report: %s\n", err.Error())
		os.Exit(1)
	}
}

func setupInput(args []string) (*os.File, error) {
	var input *os.File = os.Stdin
	if len(args) > 0 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			return nil, fmt.Errorf("unable to open file: %s\n", err.Error())
		}

		input = f
	}

	return input, nil
}

func sanitizeTemplate(tmpl string) string {
	tmpl = strings.TrimSpace(tmpl)
	r := strings.NewReplacer(`\t`, "\t", `\n`, "\n")
	tmpl = r.Replace(tmpl)

	return tmpl
}
