package shapes

import "math"
import "image"
import "image/color"
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
	
	bounds     := source.Bounds().Intersect(destination.Bounds()).Canon()
	realBounds := source.Bounds()
	if bounds.Empty() { return }
	updatedRegion = bounds

	point := image.Point { }
	for point.Y = bounds.Min.Y; point.Y < bounds.Max.Y; point.Y ++ {
	for point.X = bounds.Min.X; point.X < bounds.Max.X; point.X ++ {
		if inEllipse(point, realBounds) {
			offsetPoint := point.Add(offset)
			dstIndex := offsetPoint.X + (offsetPoint.Y) * dstStride
			srcIndex := point.X + point.Y * srcStride
			dstData[dstIndex] = srcData[srcIndex]
		}
	}}
	return
}

// StrokeEllipse is similar to FillEllipse, but it draws an elliptical inset
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
		plottingContext: plottingContext {
			dstData:   dstData,
			dstStride: dstStride,
			srcData:   srcData,
			srcStride: srcStride,
			weight:    weight,
			offset:    offset,
			bounds:    bounds.Intersect(destination.Bounds()),
		},
		radii:  image.Pt(bounds.Dx() / 2 - 1, bounds.Dy() / 2 - 1),
	}
	context.center = bounds.Min.Add(context.radii)
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
	destination canvas.Canvas,
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
	destination canvas.Canvas,
	color       color.RGBA,
	bounds      image.Rectangle,
	weight      int,
) (
	updatedRegion image.Rectangle,
) {
	if weight < 1 { return }

	dstData, dstStride := destination.Buffer()
	bounds = bounds.Inset(weight - 1)

	context := ellipsePlottingContext {
		plottingContext: plottingContext {
			dstData:   dstData,
			dstStride: dstStride,
			color:     color,
			weight:    weight,
			bounds:    bounds.Intersect(destination.Bounds()),
		},
		radii:  image.Pt(bounds.Dx() / 2 - 1, bounds.Dy() / 2 - 1),
	}
	context.center = bounds.Min.Add(context.radii)
	context.plotEllipse()
	return
}

func inEllipse (point image.Point, bounds image.Rectangle) bool {
	point = point.Sub(bounds.Min)
	x := (float64(point.X) + 0.5) / float64(bounds.Dx()) - 0.5
	y := (float64(point.Y) + 0.5) / float64(bounds.Dy()) - 0.5
	return math.Hypot(x, y) <= 0.5
}
