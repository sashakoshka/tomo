package patterns

import "image"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/artist"

// Border is a pattern that behaves similarly to border-image in CSS. It divides
// a source canvas into nine sections...
//
//                         Inset[1]
//                         ┌──┴──┐
//           ┌─┌─────┬─────┬─────┐
//  Inset[0]─┤ │  0  │  1  │  2  │
//           └─├─────┼─────┼─────┤
//             │  3  │  4  │  5  │
//             ├─────┼─────┼─────┤─┐
//             │  6  │  7  │  8  │ ├─Inset[2]
//             └─────┴─────┴─────┘─┘
//             └──┬──┘
//             Inset[3]
//
// ... Where the bounds of section 4 are defined as the application of the
// pattern's inset to the canvas's bounds. The bounds of the other eight
// sections are automatically sized around it.
//
// When drawn to a destination canvas, the bounds of sections 1, 3, 4, 5, and 7
// are expanded or contracted to fit the given drawing bounds. All sections are
// rendered as if they are Texture patterns, meaning these flexible sections
// will repeat to fill in any empty space.
//
// This pattern can be used to make a static image texture into something that
// responds well to being resized.
type Border struct {
	canvas.Canvas
	artist.Inset
}

// Draw draws the border pattern onto the destination canvas within the given
// bounds.
func (pattern Border) Draw (destination canvas.Canvas, bounds image.Rectangle) {
	drawBounds := bounds.Canon().Intersect(destination.Bounds())
	if drawBounds.Empty() { return }

	srcSections := nonasect(pattern.Bounds(), pattern.Inset)
	srcTextures := [9]Texture { }
	for index, section := range srcSections {
		srcTextures[index].Canvas = canvas.Cut(pattern, section)
	}
	
	dstSections := nonasect(bounds, pattern.Inset)
	for index, section := range dstSections {
		srcTextures[index].Draw(destination, section)
	}
}

func nonasect (bounds image.Rectangle, inset artist.Inset) [9]image.Rectangle {
	center := inset.Apply(bounds)
	return [9]image.Rectangle {
		// top
		image.Rectangle {
			bounds.Min,
			center.Min },
		image.Rect (
			center.Min.X, bounds.Min.Y,
			center.Max.X, center.Min.Y),
		image.Rect (
			center.Max.X, bounds.Min.Y,
			bounds.Max.X, center.Min.Y),
			
		// center
		image.Rect (
			bounds.Min.X, center.Min.Y,
			center.Min.X, center.Max.Y),
		center,
		image.Rect (
			center.Max.X, center.Min.Y,
			bounds.Max.X, center.Max.Y),
			
		// bottom
		image.Rect (
			bounds.Min.X, center.Max.Y,
			center.Min.X, bounds.Max.Y),
		image.Rect (
			center.Min.X, center.Max.Y,
			center.Max.X, bounds.Max.Y),
		image.Rect (
			center.Max.X, center.Max.Y,
			bounds.Max.X, bounds.Max.Y),
	}
}
