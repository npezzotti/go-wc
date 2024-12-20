package wc

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/template"
)

type WCFormatter interface {
	Write(wordCount WordCount) error
}

type JsonFormatter struct {
	Output io.Writer
}

func NewJsonFormatter(writer io.Writer) *JsonFormatter {
	return &JsonFormatter{Output: writer}
}

func (jf JsonFormatter) Write(wordCount WordCount) error {
	jsonBytes, err := json.Marshal(wordCount)
	if err != nil {
		return fmt.Errorf("unable to marshal json: %s", err.Error())
	}

	_, err = fmt.Fprintf(jf.Output, string(jsonBytes))

	return err
}

type TemplateFormatter struct {
	tmpl   *template.Template
	output io.Writer
}

func NewTemplateFormatter(tmplString string, writer io.Writer) (*TemplateFormatter, error) {
	tmpl, err := template.New("").Parse(cleanTemplateString(tmplString))
	if err != nil {
		return nil, err
	}

	return &TemplateFormatter{
		tmpl:   tmpl,
		output: writer,
	}, nil
}

func (tf TemplateFormatter) Write(wordCount WordCount) error {
	if err := tf.tmpl.Execute(tf.output, wordCount); err != nil {
		return fmt.Errorf("failed writing template: %s", err.Error())
	}

	return nil
}

func cleanTemplateString(tmplString string) string {
	tmplString = strings.TrimSpace(tmplString)
	r := strings.NewReplacer(`\t`, "\t", `\n`, "\n")
	tmplString = r.Replace(tmplString)

	return tmplString
}
