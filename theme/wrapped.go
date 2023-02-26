package theme

import "image"
import "golang.org/x/image/font"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/canvas"

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

// Padding returns how much space should be between the bounds of a
// pattern whatever an element draws inside of it.
func (wrapped Wrapped) Padding (id Pattern) artist.Inset {
	real := wrapped.ensure()
	return real.Padding(id, wrapped.Case)
}

// Margin returns the left/right (x) and top/bottom (y) margins that
// should be put between any self-contained objects drawn within this
// pattern (if applicable).
func (wrapped Wrapped) Margin (id Pattern) image.Point {
	real := wrapped.ensure()
	return real.Margin(id, wrapped.Case)
}

// Sink returns a vector that should be added to an element's inner content when
// it is pressed down (if applicable) to simulate a 3D sinking effect.
func (wrapped Wrapped) Sink (id Pattern) image.Point {
	real := wrapped.ensure()
	return real.Sink(id, wrapped.Case)
}

// Hints returns rendering optimization hints for a particular pattern.
// These are optional, but following them may result in improved
// performance.
func (wrapped Wrapped) Hints (id Pattern) Hints {
	real := wrapped.ensure()
	return real.Hints(id, wrapped.Case)
}

func (wrapped Wrapped) ensure () (real Theme) {
	real = wrapped.Theme
	if real == nil { real = Default { } }
	return
}
