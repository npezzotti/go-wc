package wc

import (
	"bytes"
	"strings"
	"testing"
)

func TestCount(t *testing.T) {
	tcases := []struct {
		name          string
		str           string
		Bytes         int64
		Runes         int64
		Lines         int64
		Words         int64
		MaxLineLength int64
	}{
		{
			name:          "standard string",
			str:           "this is a\ntest.\n",
			Bytes:         int64(16),
			Runes:         int64(16),
			Lines:         int64(2),
			Words:         int64(4),
			MaxLineLength: int64(9),
		},
		{
			name:          "string with no new lines",
			str:           "this is a test.",
			Bytes:         int64(15),
			Runes:         int64(15),
			Lines:         int64(0),
			Words:         int64(4),
			MaxLineLength: int64(0),
		},
		{
			name:          "string with unicode chars",
			str:           "This is a test 🙂\n",
			Bytes:         int64(20),
			Runes:         int64(17),
			Lines:         int64(1),
			Words:         int64(5),
			MaxLineLength: int64(19),
		},
	}

	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			b := make([]byte, 0)
			buf := bytes.NewBuffer(b)
			wc := NewWordCount(strings.NewReader(tc.str), buf, nil)

			if err := wc.Count(); err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			}

			if wc.Bytes != tc.Bytes {
				t.Errorf("incorrect byte count for %q: got %d, expected %d", tc.str, wc.Bytes, tc.Bytes)
			}

			if wc.Runes != tc.Runes {
				t.Errorf("incorrect rune count for %q: got %d, expected %d", tc.str, wc.Runes, tc.Runes)
			}

			if wc.Lines != tc.Lines {
				t.Errorf("incorrect line count for %q: got %d, expected %d", tc.str, wc.Lines, tc.Lines)
			}

			if wc.Words != tc.Words {
				t.Errorf("incorrect word count for %q: got %d, expected %d", tc.str, wc.Words, tc.Words)
			}

			if wc.MaxLineLength != tc.MaxLineLength {
				t.Errorf("incorrect max line length count for %q: got %d, expected %d", tc.str, wc.MaxLineLength, tc.MaxLineLength)
			}
		})
	}
}

func TestReport(t *testing.T) {
	want := "\t1\t4\t15\n"
	reader := strings.NewReader("this is a test\n")

	var buf bytes.Buffer
	wc := NewWordCount(reader, &buf, nil)

	if err := wc.Count(); err != nil {
		t.Error(err)
	}

	if _, err := wc.Report(); err != nil {
		t.Error(err)
	}

	if buf.String() != want {
		t.Errorf("got %s, expected %s", buf.String(), want)
	}
}
