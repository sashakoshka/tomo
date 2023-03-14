package patterns

import "image"
import "git.tebibyte.media/sashakoshka/tomo/canvas"

// Texture is a pattern that tiles the content of a canvas both horizontally and
// vertically.
type Texture struct {
	canvas.Canvas
}

// Draw tiles the pattern's canvas within the given bounds. The minimum
// points of the pattern's canvas and the destination canvas will be lined up.
func (pattern Texture) Draw (destination canvas.Canvas, bounds image.Rectangle) {
	drawBounds := bounds.Canon().Intersect(destination.Bounds())
	if drawBounds.Empty() { return }
	
	dstData, dstStride := destination.Buffer()
	srcData, srcStride := pattern.Buffer()
	srcBounds := pattern.Bounds()

	dstPoint := image.Point { }
	srcPoint := drawBounds.Min.Sub(bounds.Min).Add(srcBounds.Min)
	srcPoint.X = wrap(srcPoint.X, srcBounds.Min.X, srcBounds.Max.X)
	srcPoint.Y = wrap(srcPoint.Y, srcBounds.Min.Y, srcBounds.Max.Y)
	
	for dstPoint.Y = drawBounds.Min.Y; dstPoint.Y < drawBounds.Max.Y; dstPoint.Y ++ {
		srcPoint.X = srcBounds.Min.X
		dstPoint.X = drawBounds.Min.X
		dstYComponent := dstPoint.Y * dstStride
		srcYComponent := srcPoint.Y * srcStride
		
		for {
			dstIndex := dstYComponent + dstPoint.X
			srcIndex := srcYComponent + srcPoint.X
			dstData[dstIndex] = srcData[srcIndex]

			srcPoint.X ++
			if srcPoint.X >= srcBounds.Max.X {
				srcPoint.X = srcBounds.Min.X
			}

			dstPoint.X ++
			if dstPoint.X >= drawBounds.Max.X {
				break
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
