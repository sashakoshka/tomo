package theme

import "image"
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

var accentPattern         = artist.NewUniform(hex(0x408090FF))
var backgroundPattern     = artist.NewUniform(color.Gray16 { 0xAAAA })
var foregroundPattern     = artist.NewUniform(color.Gray16 { 0x0000 })
var weakForegroundPattern = artist.NewUniform(color.Gray16 { 0x4444 })
var strokePattern         = artist.NewUniform(color.Gray16 { 0x0000 })

var sunkenPattern = artist.NewMultiBorder (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke {
		Weight: 1,
		Pattern: artist.Chiseled {
			Highlight: artist.NewUniform(hex(0x3b534eFF)),
			Shadow:    artist.NewUniform(hex(0x97a09cFF)),
		},
	},
	artist.Stroke { Pattern: artist.NewUniform(hex(0x97a09cFF)) })

var deadPattern = artist.NewMultiBorder (
	artist.Stroke { Weight: 1, Pattern: strokePattern },
	artist.Stroke { Pattern: artist.NewUniform(hex(0x97a09cFF)) })

func AccentPattern () (artist.Pattern) { return accentPattern }
func BackgroundPattern () (artist.Pattern) { return backgroundPattern }
func SunkenPattern () (artist.Pattern) { return sunkenPattern}
func DeadPattern () (artist.Pattern) { return deadPattern }
func ForegroundPattern (enabled bool) (artist.Pattern) {
	if enabled {
		return foregroundPattern
	} else {
		return weakForegroundPattern
	}
}

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

// SinkOffsetVector specifies a vector for things such as text to move by when a
// "sinking in" effect is desired, such as a button label during a button press.
func SinkOffsetVector () image.Point {
	return image.Point { 1, 1 }
}
