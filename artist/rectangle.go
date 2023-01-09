package artist

import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo"

// Paste transfers one image onto another, offset by the specified point.
func Paste (
	destination tomo.Canvas,
	source tomo.Image,
	offset image.Point,
) (
	updatedRegion image.Rectangle,
) {
	sourceBounds := source.Bounds().Canon()
	updatedRegion = sourceBounds.Add(offset)
	for y := sourceBounds.Min.Y; y < sourceBounds.Max.Y; y ++ {
	for x := sourceBounds.Min.X; x < sourceBounds.Max.X; x ++ {
		destination.SetRGBA (
			x + offset.X, y + offset.Y,
			source.RGBAAt(x, y))
	}}

	return
}

// Rectangle draws a rectangle with an inset border. If the border image is nil,
// no border will be drawn. Likewise, if the fill image is nil, the rectangle
// will have no fill.
func Rectangle (
	destination tomo.Canvas,
	fill   tomo.Image,
	stroke tomo.Image,
	weight int,
	bounds image.Rectangle,
) (
	updatedRegion image.Rectangle,
) {
	bounds = bounds.Canon()
	updatedRegion = bounds

	fillBounds := bounds
	fillBounds.Min = fillBounds.Min.Add(image.Point { weight, weight })
	fillBounds.Max = fillBounds.Max.Sub(image.Point { weight, weight })
	fillBounds = fillBounds.Canon()

	for y := bounds.Min.Y; y < bounds.Max.Y; y ++ {
	for x := bounds.Min.X; x < bounds.Max.X; x ++ {
		var pixel color.RGBA
		if (image.Point { x, y }).In(fillBounds) {
			pixel = fill.RGBAAt(x, y)
		} else {
			pixel = stroke.RGBAAt(x, y)
		}
		destination.SetRGBA(x, y, pixel)
	}}
	
	return
}

// OffsetRectangle is the same as Rectangle, but offsets the border image to the
// top left corner of the border and the fill image to the top left corner of
// the fill.
func OffsetRectangle (
	destination tomo.Canvas,
	fill   tomo.Image,
	stroke tomo.Image,
	weight int,
	bounds image.Rectangle,
) (
	updatedRegion image.Rectangle,
) {
	bounds = bounds.Canon()
	updatedRegion = bounds

	fillBounds := bounds
	fillBounds.Min = fillBounds.Min.Add(image.Point { weight, weight })
	fillBounds.Max = fillBounds.Max.Sub(image.Point { weight, weight })
	fillBounds = fillBounds.Canon()

	strokeImageMin := stroke.Bounds().Min
	fillImageMin   := fill.Bounds().Min

	yy := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y ++ {
		xx := 0
		for x := bounds.Min.X; x < bounds.Max.X; x ++ {
			var pixel color.RGBA
			if (image.Point { x, y }).In(fillBounds) {
				pixel = fill.RGBAAt (
					xx - weight + fillImageMin.X,
					yy - weight + fillImageMin.Y)
			} else {
				pixel = stroke.RGBAAt (
					xx + strokeImageMin.X,
					yy + strokeImageMin.Y)
			}
			destination.SetRGBA(x, y, pixel)
			xx ++
		}
		yy ++
	}
	
	return
}
