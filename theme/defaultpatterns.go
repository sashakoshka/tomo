package theme

import "image/color"
import "git.tebibyte.media/sashakoshka/tomo/artist"

var accentPattern         = artist.NewUniform(hex(0x408090FF))
var backgroundPattern     = artist.NewUniform(color.Gray16 { 0xAAAA })
var foregroundPattern     = artist.NewUniform(color.Gray16 { 0x0000 })
var weakForegroundPattern = artist.NewUniform(color.Gray16 { 0x4444 })
var strokePattern         = artist.NewUniform(color.Gray16 { 0x0000 })

var sunkenPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Beveled {
			artist.NewUniform(hex(0x3b534eFF)),
			artist.NewUniform(hex(0x97a09cFF)),
		},
	},
	artist.Stroke { Pattern: artist.NewUniform(hex(0x97a09cFF)) })
	
var focusedSunkenPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: accentPattern,
	},
	artist.Stroke { Pattern: artist.NewUniform(hex(0x97a09cFF)) })

var texturedSunkenPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Beveled {
			artist.NewUniform(hex(0x3b534eFF)),
			artist.NewUniform(hex(0x97a09cFF)),
		},
	},
	// artist.Stroke { Pattern: artist.Striped {
		// First: artist.Stroke {
			// Weight: 2,
			// Pattern: artist.NewUniform(hex(0x97a09cFF)),
		// },
		// Second: artist.Stroke {
			// Weight: 1,
			// Pattern: artist.NewUniform(hex(0x6e8079FF)),
		// },
	// }})
	
	artist.Stroke { Pattern: artist.Noisy {
		Low:  artist.NewUniform(hex(0x97a09cFF)),
		High: artist.NewUniform(hex(0x6e8079FF)),
	}})

var raisedPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Beveled {
			artist.NewUniform(hex(0xDBDBDBFF)),
			artist.NewUniform(hex(0x383C3AFF)),
		},
	},
	artist.Stroke { Pattern: artist.NewUniform(hex(0xAAAAAAFF)) })

var selectedRaisedPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Beveled {
			artist.NewUniform(hex(0xDBDBDBFF)),
			artist.NewUniform(hex(0x383C3AFF)),
		},
	},
	artist.Stroke { Weight: 1, Pattern: accentPattern },
	artist.Stroke { Pattern: artist.NewUniform(hex(0xAAAAAAFF)) })

var deadPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke { Pattern: artist.NewUniform(hex(0x97a09cFF)) })

var buttonPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Beveled {
			artist.NewUniform(hex(0xCCD5D2FF)),
			artist.NewUniform(hex(0x4B5B59FF)),
		},
	},
	artist.Stroke { Pattern: artist.NewUniform(hex(0x8D9894FF)) })
var selectedButtonPattern = artist.NewMultiBordered (
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
var pressedButtonPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Beveled {
			artist.NewUniform(hex(0x4B5B59FF)),
			artist.NewUniform(hex(0x8D9894FF)),
		},
	},
	artist.Stroke { Pattern: artist.NewUniform(hex(0x8D9894FF)) })
var pressedSelectedButtonPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Beveled {
			artist.NewUniform(hex(0x4B5B59FF)),
			artist.NewUniform(hex(0x8D9894FF)),
		},
	},
	artist.Stroke { Pattern: artist.NewUniform(hex(0x8D9894FF)) })
var disabledButtonPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: weakForegroundPattern },
	artist.Stroke { Pattern: backgroundPattern })

var darkButtonPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Beveled {
			artist.NewUniform(hex(0xaebdb9FF)),
			artist.NewUniform(hex(0x3b4947FF)),
		},
	},
	artist.Stroke { Pattern: artist.NewUniform(hex(0x6b7a75FF)) })
var pressedDarkButtonPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Beveled {
			artist.NewUniform(hex(0x3b4947FF)),
			artist.NewUniform(hex(0x6b7a75FF)),
		},
	},
	artist.Stroke { Pattern: artist.NewUniform(hex(0x6b7a75FF)) })

var inputPattern = artist.NewMultiBordered (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Beveled {
			artist.NewUniform(hex(0x89925AFF)),
			artist.NewUniform(hex(0xD2CB9AFF)),
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
