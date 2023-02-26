package shapes

import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/shatter"

// TODO: return updatedRegion for all routines in this package

func FillRectangle (
	destination canvas.Canvas,
	source      canvas.Canvas,
) (
	updatedRegion image.Rectangle,
) {
	dstData, dstStride := destination.Buffer()
	srcData, srcStride := source.Buffer()

	offset := source.Bounds().Min.Sub(destination.Bounds().Min)
	bounds     := source.Bounds().Sub(offset).Intersect(destination.Bounds())
	if bounds.Empty() { return }
	updatedRegion = bounds
	
	point := image.Point { }
	for point.Y = bounds.Min.Y; point.Y < bounds.Max.Y; point.Y ++ {
	for point.X = bounds.Min.X; point.X < bounds.Max.X; point.X ++ {
		offsetPoint := point.Add(offset)
		dstIndex := point.X       + point.Y       * dstStride
		srcIndex := offsetPoint.X + offsetPoint.Y * srcStride
		dstData[dstIndex] = srcData[srcIndex]
	}}

	return
}

func StrokeRectangle (
	destination canvas.Canvas,
	source      canvas.Canvas,
	weight      int,
) {
	bounds := destination.Bounds()
	insetBounds := bounds.Inset(weight)
	if insetBounds.Empty() {
		FillRectangle(destination, source)
		return
	}
	FillRectangleShatter(destination, source, insetBounds)
}

// FillRectangleShatter is like FillRectangle, but it does not draw in areas
// specified in "rocks".
func FillRectangleShatter (
	destination canvas.Canvas,
	source      canvas.Canvas,
	rocks       ...image.Rectangle,
) {
	tiles  := shatter.Shatter(destination.Bounds(), rocks...)
	offset := source.Bounds().Min.Sub(destination.Bounds().Min)
	for _, tile := range tiles {
		FillRectangle (
			canvas.Cut(destination, tile),
			canvas.Cut(source, tile.Add(offset)))
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
