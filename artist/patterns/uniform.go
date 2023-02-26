package patterns

import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/artist/shapes"

// Uniform is a pattern that draws a solid color.
type Uniform color.RGBA

// Draw fills the clipping rectangle with the pattern's color.
func (pattern Uniform) Draw (destination canvas.Canvas, clip image.Rectangle) {
	shapes.FillColorRectangle(destination, color.RGBA(pattern), clip)
}

// Uhex creates a new Uniform pattern from an RGBA integer value.
func Uhex (color uint32) (uniform Uniform) {
	return Uniform(artist.Hex(color))
}
