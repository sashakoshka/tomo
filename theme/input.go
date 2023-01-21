package theme

import "git.tebibyte.media/sashakoshka/tomo/artist"

var inputPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Beveled {
			Highlight: artist.NewUniform(hex(0x89925AFF)),
			Shadow:    artist.NewUniform(hex(0xD2CB9AFF)),
		},
	},
	artist.Stroke { Pattern: artist.NewUniform(hex(0xD2CB9AFF)) })
var selectedInputPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke { Weight: 1, Pattern: accentPattern },
	artist.Stroke { Pattern: artist.NewUniform(hex(0xD2CB9AFF)) })
var disabledInputPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: weakForegroundPattern },
	artist.Stroke { Pattern: backgroundPattern })

func InputPattern (enabled, selected bool) (artist.Pattern) {
	if enabled {
		if selected {
			return selectedInputPattern
		} else {
			return inputPattern
		}
	} else {
		return disabledInputPattern
	}
}
