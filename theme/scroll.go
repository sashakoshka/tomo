package theme

import "git.tebibyte.media/sashakoshka/tomo/artist"

var scrollGutterPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Beveled {
			artist.NewUniform(hex(0x3b534eFF)),
			artist.NewUniform(hex(0x6e8079FF)),
		},
	},
	artist.Stroke { Pattern: artist.NewUniform(hex(0x6e8079FF)) })
var disabledScrollGutterPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: weakForegroundPattern },
	artist.Stroke { Pattern: backgroundPattern })
var scrollBarPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Beveled {
			artist.NewUniform(hex(0xCCD5D2FF)),
			artist.NewUniform(hex(0x4B5B59FF)),
		},
	},
	artist.Stroke { Pattern: artist.NewUniform(hex(0x8D9894FF)) })
var selectedScrollBarPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Beveled {
			artist.NewUniform(hex(0xCCD5D2FF)),
			artist.NewUniform(hex(0x4B5B59FF)),
		},
	},
	artist.Stroke { Weight: 1, Pattern: accentPattern },
	artist.Stroke { Pattern: artist.NewUniform(hex(0x8D9894FF)) })
var pressedScrollBarPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Beveled {
			artist.NewUniform(hex(0xCCD5D2FF)),
			artist.NewUniform(hex(0x4B5B59FF)),
		},
	},
	artist.Stroke { Weight: 1, Pattern: artist.NewUniform(hex(0x8D9894FF)) },
	artist.Stroke { Pattern: artist.NewUniform(hex(0x7f8c89FF)) })
var pressedSelectedScrollBarPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Beveled {
			artist.NewUniform(hex(0xCCD5D2FF)),
			artist.NewUniform(hex(0x4B5B59FF)),
		},
	},
	artist.Stroke { Weight: 1, Pattern: accentPattern },
	artist.Stroke { Pattern: artist.NewUniform(hex(0x7f8c89FF)) })
var disabledScrollBarPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: weakForegroundPattern },
	artist.Stroke { Pattern: backgroundPattern })
