package wc

import (
	"strings"
	"testing"
)

func TestAddFile(t *testing.T) {
	tcases := []struct {
		content         []string
		outputWordCount WordCount
	}{
		{
			content: []string{"this is\n a test\n"},
			outputWordCount: WordCount{
				Files: []File{{
					Name:          "test",
					Lines:         2,
					Words:         4,
					Bytes:         16,
					Runes:         16,
					MaxLineLength: 7,
				},
				},
				TotalLines:    2,
				TotalWords:    4,
				TotalBytes:    16,
				TotalRunes:    16,
				MaxLineLength: 7,
			},
		},
	}

	for _, tc := range tcases {
		for _, file := range tc.content {
			sr := strings.NewReader(file)

			wordCount := WordCount{}
			wordCount.AddFile(sr, File{Name: "test"})

			if !wordCount.Equal(tc.outputWordCount) {
				t.Errorf("got %+v, expected %+v", wordCount, tc.outputWordCount)
			}
		}
	}
}
