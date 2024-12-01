package wc

import (
	"testing"
)

func TestNewTemplateFormatter(t *testing.T) {
	tmpl := "{{.Bytes}}"
	tc := NewTemplateFormatter(tmpl)

	if tc.tmpl != tmpl {
		t.Errorf("\"tc.tmpl\" is %q, should be %q", tc.tmpl, tc.tmpl)
	}
}

func TestTemplateFormatterFormat(t *testing.T) {
	rp := report{
		Lines:         2,
		Words:         10,
		Bytes:         15,
		Runes:         15,
		MaxLineLength: 6,
	}

	tf := NewTemplateFormatter("{{.Bytes}}")
	format, err := tf.Format(rp)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if output := string(format); output != "15" {
		t.Errorf("expected %q, got %q", "15", output)
	}
}
