package shapes

import "math"
import "image"
import "image/color"
// import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/canvas"

// FillEllipse draws the content of one canvas onto another, clipped by an
// ellipse stretched to the bounds of the source canvas. The offset point
// defines where the origin point of the source canvas is positioned in relation
// to the origin point of the destination canvas. To prevent the entire source
// canvas's bounds from being used, it must be cut with canvas.Cut().
func FillEllipse (
	destination canvas.Canvas,
	source      canvas.Canvas,
	offset      image.Point,
) (
	updatedRegion image.Rectangle,
) {
	dstData, dstStride := destination.Buffer()
	srcData, srcStride := source.Buffer()
	
	bounds := source.Bounds()
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
			dstData[x + offset.X + (y + offset.Y) * dstStride] =
				srcData[x + y * srcStride]
		}
	}}
	return
}

// StrokeRectangle is similar to FillEllipse, but it draws an elliptical inset
// outline of the source canvas onto the destination canvas. To prevent the
// entire source canvas's bounds from being used, it must be cut with
// canvas.Cut().
func StrokeEllipse (
	destination canvas.Canvas,
	source      canvas.Canvas,
	offset      image.Point,
	weight      int,
) {
	if weight < 1 { return }

	dstData, dstStride := destination.Buffer()
	srcData, srcStride := source.Buffer()
	
	bounds := source.Bounds().Inset(weight - 1)

	context := ellipsePlottingContext {
		dstData:   dstData,
		dstStride: dstStride,
		srcData:   srcData,
		srcStride: srcStride,
		weight:    weight,
		offset:    offset,
		bounds:    bounds.Intersect(destination.Bounds()),
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
	dstData   []color.RGBA
	dstStride int
	srcData   []color.RGBA
	srcStride int
	weight    int
	offset    image.Point
	bounds    image.Rectangle
}

func (context ellipsePlottingContext) plot (x, y int) {
	square :=
		image.Rect(0, 0, context.weight, context.weight).
		Sub(image.Pt(context.weight / 2, context.weight / 2)).
		Add(image.Pt(x, y)).
		Intersect(context.bounds)
	
	for y := square.Min.Y; y < square.Min.Y; y ++ {
	for x := square.Min.X; x < square.Min.X; x ++ {
		context.dstData[x + y * context.dstStride] =
			context.srcData [
				x + y * context.dstStride]
	}}
}
