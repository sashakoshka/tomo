package theme

import "git.tebibyte.media/sashakoshka/tomo/artist"

var buttonPattern = artist.NewMultiBorder (
	artist.Border { Weight: 1, Stroke: strokePattern },
	artist.Border {
		Weight: 1,
		Stroke: artist.Chiseled {
			Highlight: artist.NewUniform(hex(0xCCD5D2FF)),
			Shadow:    artist.NewUniform(hex(0x4B5B59FF)),
		},
	},
	artist.Border { Stroke: artist.NewUniform(hex(0x8D9894FF)) })
var selectedButtonPattern = artist.NewMultiBorder (
	artist.Border { Weight: 1, Stroke: strokePattern },
	artist.Border {
		Weight: 1,
		Stroke: artist.Chiseled {
			Highlight: artist.NewUniform(hex(0xCCD5D2FF)),
			Shadow:    artist.NewUniform(hex(0x4B5B59FF)),
		},
	},
	artist.Border { Weight: 1, Stroke: accentPattern },
	artist.Border { Stroke: artist.NewUniform(hex(0x8D9894FF)) })
var pressedButtonPattern = artist.NewMultiBorder (
	artist.Border { Weight: 1, Stroke: strokePattern },
	artist.Border {
		Weight: 1,
		Stroke: artist.Chiseled {
			Highlight: artist.NewUniform(hex(0x4B5B59FF)),
			Shadow:    artist.NewUniform(hex(0x8D9894FF)),
		},
	},
	artist.Border { Stroke: artist.NewUniform(hex(0x8D9894FF)) })
var disabledButtonPattern = artist.NewMultiBorder (
	artist.Border { Weight: 1, Stroke: weakForegroundPattern },
	artist.Border { Stroke: backgroundPattern })

func ButtonPattern (enabled, selected, pressed bool) (artist.Pattern) {
	if enabled {
		if pressed {
			return pressedButtonPattern
		} else {
			if selected {
				return selectedButtonPattern
			} else {
				return buttonPattern
			}
		}
	} else {
		return disabledButtonPattern
	}
}
