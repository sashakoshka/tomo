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

var sunkenPattern = artist.NewMultiBorder (
	artist.Border { Weight: 1, Stroke: strokePattern },
	artist.Border {
		Weight: 1,
		Stroke: artist.Chiseled {
			Highlight: artist.NewUniform(hex(0x373C3AFF)),
			Shadow:    artist.NewUniform(hex(0xDBDBDBFF)),
		},
	},
	artist.Border { Stroke: backgroundPattern })

func AccentPattern () (artist.Pattern) { return accentPattern }
func BackgroundPattern () (artist.Pattern) { return backgroundPattern }
func SunkenPattern () (artist.Pattern) { return sunkenPattern}
func ForegroundPattern (enabled bool) (artist.Pattern) {
	if enabled {
		return foregroundPattern
	} else {
		return weakForegroundPattern
	}
}
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
