package artist

import "image"
import "git.tebibyte.media/sashakoshka/tomo"

// Paste transfers one canvas onto another, offset by the specified point.
func Paste (
	destination tomo.Canvas,
	source tomo.Canvas,
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

func FillRectangle (
	destination tomo.Canvas,
	source Pattern,
	bounds image.Rectangle,
) (
	updatedRegion image.Rectangle,
) {
	data, stride := destination.Buffer()
	bounds = bounds.Canon().Intersect(destination.Bounds()).Canon()
	if bounds.Empty() { return }
	updatedRegion = bounds

	width, height := bounds.Dx(), bounds.Dy()
	for y := 0; y < height; y ++ {
	for x := 0; x < width;  x ++ {
		data[x + bounds.Min.X + (y + bounds.Min.Y) * stride] =
			source.AtWhen(x, y, width, height)
	}}
	return
}
