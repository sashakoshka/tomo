package shapes

import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/shatter"

// TODO: return updatedRegion for all routines in this package

// FillRectangle draws the content of one canvas onto another. The offset point
// defines where the origin point of the source canvas is positioned in relation
// to the origin point of the destination canvas. To prevent the entire source
// canvas from being drawn, it must be cut with canvas.Cut().
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
	FillRectangleShatter(destination, source, offset, insetBounds)
}

// FillRectangleShatter is like FillRectangle, but it does not draw in areas
// specified in "rocks".
func FillRectangleShatter (
	destination canvas.Canvas,
	source      canvas.Canvas,
	offset      image.Point,
	rocks       ...image.Rectangle,
) {
	tiles := shatter.Shatter(source.Bounds().Sub(offset), rocks...)
	for _, tile := range tiles {
		FillRectangle(destination, canvas.Cut(source, tile), offset)
	}
}

// FillColorRectangle fills a rectangle within the destination canvas with a
// solid color.
func FillColorRectangle (
	destination canvas.Canvas,
	color       color.RGBA,
	bounds      image.Rectangle,
) (
	updatedRegion image.Rectangle,
) {
	dstData, dstStride := destination.Buffer()
	bounds = bounds.Canon().Intersect(destination.Bounds())
	if bounds.Empty() { return }
	
	updatedRegion = bounds
	for y := bounds.Min.Y; y < bounds.Max.Y; y ++ {
	for x := bounds.Min.X; x < bounds.Max.X; x ++ {
		dstData[x + y * dstStride] = color
	}}
	
	return
}

// FillColorRectangleShatter is like FillColorRectangle, but it does not draw in
// areas specified in "rocks".
func FillColorRectangleShatter (
	destination canvas.Canvas,
	color       color.RGBA,
	bounds      image.Rectangle,
	rocks       ...image.Rectangle,
) {
	tiles := shatter.Shatter(bounds, rocks...)
	for _, tile := range tiles {
		FillColorRectangle(destination, color, tile)
	}
}

// StrokeColorRectangle is similar to FillColorRectangle, but it draws an inset
// outline of the given rectangle instead.
func StrokeColorRectangle (
	destination canvas.Canvas,
	color       color.RGBA,
	bounds      image.Rectangle,
	weight      int,
) {
	insetBounds := bounds.Inset(weight)
	if insetBounds.Empty() {
		FillColorRectangle(destination, color, bounds)
		return
	}
	FillColorRectangleShatter(destination, color, bounds, insetBounds)
}
