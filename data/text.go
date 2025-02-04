package data

import "fyne.io/fyne/v2/widget"

// NOTE: This package isn't the right spot for this, but I don't want to have a dedicated package for it, so here it be.

type textTag struct {
	adjust func(*widget.RichTextStyle)
}

var textTags = map[string]textTag{
	"b": {
		adjust: func(style *widget.RichTextStyle) {
			style.TextStyle.Bold = true
		},
	},
	"/b": {
		adjust: func(style *widget.RichTextStyle) {
			style.TextStyle.Bold = false
		},
	},
	"i": {
		adjust: func(style *widget.RichTextStyle) {
			style.TextStyle.Italic = true
		},
	},
	"/i": {
		adjust: func(style *widget.RichTextStyle) {
			style.TextStyle.Italic = false
		},
	},
	"ul": {
		adjust: func(style *widget.RichTextStyle) {
			style.TextStyle.Underline = true
		},
	},
	"/ul": {
		adjust: func(style *widget.RichTextStyle) {
			style.TextStyle.Underline = false
		},
	},
	"fixed": {
		adjust: func(style *widget.RichTextStyle) {
			style.TextStyle.Monospace = true
		},
	},
	"print": {
		adjust: func(style *widget.RichTextStyle) {
			style.TextStyle.Monospace = false
		},
	},
}

// TextToRichTextSegments converts a CF-friendly text string to a list of RichTextSegments.
func TextToRichTextSegments(text string) []widget.RichTextSegment {
	var segments []widget.RichTextSegment
	var style widget.RichTextStyle
	style.Inline = true
	pos := 0
	for i := 0; i < len(text); i++ {
		r := text[i]
		switch r {
		case '[':
			// Find end.
			for j := i + 1; j < len(text); j++ {
				if text[j] == ']' {
					tag := text[i+1 : j]
					// Submit as inline.
					text := text[pos:i]
					if text != "" {
						segments = append(segments, &widget.TextSegment{Text: text, Style: style})
					}
					pos = j + 1
					if t, ok := textTags[tag]; ok {
						t.adjust(&style)
					}
					break
				}
			}
		case '\n':
			if pos != i {
				style.Inline = false
				segments = append(segments, &widget.TextSegment{Text: text[pos:i], Style: style})
				style.Inline = true
				pos = i + 1
			}
		}
		if i == len(text)-1 && pos != i {
			style.Inline = false
			segments = append(segments, &widget.TextSegment{Text: text[pos:], Style: style})
		}
	}
	return segments
}
