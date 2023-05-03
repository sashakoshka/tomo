package shapes

import "math"
import "image"
import "image/color"
import "tomo/artist"

// TODO: redo fill ellipse, stroke ellipse, etc. so that it only takes in
// destination and source, using the bounds of destination as the bounds of the
// ellipse and the bounds of source as the "clipping rectangle". Line up the Min
// of both canvases.

func FillEllipse (
	destination artist.Canvas,
	source      artist.Canvas,
	bounds      image.Rectangle,
) (
	updatedRegion image.Rectangle,
) {
	dstData, dstStride := destination.Buffer()
	srcData, srcStride := source.Buffer()
	
	offset := source.Bounds().Min.Sub(destination.Bounds().Min)
	drawBounds :=
		source.Bounds().Sub(offset).
		Intersect(destination.Bounds()).
		Intersect(bounds)
	if bounds.Empty() { return }
	updatedRegion = bounds

	point := image.Point { }
	for point.Y = drawBounds.Min.Y; point.Y < drawBounds.Max.Y; point.Y ++ {
	for point.X = drawBounds.Min.X; point.X < drawBounds.Max.X; point.X ++ {
		if inEllipse(point, bounds) {
			offsetPoint := point.Add(offset)
			dstIndex := point.X       + point.Y       * dstStride
			srcIndex := offsetPoint.X + offsetPoint.Y * srcStride
			dstData[dstIndex] = srcData[srcIndex]
		}
	}}
	return
}

func StrokeEllipse (
	destination artist.Canvas,
	source      artist.Canvas,
	bounds      image.Rectangle,
	weight      int,
) {
	if weight < 1 { return }

	dstData, dstStride := destination.Buffer()
	srcData, srcStride := source.Buffer()
	
	drawBounds := destination.Bounds().Inset(weight - 1)
	offset := source.Bounds().Min.Sub(destination.Bounds().Min)
	if drawBounds.Empty() { return }

	context := ellipsePlottingContext {
		plottingContext: plottingContext {
			dstData:   dstData,
			dstStride: dstStride,
			srcData:   srcData,
			srcStride: srcStride,
			weight:    weight,
			offset:    offset,
			bounds:    bounds,
		},
		radii:  image.Pt(drawBounds.Dx() / 2, drawBounds.Dy() / 2),
	}
	context.center = drawBounds.Min.Add(context.radii)
	context.plotEllipse()
}

type ellipsePlottingContext struct {
	plottingContext
	radii  image.Point
	center image.Point
}

func (context ellipsePlottingContext) plotEllipse () {
	x := float64(0)
	y := float64(context.radii.Y)

	// region 1 decision parameter
	decision1 :=
		float64(context.radii.Y * context.radii.Y) -
		float64(context.radii.X * context.radii.X * context.radii.Y) +
		(0.25 * float64(context.radii.X) * float64(context.radii.X))
	decisionX := float64(2 * context.radii.Y * context.radii.Y * int(x))
	decisionY := float64(2 * context.radii.X * context.radii.X * int(y))

	// draw region 1
	for decisionX < decisionY {
		points := []image.Point {
			image.Pt(-int(x) + context.center.X, -int(y) + context.center.Y),
			image.Pt( int(x) + context.center.X, -int(y) + context.center.Y),
			image.Pt(-int(x) + context.center.X,  int(y) + context.center.Y),
			image.Pt( int(x) + context.center.X,  int(y) + context.center.Y),			
		}
		if context.srcData == nil {
			context.plotColor(points[0])
			context.plotColor(points[1])
			context.plotColor(points[2])
			context.plotColor(points[3])
		} else {
			context.plotSource(points[0])
			context.plotSource(points[1])
			context.plotSource(points[2])
			context.plotSource(points[3])
		}

		if (decision1 < 0) {
			x ++
			decisionX += float64(2 * context.radii.Y * context.radii.Y)
			decision1 += decisionX + float64(context.radii.Y * context.radii.Y)
		} else {
			x ++
			y --
			decisionX += float64(2 * context.radii.Y * context.radii.Y)
			decisionY -= float64(2 * context.radii.X * context.radii.X)
			decision1 +=
				decisionX - decisionY +
				float64(context.radii.Y * context.radii.Y)
		}
	}

	// region 2 decision parameter
	decision2 :=
		float64(context.radii.Y * context.radii.Y) * (x + 0.5) * (x + 0.5) +
		float64(context.radii.X * context.radii.X) * (y - 1)   * (y - 1) -
		float64(context.radii.X * context.radii.X * context.radii.Y * context.radii.Y)

	// draw region 2
	for y >= 0 {
		points := []image.Point {
			image.Pt( int(x) + context.center.X,  int(y) + context.center.Y),
			image.Pt(-int(x) + context.center.X,  int(y) + context.center.Y),
			image.Pt( int(x) + context.center.X, -int(y) + context.center.Y),
			image.Pt(-int(x) + context.center.X, -int(y) + context.center.Y),
		}
		if context.srcData == nil {
			context.plotColor(points[0])
			context.plotColor(points[1])
			context.plotColor(points[2])
			context.plotColor(points[3])
		} else {
			context.plotSource(points[0])
			context.plotSource(points[1])
			context.plotSource(points[2])
			context.plotSource(points[3])
		}

		if decision2 > 0 {
			y --
			decisionY -= float64(2 * context.radii.X * context.radii.X)
			decision2 += float64(context.radii.X * context.radii.X) - decisionY
		} else {
			y --
			x ++
			decisionX += float64(2 * context.radii.Y * context.radii.Y)
			decisionY -= float64(2 * context.radii.X * context.radii.X)
			decision2 +=
				decisionX - decisionY +
				float64(context.radii.X * context.radii.X)
		}
	}
}

// FillColorEllipse fills an ellipse within the destination canvas with a solid
// color.
func FillColorEllipse (
	destination artist.Canvas,
	color       color.RGBA,
	bounds      image.Rectangle,
) (
	updatedRegion image.Rectangle,
) {
	dstData, dstStride := destination.Buffer()
	
	realBounds := bounds
	bounds = bounds.Intersect(destination.Bounds()).Canon()
	if bounds.Empty() { return }
	updatedRegion = bounds

	point := image.Point { }
	for point.Y = bounds.Min.Y; point.Y < bounds.Max.Y; point.Y ++ {
	for point.X = bounds.Min.X; point.X < bounds.Max.X; point.X ++ {
		if inEllipse(point, realBounds) {
			dstData[point.X + point.Y * dstStride] = color
		}
	}}
	return
}

// StrokeColorEllipse is similar to FillColorEllipse, but it draws an inset
// outline of an ellipse instead.
func StrokeColorEllipse (
	destination artist.Canvas,
	color       color.RGBA,
	bounds      image.Rectangle,
	weight      int,
) (
	updatedRegion image.Rectangle,
) {
	if weight < 1 { return }

	dstData, dstStride := destination.Buffer()
	insetBounds := bounds.Inset(weight - 1)

	context := ellipsePlottingContext {
		plottingContext: plottingContext {
			dstData:   dstData,
			dstStride: dstStride,
			color:     color,
			weight:    weight,
			bounds:    bounds.Intersect(destination.Bounds()),
		},
		radii: image.Pt(insetBounds.Dx() / 2, insetBounds.Dy() / 2),
	}
	context.center = insetBounds.Min.Add(context.radii)
	context.plotEllipse()
	return
}

func inEllipse (point image.Point, bounds image.Rectangle) bool {
	point = point.Sub(bounds.Min)
	x := (float64(point.X) + 0.5) / float64(bounds.Dx()) - 0.5
	y := (float64(point.Y) + 0.5) / float64(bounds.Dy()) - 0.5
	return math.Hypot(x, y) <= 0.5
}
