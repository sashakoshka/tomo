package theme

import "image/color"
import "golang.org/x/image/font"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/defaultfont"

// none of these colors are final! TODO: generate these values from a theme
// file at startup.

func hex (color uint32) (c color.RGBA) {
	c.A = uint8(color)
	c.B = uint8(color >>  8)
	c.G = uint8(color >> 16)
	c.R = uint8(color >> 24)
	return
}

func uhex (color uint32) (pattern artist.Pattern) {
	return artist.NewUniform(hex(color))
}

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

// TODO: load fonts from an actual source instead of using defaultfont

// FontFaceRegular returns the font face to be used for normal text.
func FontFaceRegular () font.Face {
	return defaultfont.FaceRegular
}

// FontFaceBold returns the font face to be used for bolded text.
func FontFaceBold () font.Face {
	return defaultfont.FaceBold
}

// FontFaceItalic returns the font face to be used for italicized text.
func FontFaceItalic () font.Face {
	return defaultfont.FaceItalic
}

// FontFaceBoldItalic returns the font face to be used for text that is both
// bolded and italicized.
func FontFaceBoldItalic () font.Face {
	return defaultfont.FaceBoldItalic
}

// Padding returns how spaced out things should be on the screen. Generally,
// text should be offset from its container on all sides by this amount.
func Padding () int {
	return 8
}

// HandleWidth returns how large grab handles should typically be. This is
// important for accessibility reasons.
func HandleWidth () int {
	return Padding() * 2
}
