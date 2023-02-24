package shapes

import "image"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/shatter"

// FillRectangle draws a rectangular subset of one canvas onto the other. The
// offset point defines where the origin point of the source canvas is
// positioned in relation to the origin point of the destination canvas. To
// prevent the entire source canvas from being drawn, it must be cut with
// canvas.Cut().
func FillRectangle (
	destination canvas.Canvas,
	source      canvas.Canvas,
	offset      image.Point,
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

// StrokeRectangle is similar to FillRectangle, but it draws an inset outline of
// the source canvas onto the destination canvas. To prevent the entire source
// canvas's bounds from being used, it must be cut with canvas.Cut().
func StrokeRectangle (
	destination canvas.Canvas,
	source      canvas.Canvas,
	offset      image.Point,
	weight      int,
) {
	bounds := source.Bounds()
	insetBounds := bounds.Inset(weight)
	if insetBounds.Empty() {
		FillRectangle(destination, source, offset)
		return
	}

	top :=  image.Rect (
		bounds.Min.X, bounds.Min.Y,
		bounds.Max.X, insetBounds.Min.Y)
	bottom := image.Rect (
		bounds.Min.X, insetBounds.Max.Y,
		bounds.Max.X, bounds.Max.Y)
	left := image.Rect (
		bounds.Min.X, insetBounds.Min.Y,
		insetBounds.Min.X, insetBounds.Max.Y)
	right := image.Rect (
		insetBounds.Max.X, insetBounds.Min.Y,
		bounds.Max.X, insetBounds.Max.Y)
	
	FillRectangle (destination, canvas.Cut(source, top),    offset)
	FillRectangle (destination, canvas.Cut(source, bottom), offset)
	FillRectangle (destination, canvas.Cut(source, left),   offset)
	FillRectangle (destination, canvas.Cut(source, right),  offset)
}

// FillRectangleShatter is like FillRectangle, but it does not draw in areas
// specified in "rocks".
func FillRectangleShatter (
	destination canvas.Canvas,
	source      canvas.Canvas,
	offset      image.Point,
	rocks       []image.Rectangle,
) {
	tiles := shatter.Shatter(source.Bounds())
	for _, tile := range tiles {
		tile
	}
}
