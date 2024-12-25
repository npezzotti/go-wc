package wc

import (
	"bytes"
	"testing"
)

func TestTemplateFormatter(t *testing.T) {
	tcases := []struct {
		tmplString   string
		wordCount    WordCount
		expectedTmpl string
	}{
		{
			tmplString: `{{range .Files}}{{.Lines}}\t{{.Words}}\t{{.Bytes}} {{.Name}}\n{{end}}`,
			wordCount: WordCount{
				Files: []File{
					{
						Name:  "test",
						Lines: 5,
						Words: 10,
						Bytes: 18,
					},
				},
			},
			expectedTmpl: "5\t10\t18 test\n",
		},
	}

	for _, tt := range tcases {
		byteSlice := make([]byte, 0)
		buf := bytes.NewBuffer(byteSlice)

		fmtr, err := NewTemplateFormatter(tt.tmplString, buf)
		if err != nil {
			t.Errorf("failed created template formatter: %s", err)
		}

		if err := fmtr.Write(tt.wordCount); err != nil {
			t.Errorf("unexpected error writing: %s", err)
		}

		if buf.String() != tt.expectedTmpl {
			t.Errorf("got %q, expected %q", buf.String(), tt.expectedTmpl)
		}
	}
}

func TestJsonFormatter(t *testing.T) {
	tcases := []struct {
		inputWordCount WordCount
		outputJson     string
	}{
		{
			inputWordCount: WordCount{
				Files: []File{
					{
						Name:          "test",
						Lines:         5,
						Words:         10,
						Bytes:         18,
						Runes:         18,
						MaxLineLength: 7,
					},
				},
				TotalLines:    5,
				TotalWords:    10,
				TotalBytes:    18,
				TotalRunes:    18,
				MaxLineLength: 7,
			},
			outputJson: `{"files":[{"name":"test","lines":5,"words":10,"bytes":18,"runes":18,"max-line-length":7}],` +
				`"total-lines":5,"total-words":10,"total-bytes":18,"total-runes":18,"max-line-length":7}`,
		},
	}

	for _, tt := range tcases {
		byteSlice := make([]byte, 0)
		buf := bytes.NewBuffer(byteSlice)

		fmtr := NewJsonFormatter(buf)
		if err := fmtr.Write(tt.inputWordCount); err != nil {
			t.Errorf("unexpected error writing json: %s", err)
		}

		if buf.String() != tt.outputJson {
			t.Errorf("got %s, expected %s", buf.String(), tt.outputJson)
		}
	}
}
