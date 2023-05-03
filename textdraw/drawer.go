package textdraw

import "image"
import "unicode"
import "image/draw"
import "image/color"
import "golang.org/x/image/math/fixed"
import "tomo/artist"

// Drawer is an extended TypeSetter that is able to draw text. Much like
// TypeSetter, It has no constructor and its zero value can be used safely.
type Drawer struct { TypeSetter }

// Draw draws the drawer's text onto the specified canvas at the given offset.
func (drawer Drawer) Draw (
	destination artist.Canvas,
	color       color.RGBA,
	offset      image.Point,
) (
	updatedRegion image.Rectangle,
) {
	drawer.For (func (
		index    int,
		char     rune,
		position fixed.Point26_6,
	) bool {
		destinationRectangle,
		mask, maskPoint, _, ok := drawer.face.Glyph (
			fixed.P (
				offset.X + position.X.Round(),
				offset.Y + position.Y.Round()),
			char)
		if !ok || unicode.IsSpace(char) || char == 0 {
			return true
		}

		// FIXME:? clip destination rectangle if we are on the cusp of
		// the maximum height.

		draw.DrawMask (
			destination,
			destinationRectangle,
			image.NewUniform(color), image.Point { },
			mask, maskPoint,
			draw.Over)

		updatedRegion = updatedRegion.Union(destinationRectangle)
		return true
	})
	return
}
