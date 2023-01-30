package theme

import "git.tebibyte.media/sashakoshka/tomo/artist"

var listPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Beveled {
			artist.NewUniform(hex(0x383C3AFF)),
			artist.NewUniform(hex(0x999C99FF)),
		},
	},
	artist.Stroke { Pattern: artist.NewUniform(hex(0x999C99FF)) })

var selectedListPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke { Weight: 1, Pattern: accentPattern },
	artist.Stroke { Pattern: artist.NewUniform(hex(0x999C99FF)) })

// TODO: make these better, making use of the padded pattern. also, create
// selected variations for both of these.

var listEntryPattern = artist.NewMultiBordered (
	artist.Stroke { Pattern: artist.NewUniform(hex(0x999C99FF)) })

var onListEntryPattern = artist.NewMultiBordered (
	artist.Stroke { Pattern: artist.NewUniform(hex(0x6e8079FF)) })

var selectedListEntryPattern = artist.NewMultiBordered (
	artist.Stroke { Pattern: artist.NewUniform(hex(0x999C99FF)) })

var selectedOnListEntryPattern = artist.NewMultiBordered (
	artist.Stroke { Pattern: artist.NewUniform(hex(0x6e8079FF)) })
