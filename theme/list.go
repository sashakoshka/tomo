package theme

import "git.tebibyte.media/sashakoshka/tomo/artist"

var listPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Beveled {
			Highlight: artist.NewUniform(hex(0x383C3AFF)),
			Shadow:    artist.NewUniform(hex(0x999C99FF)),
		},
	},
	artist.Stroke { Pattern: artist.NewUniform(hex(0x999C99FF)) })


var listEntryPattern = artist.NewUniform(hex(0x999C99FF))

var selectedListEntryPattern = accentPattern

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
