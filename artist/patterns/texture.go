package patterns

import "image"
import "git.tebibyte.media/sashakoshka/tomo/canvas"

// Texture is a pattern that tiles the content of a canvas both horizontally and
// vertically.
type Texture struct {
	canvas.Canvas
}

// Draw tiles the pattern's canvas within the clipping bounds. The minimum
// points of the pattern's canvas and the destination canvas will be lined up.
func (pattern Texture) Draw (destination canvas.Canvas, clip image.Rectangle) {
	realBounds := destination.Bounds()
	bounds := clip.Canon().Intersect(realBounds)
	if bounds.Empty() { return }
	
	dstData, dstStride := destination.Buffer()
	srcData, srcStride := pattern.Buffer()
	srcBounds := pattern.Bounds()

	point := image.Point { }
	for point.Y = bounds.Min.Y; point.Y < bounds.Max.Y; point.Y ++ {
	for point.X = bounds.Min.X; point.X < bounds.Max.X; point.X ++ {
		srcPoint := point.Sub(realBounds.Min).Add(srcBounds.Min)
		
		dstIndex := point.X + point.Y * dstStride
		srcIndex :=
			wrap(srcPoint.X, srcBounds.Min.X, srcBounds.Max.X) +
			wrap(srcPoint.Y, srcBounds.Min.Y, srcBounds.Max.Y) * srcStride
		dstData[dstIndex] = srcData[srcIndex]
	}}
}

func wrap (value, min, max int) int {
	difference := max - min
	value = (value - min) % difference
	if value < 0 { value += difference }
	return value + min
}
