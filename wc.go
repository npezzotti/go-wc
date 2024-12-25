package wc

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"unicode/utf8"
)

type WordCount struct {
	Files         []File `json:"files"`
	TotalLines    int64  `json:"total-lines"`
	TotalWords    int64  `json:"total-words"`
	TotalBytes    int64  `json:"total-bytes"`
	TotalRunes    int64  `json:"total-runes"`
	MaxLineLength int64  `json:"max-line-length"`
}

type File struct {
	Name          string `json:"name"`
	Lines         int64  `json:"lines"`
	Words         int64  `json:"words"`
	Bytes         int64  `json:"bytes"`
	Runes         int64  `json:"runes"`
	MaxLineLength int64  `json:"max-line-length"`
}

func (w *WordCount) AddFile(input io.Reader, f File) {
	reader := bufio.NewReader(input)

	hasLines := true
	for hasLines {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Printf("read string: %s", err)
			}
			hasLines = false
		}

		f.Bytes += int64(len(line))
		f.Runes += int64(utf8.RuneCount(line))
		f.Words += int64(len(bytes.Fields(line)))

		if hasLines {
			f.Lines++

			lineLen := int64(len(line))
			if lineLen > f.MaxLineLength {
				f.MaxLineLength = lineLen - 1
			}
		}
	}

	w.Files = append(w.Files, f)
	w.TotalLines += f.Lines
	w.TotalWords += f.Words
	w.TotalBytes += f.Bytes
	w.TotalRunes += f.Runes

	if w.MaxLineLength < f.MaxLineLength {
		w.MaxLineLength = f.MaxLineLength
	}
}

func (wc WordCount) Equal(wc2 WordCount) bool {
	if wc.TotalLines != wc2.TotalLines ||
		wc.TotalWords != wc2.TotalWords ||
		wc.TotalBytes != wc2.TotalBytes ||
		wc.TotalRunes != wc2.TotalRunes ||
		wc.MaxLineLength != wc2.MaxLineLength {
		return false
	}

	if wc.Files == nil && wc2.Files == nil {
		return true
	}

	if wc.Files == nil || wc2.Files == nil {
		return false
	}

	if len(wc.Files) != len(wc2.Files) {
		return false
	}

	for i := 0; i < len(wc.Files); i++ {
		if wc.Files[i] != wc2.Files[i] {
			return false
		}
	}

	return true
}
