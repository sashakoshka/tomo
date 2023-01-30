package theme

import "git.tebibyte.media/sashakoshka/tomo/artist"

var listPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Beveled {
			uhex(0x383C3AFF),
			uhex(0x999C99FF),
		},
	},
	artist.Stroke { Pattern: uhex(0x999C99FF) })

var focusedListPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke { Weight: 1, Pattern: accentPattern },
	artist.Stroke { Pattern: uhex(0x999C99FF) })

var listEntryPattern = artist.Padded {
	Stroke: uhex(0x383C3AFF),
	Fill:   uhex(0x999C99FF),
	Sides:  []int { 0, 0, 0, 1 },
}

var onListEntryPattern = artist.Padded {
	Stroke: uhex(0x383C3AFF),
	Fill:   uhex(0x6e8079FF),
	Sides:  []int { 0, 0, 0, 1 },
}

var focusedListEntryPattern = artist.Padded {
	Stroke: accentPattern,
	Fill:   uhex(0x999C99FF),
	Sides:  []int { 0, 1, 0, 1 },
}

var focusedOnListEntryPattern = artist.Padded {
	Stroke: accentPattern,
	Fill:   uhex(0x6e8079FF),
	Sides:  []int { 0, 1, 0, 1 },
}

