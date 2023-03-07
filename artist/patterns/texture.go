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

	dstPoint := image.Point { }
	srcPoint := bounds.Min.Sub(realBounds.Min).Add(srcBounds.Min)
	srcPoint.X = wrap(srcPoint.X, srcBounds.Min.X, srcBounds.Max.X)
	srcPoint.Y = wrap(srcPoint.Y, srcBounds.Min.Y, srcBounds.Max.Y)
	srcPointXStart := srcPoint.X
	
	for dstPoint.Y = bounds.Min.Y; dstPoint.Y < bounds.Max.Y; dstPoint.Y ++ {
		srcPoint.X = srcPointXStart
		
		for dstPoint.X = bounds.Min.X; dstPoint.X < bounds.Max.X; dstPoint.X ++ {
			dstIndex := dstPoint.X + dstPoint.Y * dstStride
			srcIndex :=
				srcPoint.X +
				srcPoint.Y * srcStride
			dstData[dstIndex] = srcData[srcIndex]

			srcPoint.X ++
			if srcPoint.X >= srcBounds.Max.X {
				srcPoint.X = srcBounds.Min.X
			}
		}

		srcPoint.Y ++
		if srcPoint.Y >= srcBounds.Max.Y {
			srcPoint.Y = srcBounds.Min.Y
		}
	}
}

func wrap (value, min, max int) int {
	difference := max - min
	value = (value - min) % difference
	if value < 0 { value += difference }
	return value + min
}
