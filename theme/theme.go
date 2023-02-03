package theme

import "image"
import "golang.org/x/image/font"
import "git.tebibyte.media/sashakoshka/tomo/artist"

// FontStyle specifies stylistic alterations to a font face.
type FontStyle int; const (
	FontStyleRegular    FontStyle = 0
	FontStyleBold       FontStyle = 1
	FontStyleItalic     FontStyle = 2
	FontStyleBoldItalic FontStyle = 1 | 2
)

// FontSize specifies the general size of a font face in a semantic way.
type FontSize int; const (
	// FontSizeNormal is the default font size that should be used for most
	// things.
	FontSizeNormal FontSize = iota

	// FontSizeLarge is a larger font size suitable for things like section
	// headings.
	FontSizeLarge

	// FontSizeHuge is a very large font size suitable for things like
	// titles, wizard step names, digital clocks, etc.
	FontSizeHuge

	// FontSizeSmall is a smaller font size. Try not to use this unless it
	// makes a lot of sense to do so, because it can negatively impact
	// accessibility. It is useful for things like copyright notices at the
	// bottom of some window that the average user doesn't actually care
	// about.
	FontSizeSmall
)

// Pattern lists a number of cannonical pattern types, each with its own ID.
// This allows custom elements to follow themes, even those that do not
// explicitly support them.
type Pattern int; const (
	// PatternAccent is the accent color of the theme. It is safe to assume
	// that this is, by default, a solid color.
	PatternAccent Pattern = iota

	// PatternBackground is the background color of the theme. It is safe to
	// assume that this is, by default, a solid color.
	PatternBackground

	// PatternForeground is the foreground text color of the theme. It is
	// safe to assume that this is, by default, a solid color.
	PatternForeground

	// PatternDead is a pattern that is displayed on a "dead area" where no
	// controls exist, but there still must be some indication of visual
	// structure (such as in the corner between two scroll bars).
	PatternDead

	// PatternRaised is a generic raised pattern.
	PatternRaised

	// PatternSunken is a generic sunken pattern.
	PatternSunken

	// PatternPinboard is similar to PatternSunken, but it is textured.
	PatternPinboard

	// PatternButton is a button pattern.
	PatternButton

	// PatternInput is a pattern for input fields, editable text areas, etc.
	PatternInput

	// PatternGutter is a track for things to slide on.
	PatternGutter

	// PatternHandle is a handle that slides along a gutter.
	PatternHandle
)

// Theme represents a visual style configuration,
type Theme interface {
	// FontFace returns the proper font for a given style, size, and case.
	FontFace (FontStyle, FontSize, Case) font.Face

	// Icon returns an appropriate icon given an icon name and case.
	Icon (string, Case) artist.Pattern

	// Pattern returns an appropriate pattern given a pattern name, case,
	// and state.
	Pattern (Pattern, Case, PatternState) artist.Pattern

	// Inset returns the area on all sides of a given pattern that is not
	// meant to be drawn on.
	Inset (Pattern, Case) Inset

	// Sink returns a vector that should be added to an element's inner
	// content when it is pressed down (if applicable) to simulate a 3D
	// sinking effect.
	Sink (Pattern, Case) image.Point
}
