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


var listEntryPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: artist.QuadBeveled {
		artist.NewUniform(hex(0x999C99FF)),
		strokePattern,
		artist.NewUniform(hex(0x999C99FF)),
		strokePattern,
	}},
	artist.Stroke { Pattern: artist.NewUniform(hex(0x999C99FF)) })

var selectedListEntryPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Beveled {
			artist.NewUniform(hex(0x3b534eFF)),
			artist.NewUniform(hex(0x97a09cFF)),
		},
	},
	artist.Stroke { Pattern: artist.NewUniform(hex(0x97a09cFF)) })

func ListPattern () (pattern artist.Pattern) {
	return listPattern
}

func ListEntryPattern (selected bool) (pattern artist.Pattern) {
	if selected {
		return selectedListEntryPattern
	} else {
		return listEntryPattern
	}
}
