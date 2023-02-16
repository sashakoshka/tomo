package theme

import "image"
import "image/color"
import "golang.org/x/image/font"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/canvas"

// IconSize is a type representing valid icon sizes.
type IconSize int

const (
	IconSizeSmall IconSize = 16
	IconSizeLarge IconSize = 48
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

// Hints specifies rendering hints for a particular pattern. Elements can take
// these into account in order to gain extra performance.
type Hints struct {
	// StaticInset defines an inset rectangular area in the middle of the
	// pattern that does not change between PatternStates. If the inset is
	// zero on all sides, this hint does not apply.
	StaticInset Inset

	// Uniform specifies a singular color for the entire pattern. If the
	// alpha channel is zero, this hint does not apply.
	Uniform color.RGBA
}

// Theme represents a visual style configuration,
type Theme interface {
	// FontFace returns the proper font for a given style, size, and case.
	FontFace (FontStyle, FontSize, Case) font.Face

	// Icon returns an appropriate icon given an icon name, size, and case.
	Icon (string, IconSize, Case) canvas.Image

	// Pattern returns an appropriate pattern given a pattern name, case,
	// and state.
	Pattern (Pattern, PatternState, Case) artist.Pattern

	// Inset returns the area on all sides of a given pattern that is not
	// meant to be drawn on.
	Inset (Pattern, Case) Inset

	// Sink returns a vector that should be added to an element's inner
	// content when it is pressed down (if applicable) to simulate a 3D
	// sinking effect.
	Sink (Pattern, Case) image.Point

	// Hints returns rendering optimization hints for a particular pattern.
	// These are optional, but following them may result in improved
	// performance.
	Hints (Pattern, Case) Hints
}

// Wrapped wraps any theme and injects a case into it automatically so that it
// doesn't need to be specified for each query. Additionally, if the underlying
// theme is nil, it just uses the default theme instead.
type Wrapped struct {
	Theme
	Case
}

// FontFace returns the proper font for a given style and size.
func (wrapped Wrapped) FontFace (style FontStyle, size FontSize) font.Face {
	real := wrapped.ensure()
	return real.FontFace(style, size, wrapped.Case)
}

// Icon returns an appropriate icon given an icon name.
func (wrapped Wrapped) Icon (name string, size IconSize) canvas.Image {
	real := wrapped.ensure()
	return real.Icon(name, size, wrapped.Case)
}

// Pattern returns an appropriate pattern given a pattern name and state.
func (wrapped Wrapped) Pattern (id Pattern, state PatternState) artist.Pattern {
	real := wrapped.ensure()
	return real.Pattern(id, state, wrapped.Case)
}

// Inset returns the area on all sides of a given pattern that is not meant to
// be drawn on.
func (wrapped Wrapped) Inset (id Pattern) Inset {
	real := wrapped.ensure()
	return real.Inset(id, wrapped.Case)
}

// Sink returns a vector that should be added to an element's inner content when
// it is pressed down (if applicable) to simulate a 3D sinking effect.
func (wrapped Wrapped) Sink (id Pattern) image.Point {
	real := wrapped.ensure()
	return real.Sink(id, wrapped.Case)
}

func (wrapped Wrapped) ensure () (real Theme) {
	real = wrapped.Theme
	if real == nil { real = Default { } }
	return
}
