package wc

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"sync"
	"unicode/utf8"
)

var defaultTemplate = "\t{{.Lines}}\t{{.Words}}\t{{.Bytes}}\n"

type WordCount struct {
	Bytes         int64
	Runes         int64
	Lines         int64
	Words         int64
	MaxLineLength int64
	input         *bufio.Reader
	output        io.Writer
	formatter     Formatter
	once          sync.Once
}

func NewWordCount(input io.Reader, output io.Writer, formatter Formatter) *WordCount {
	reader := bufio.NewReader(input)

	wc := &WordCount{
		input:     reader,
		output:    output,
		formatter: NewTemplateFormatter(defaultTemplate),
	}

	if formatter != nil {
		wc.formatter = formatter
	}

	return wc
}

func (wc *WordCount) Count() error {
	var err error

	hasLines := true
	for hasLines {
		line, err := wc.input.ReadBytes('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				err = fmt.Errorf("read string: %s", err)
			}
			hasLines = false
		}

		wc.Bytes += int64(len(line))
		wc.Runes += int64(utf8.RuneCount(line))
		wc.Words += int64(len(bytes.Fields(line)))

		if hasLines {
			wc.Lines++

			lineLen := int64(len(line))
			if lineLen > wc.MaxLineLength {
				wc.MaxLineLength = lineLen - 1
			}
		}
	}

	return err
}

func (wc *WordCount) Report() (int, error) {
	rp := report{
		Lines:         wc.Lines,
		Words:         wc.Words,
		Bytes:         wc.Bytes,
		Runes:         wc.Runes,
		MaxLineLength: wc.MaxLineLength,
	}

	rpBytes, err := wc.formatter.Format(rp)
	if err != nil {
		return 0, err
	}

	return wc.output.Write(rpBytes)
}
