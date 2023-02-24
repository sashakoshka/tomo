package shapes

import "math"
import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/canvas"

// FillEllipse draws a filled ellipse with the specified pattern.
func FillEllipse (
	destination canvas.Canvas,
	source artist.Pattern,
	bounds image.Rectangle,
) (
	updatedRegion image.Rectangle,
) {
	bounds = bounds.Canon()
	data, stride := destination.Buffer()
	realWidth, realHeight := bounds.Dx(), bounds.Dy()
	bounds = bounds.Intersect(destination.Bounds()).Canon()
	if bounds.Empty() { return }
	updatedRegion = bounds

	width, height := bounds.Dx(), bounds.Dy()
	for y := 0; y < height; y ++ {
	for x := 0; x < width;  x ++ {
		xf := (float64(x) + 0.5) / float64(realWidth)  - 0.5
		yf := (float64(y) + 0.5) / float64(realHeight) - 0.5
		if math.Sqrt(xf * xf + yf * yf) <= 0.5 {
			data[x + bounds.Min.X + (y + bounds.Min.Y) * stride] =
				source.AtWhen(x, y, realWidth, realHeight)
		}
	}}
	return
}

// StrokeEllipse draws the outline of an ellipse with the specified line weight
// and pattern.
func StrokeEllipse (
	destination canvas.Canvas,
	source artist.Pattern,
	weight int,
	bounds image.Rectangle,
) {
	if weight < 1 { return }

	data, stride := destination.Buffer()
	bounds = bounds.Canon().Inset(weight - 1)
	width, height := bounds.Dx(), bounds.Dy()

	context := ellipsePlottingContext {
		data: data,
		stride: stride,
		source: source,
		width: width,
		height: height,
		weight: weight,
		bounds: bounds,
	}
	
	bounds.Max.X -= 1
	bounds.Max.Y -= 1

	radii := image.Pt (
		bounds.Dx() / 2,
		bounds.Dy() / 2)
	center := bounds.Min.Add(radii)

	x := float64(0)
	y := float64(radii.Y)

	// region 1 decision parameter
	decision1 :=
		float64(radii.Y * radii.Y) -
		float64(radii.X * radii.X * radii.Y) +
		(0.25 * float64(radii.X) * float64(radii.X))
	decisionX := float64(2 * radii.Y * radii.Y * int(x))
	decisionY := float64(2 * radii.X * radii.X * int(y))

	// draw region 1
	for decisionX < decisionY {
		context.plot( int(x) + center.X,  int(y) + center.Y)
		context.plot(-int(x) + center.X,  int(y) + center.Y)
		context.plot( int(x) + center.X, -int(y) + center.Y)
		context.plot(-int(x) + center.X, -int(y) + center.Y)

		if (decision1 < 0) {
			x ++
			decisionX += float64(2 * radii.Y * radii.Y)
			decision1 += decisionX + float64(radii.Y * radii.Y)
		} else {
			x ++
			y --
			decisionX += float64(2 * radii.Y * radii.Y)
			decisionY -= float64(2 * radii.X * radii.X)
			decision1 +=
				decisionX - decisionY +
				float64(radii.Y * radii.Y)
		}
	}

	// region 2 decision parameter
	decision2 :=
		float64(radii.Y * radii.Y) * (x + 0.5) * (x + 0.5) +
		float64(radii.X * radii.X) * (y - 1)   * (y - 1) -
		float64(radii.X * radii.X * radii.Y * radii.Y)

	// draw region 2
	for y >= 0 {
		context.plot( int(x) + center.X,  int(y) + center.Y)
		context.plot(-int(x) + center.X,  int(y) + center.Y)
		context.plot( int(x) + center.X, -int(y) + center.Y)
		context.plot(-int(x) + center.X, -int(y) + center.Y)

		if decision2 > 0 {
			y --
			decisionY -= float64(2 * radii.X * radii.X)
			decision2 += float64(radii.X * radii.X) - decisionY
		} else {
			y --
			x ++
			decisionX += float64(2 * radii.Y * radii.Y)
			decisionY -= float64(2 * radii.X * radii.X)
			decision2 +=
				decisionX - decisionY +
				float64(radii.X * radii.X)
		}
	}
}

type ellipsePlottingContext struct {
	data []color.RGBA
	stride int
	source artist.Pattern
	width, height int
	weight int
	bounds image.Rectangle
}

func (context ellipsePlottingContext) plot (x, y int) {
	if (image.Point { x, y }).In(context.bounds) {
		squareAround (
			context.data, context.stride, context.source, x, y,
			context.width, context.height, context.weight)
	}
}
