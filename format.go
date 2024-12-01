package wc

import (
	"bytes"
	"encoding/json"
	"text/template"
)

type Formatter interface {
	Format(rp report) ([]byte, error)
}

type report struct {
	Lines         int64 `json:"lines"`
	Words         int64 `json:"words"`
	Bytes         int64 `json:"bytes"`
	Runes         int64 `json:"runes"`
	MaxLineLength int64 `json:"max-line-length"`
}

type TemplateFormatter struct {
	tmpl string
}

func NewTemplateFormatter(tmpl string) TemplateFormatter {
	return TemplateFormatter{tmpl}
}

func (tf TemplateFormatter) Format(rp report) ([]byte, error) {
	var buf []byte
	writer := bytes.NewBuffer(buf)

	tmpl := template.New("")
	tt, err := tmpl.Parse(tf.tmpl)
	if err != nil {
		return nil, err
	}

	if err := tt.Execute(writer, rp); err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}

type JsonFormatter struct{}

func (jf JsonFormatter) Format(rp report) ([]byte, error) {
	jsonBytes, err := json.Marshal(rp)
	if err != nil {
		return nil, err
	}

	return jsonBytes, nil
}
