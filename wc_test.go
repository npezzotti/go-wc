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

func TestEqual(t *testing.T) {
	tcases := []struct {
		wc1 WordCount
		wc2 WordCount
		res bool
	}{
		{
			wc1: WordCount{
				Files: []File{
					{
						Name:          "test",
						Lines:         15,
						Words:         15,
						Bytes:         14,
						Runes:         25,
						MaxLineLength: 10,
					},
				},
				TotalLines:    10,
				TotalWords:    5,
				TotalBytes:    19,
				TotalRunes:    20,
				MaxLineLength: 15,
			},
			wc2: WordCount{
				Files: []File{
					{
						Name:          "test",
						Lines:         15,
						Words:         15,
						Bytes:         14,
						Runes:         25,
						MaxLineLength: 10,
					},
				},
				TotalLines:    10,
				TotalWords:    5,
				TotalBytes:    19,
				TotalRunes:    20,
				MaxLineLength: 15,
			},
			res: true,
		},
		{
			wc1: WordCount{},
			wc2: WordCount{},
			res: true,
		},
		{
			wc1: WordCount{
				TotalLines:    10,
				TotalWords:    5,
				TotalBytes:    19,
				TotalRunes:    20,
				MaxLineLength: 15,
			},
			wc2: WordCount{
				TotalLines:    10,
				TotalWords:    5,
				TotalBytes:    19,
				TotalRunes:    20,
				MaxLineLength: 16,
			},
			res: false,
		},
		{
			wc1: WordCount{
				Files: []File{
					{
						Name:          "test",
						Lines:         15,
						Words:         15,
						Bytes:         14,
						Runes:         25,
						MaxLineLength: 10,
					},
				},
			},
			wc2: WordCount{},
			res: false,
		},
		{
			wc1: WordCount{
				Files: []File{
					{
						Name:          "test",
						Lines:         15,
						Words:         15,
						Bytes:         14,
						Runes:         25,
						MaxLineLength: 10,
					},
				},
				TotalLines:    10,
				TotalWords:    5,
				TotalBytes:    19,
				TotalRunes:    20,
				MaxLineLength: 15,
			},
			wc2: WordCount{
				Files: []File{
					{
						Name: "test",
					},
				},
				TotalLines:    10,
				TotalWords:    5,
				TotalBytes:    19,
				TotalRunes:    20,
				MaxLineLength: 15,
			},
			res: false,
		},
		{
			wc1: WordCount{
				Files: []File{
					{
						Name: "test",
					},
				},
			},
			wc2: WordCount{
				Files: []File{
					{
						Name: "test1",
					},
					{
						Name: "test2",
					},
				},
			},
			res: false,
		},
	}

	for _, tt := range tcases {
		if res := tt.wc1.Equal(tt.wc2); res != tt.res {
			t.Errorf("unexpected %t result for %+v and %+v", res, tt.wc1, tt.wc2)
		}
	}
}
