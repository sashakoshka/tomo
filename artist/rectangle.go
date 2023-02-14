package artist

import "image"
import "git.tebibyte.media/sashakoshka/tomo/canvas"

// Paste transfers one canvas onto another, offset by the specified point.
func Paste (
	destination canvas.Canvas,
	source canvas.Canvas,
	offset image.Point,
) (
	updatedRegion image.Rectangle,
) {
	dstData, dstStride := destination.Buffer()
	srcData, srcStride := source.Buffer()

	sourceBounds :=
		source.Bounds().Canon().
		Intersect(destination.Bounds().Sub(offset))
	if sourceBounds.Empty() { return }
	
	updatedRegion = sourceBounds.Add(offset)
	for y := sourceBounds.Min.Y; y < sourceBounds.Max.Y; y ++ {
	for x := sourceBounds.Min.X; x < sourceBounds.Max.X; x ++ {
		dstData[x + offset.X + (y + offset.Y) * dstStride] =
			srcData[x + y * srcStride]
	}}

	return
}

// FillRectangle draws a filled rectangle with the specified pattern.
func FillRectangle (
	destination canvas.Canvas,
	source Pattern,
	bounds image.Rectangle,
) (
	updatedRegion image.Rectangle,
) {
	return FillRectangleClip(destination, source, bounds, bounds)
}

// FillRectangleClip is similar to FillRectangle, but it clips the pattern to
// a specified rectangle mask. That is—the pattern will be queried as if it
// were drawn without the mask, but only the area specified by the intersection
// of bounds and mask will be drawn to.
func FillRectangleClip (
	destination canvas.Canvas,
	source Pattern,
	bounds image.Rectangle,
	mask   image.Rectangle,
) (
	updatedRegion image.Rectangle,
) {
	data, stride := destination.Buffer()
	realBounds := bounds
	bounds =
		bounds.Canon().
		Intersect(mask.Canon()).
		Intersect(destination.Bounds())
	if bounds.Empty() { return }
	updatedRegion = bounds

	realWidth, realHeight := realBounds.Dx(), realBounds.Dy()
	patternOffset := realBounds.Min.Sub(bounds.Min)

	width, height := bounds.Dx(), bounds.Dy()
	for y := 0; y < height; y ++ {
	for x := 0; x < width;  x ++ {
		data[x + bounds.Min.X + (y + bounds.Min.Y) * stride] =
			source.AtWhen (
				x - patternOffset.X, y - patternOffset.Y,
				realWidth, realHeight)
	}}
	return
}

// StrokeRectangle draws the outline of a rectangle with the specified line
// weight and pattern.
func StrokeRectangle (
	destination canvas.Canvas,
	source Pattern,
	weight int,
	bounds image.Rectangle,
) {
	bounds = bounds.Canon()
	insetBounds := bounds.Inset(weight)
	if insetBounds.Empty() {
		FillRectangle(destination, source, bounds)
		return
	}

	// top
	FillRectangle (destination, source, image.Rect (
		bounds.Min.X, bounds.Min.Y,
		bounds.Max.X, insetBounds.Min.Y))
		
	// bottom
	FillRectangle (destination, source, image.Rect (
		bounds.Min.X, insetBounds.Max.Y,
		bounds.Max.X, bounds.Max.Y))

	// left
	FillRectangle (destination, source, image.Rect (
		bounds.Min.X, insetBounds.Min.Y,
		insetBounds.Min.X, insetBounds.Max.Y))
		
	// right
	FillRectangle (destination, source, image.Rect (
		insetBounds.Max.X, insetBounds.Min.Y,
		bounds.Max.X, insetBounds.Max.Y))
}
