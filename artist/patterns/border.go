package patterns

import "image"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/artist"

type Border struct {
	canvas.Canvas
	artist.Inset
}

func (pattern Border) Draw (destination canvas.Canvas, clip image.Rectangle) {
	bounds := clip.Canon().Intersect(destination.Bounds())
	if bounds.Empty() { return }

	srcSections := nonasect(pattern.Bounds(), pattern.Inset)
	srcTextures := [9]Texture { }
	for index, section := range srcSections {
		srcTextures[index] = Texture {
			Canvas: canvas.Cut(pattern, section),
		}
	}
	
	dstSections := nonasect(destination.Bounds(), pattern.Inset)
	for index, section := range dstSections {
		srcTextures[index].Draw(canvas.Cut(destination, section), clip)
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
