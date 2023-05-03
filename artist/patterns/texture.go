package patterns

import "image"
import "git.tebibyte.media/sashakoshka/tomo/artist"

// Texture is a pattern that tiles the content of a canvas both horizontally and
// vertically.
type Texture struct {
	artist.Canvas
}

// Draw tiles the pattern's canvas within the given bounds. The minimum
// point of the pattern's canvas will be lined up with the minimum point of the
// bounding rectangle.
func (pattern Texture) Draw (destination artist.Canvas, bounds image.Rectangle) {
	dstBounds := bounds.Canon().Intersect(destination.Bounds())
	if dstBounds.Empty() { return }

	dstData, dstStride := destination.Buffer()
	srcData, srcStride := pattern.Buffer()
	srcBounds := pattern.Bounds()

	// offset is a vector that is added to points in destination space to
	// convert them to points in source space
	offset := srcBounds.Min.Sub(bounds.Min)

	// calculate the starting position in source space
	srcPoint := dstBounds.Min.Add(offset)
	srcPoint.X = wrap(srcPoint.X, srcBounds.Min.X, srcBounds.Max.X)
	srcPoint.Y = wrap(srcPoint.Y, srcBounds.Min.Y, srcBounds.Max.Y)
	srcStartPoint := srcPoint

	// for each row
	dstPoint := image.Point { }
	for dstPoint.Y = dstBounds.Min.Y; dstPoint.Y < dstBounds.Max.Y; dstPoint.Y ++ {
		srcPoint.X = srcStartPoint.X
		dstPoint.X = dstBounds.Min.X
		dstYComponent := dstPoint.Y * dstStride
		srcYComponent := srcPoint.Y * srcStride

		// for each pixel in the row
		for {
			// draw pixel
			dstIndex := dstYComponent + dstPoint.X
			srcIndex := srcYComponent + srcPoint.X
			dstData[dstIndex] = srcData[srcIndex]

			// increment X in source space. wrap to start if out of
			// bounds.
			srcPoint.X ++
			if srcPoint.X >= srcBounds.Max.X {
				srcPoint.X = srcBounds.Min.X
			}

			// increment X in destination space. stop drawing this
			// row if out of bounds.
			dstPoint.X ++
			if dstPoint.X >= dstBounds.Max.X {
				break
			}
		}

		// increment row in source space. wrap to start if out of
		// bounds.
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
