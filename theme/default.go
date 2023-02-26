package theme

import "image"
import "golang.org/x/image/font"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/defaultfont"
import "git.tebibyte.media/sashakoshka/tomo/artist/patterns"

// Default is the default theme.
type Default struct { }

// FontFace returns the default font face.
func (Default) FontFace (style FontStyle, size FontSize, c Case) font.Face {
	switch style {
	case FontStyleBold:
		return defaultfont.FaceBold
	case FontStyleItalic:
		return defaultfont.FaceItalic
	case FontStyleBoldItalic:
		return defaultfont.FaceBoldItalic
	default:
		return defaultfont.FaceRegular
	}
}

// Icon returns an icon from the default set corresponding to the given name.
func (Default) Icon (string, IconSize, Case) canvas.Image {
	// TODO
	return nil
}

// Pattern returns a pattern from the default theme corresponding to the given
// pattern ID.
func (Default) Pattern (
	pattern Pattern,
	state PatternState,
	c Case,
) artist.Pattern {
	switch pattern {
	case PatternAccent:
	return patterns.Uhex(0xFF8800FF)
	case PatternBackground:
	return patterns.Uhex(0x000000FF)
	case PatternForeground:
	return patterns.Uhex(0xFFFFFFFF)
	// case PatternDead:
	// case PatternRaised:
	// case PatternSunken:
	// case PatternPinboard:
	// case PatternButton:
	// case PatternInput:
	// case PatternGutter:
	// case PatternHandle:
	default: return patterns.Uhex(0x888888FF)
	}
}

// Padding returns the default padding value for the given pattern.
func (Default) Padding (pattern Pattern, c Case) artist.Inset {
	return artist.Inset { 4, 4, 4, 4}
}

// Margin returns the default margin value for the given pattern.
func (Default) Margin (id Pattern, c Case) image.Point {
	return image.Pt(4, 4)
}

// Hints returns rendering optimization hints for a particular pattern.
// These are optional, but following them may result in improved
// performance.
func (Default) Hints (pattern Pattern, c Case) (hints Hints) {
	return
}

// Sink returns the default sink vector for the given pattern.
func (Default) Sink (pattern Pattern, c Case) image.Point {
	return image.Point { 1, 1 }
}
