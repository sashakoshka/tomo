package shapes

import "image"
import "image/color"
import "tomo/artist"
import "tomo/shatter"

// TODO: return updatedRegion for all routines in this package

func FillRectangle (
	destination artist.Canvas,
	source      artist.Canvas,
	bounds      image.Rectangle,
) (
	updatedRegion image.Rectangle,
) {
	dstData, dstStride := destination.Buffer()
	srcData, srcStride := source.Buffer()

	offset     := source.Bounds().Min.Sub(destination.Bounds().Min)
	drawBounds :=
		source.Bounds().Sub(offset).
		Intersect(destination.Bounds()).
		Intersect(bounds)
	if drawBounds.Empty() { return }
	updatedRegion = drawBounds
	
	point := image.Point { }
	for point.Y = drawBounds.Min.Y; point.Y < drawBounds.Max.Y; point.Y ++ {
	for point.X = drawBounds.Min.X; point.X < drawBounds.Max.X; point.X ++ {
		offsetPoint := point.Add(offset)
		dstIndex := point.X       + point.Y       * dstStride
		srcIndex := offsetPoint.X + offsetPoint.Y * srcStride
		dstData[dstIndex] = srcData[srcIndex]
	}}

	return
}

func StrokeRectangle (
	destination artist.Canvas,
	source      artist.Canvas,
	bounds      image.Rectangle,
	weight      int,
) (
	updatedRegion image.Rectangle,
) {
	insetBounds := bounds.Inset(weight)
	if insetBounds.Empty() {
		return FillRectangle(destination, source, bounds)
	}
	return FillRectangleShatter(destination, source, bounds, insetBounds)
}

// FillRectangleShatter is like FillRectangle, but it does not draw in areas
// specified in "rocks".
func FillRectangleShatter (
	destination artist.Canvas,
	source      artist.Canvas,
	bounds      image.Rectangle,
	rocks       ...image.Rectangle,
) (
	updatedRegion image.Rectangle,
) {
	tiles := shatter.Shatter(bounds, rocks...)
	for _, tile := range tiles {
		FillRectangle (
			artist.Cut(destination, tile),
			source, tile)
		updatedRegion = updatedRegion.Union(tile)
	}
	return
}

// FillColorRectangle fills a rectangle within the destination canvas with a
// solid color.
func FillColorRectangle (
	destination artist.Canvas,
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
	destination artist.Canvas,
	color       color.RGBA,
	bounds      image.Rectangle,
	rocks       ...image.Rectangle,
) (
	updatedRegion image.Rectangle,
) {
	tiles := shatter.Shatter(bounds, rocks...)
	for _, tile := range tiles {
		FillColorRectangle(destination, color, tile)
		updatedRegion = updatedRegion.Union(tile)
	}
	return
}

// StrokeColorRectangle is similar to FillColorRectangle, but it draws an inset
// outline of the given rectangle instead.
func StrokeColorRectangle (
	destination artist.Canvas,
	color       color.RGBA,
	bounds      image.Rectangle,
	weight      int,
) (
	updatedRegion image.Rectangle,
) {
	insetBounds := bounds.Inset(weight)
	if insetBounds.Empty() {
		return FillColorRectangle(destination, color, bounds)
	}
	return FillColorRectangleShatter(destination, color, bounds, insetBounds)
}
