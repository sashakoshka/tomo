package theme

import "git.tebibyte.media/sashakoshka/tomo/artist"

var buttonPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Chiseled {
			Highlight: artist.NewUniform(hex(0xCCD5D2FF)),
			Shadow:    artist.NewUniform(hex(0x4B5B59FF)),
		},
	},
	artist.Stroke { Pattern: artist.NewUniform(hex(0x8D9894FF)) })
var selectedButtonPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Chiseled {
			Highlight: artist.NewUniform(hex(0xCCD5D2FF)),
			Shadow:    artist.NewUniform(hex(0x4B5B59FF)),
		},
	},
	artist.Stroke { Weight: 1, Pattern: accentPattern },
	artist.Stroke { Pattern: artist.NewUniform(hex(0x8D9894FF)) })
var pressedButtonPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Chiseled {
			Highlight: artist.NewUniform(hex(0x4B5B59FF)),
			Shadow:    artist.NewUniform(hex(0x8D9894FF)),
		},
	},
	artist.Stroke { Pattern: artist.NewUniform(hex(0x8D9894FF)) })
var pressedSelectedButtonPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Chiseled {
			Highlight: artist.NewUniform(hex(0x4B5B59FF)),
			Shadow:    artist.NewUniform(hex(0x8D9894FF)),
		},
	},
	artist.Stroke { Pattern: artist.NewUniform(hex(0x8D9894FF)) })
var disabledButtonPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: weakForegroundPattern },
	artist.Stroke { Pattern: backgroundPattern })

func ButtonPattern (enabled, selected, pressed bool) (artist.Pattern) {
	if enabled {
		if pressed {
			if selected {
				return pressedSelectedButtonPattern
			} else {
				return pressedButtonPattern
			}
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
