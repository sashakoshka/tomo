package artist

import "math"
import "image"
import "git.tebibyte.media/sashakoshka/tomo"

// FillEllipse draws a filled ellipse with the specified pattern.
func FillEllipse (
	destination tomo.Canvas,
	source Pattern,
	bounds image.Rectangle,
) (
	updatedRegion image.Rectangle,
) {
	data, stride := destination.Buffer()
	realWidth, realHeight := bounds.Dx(), bounds.Dy()
	bounds = bounds.Canon().Intersect(destination.Bounds()).Canon()
	if bounds.Empty() { return }
	updatedRegion = bounds

	width,  height := bounds.Dx(), bounds.Dy()
	for y := 0; y < height; y ++ {
	for x := 0; x < width;  x ++ {
		xf := float64(x) / float64(width)  - 0.5
		yf := float64(y) / float64(height) - 0.5
		if math.Sqrt(xf * xf + yf * yf) <= 0.5 {
			data[x + bounds.Min.X + (y + bounds.Min.Y) * stride] =
				source.AtWhen(x, y, realWidth, realHeight)
		}
	}}
	return
}

// TODO: StrokeEllipse
