package theme

import "git.tebibyte.media/sashakoshka/tomo/artist"

var inputPattern = artist.NewMultiBorder (
	artist.Border { Weight: 1, Stroke: strokePattern },
	artist.Border {
		Weight: 1,
		Stroke: artist.Chiseled {
			Highlight: artist.NewUniform(hex(0x89925AFF)),
			Shadow:    artist.NewUniform(hex(0xD2CB9AFF)),
		},
	},
	artist.Border { Stroke: artist.NewUniform(hex(0xD2CB9AFF)) })
var selectedInputPattern = artist.NewMultiBorder (
	artist.Border { Weight: 1, Stroke: strokePattern },
	artist.Border { Weight: 1, Stroke: accentPattern },
	artist.Border { Stroke: artist.NewUniform(hex(0xD2CB9AFF)) })
var disabledInputPattern = artist.NewMultiBorder (
	artist.Border { Weight: 1, Stroke: weakForegroundPattern },
	artist.Border { Stroke: backgroundPattern })

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
