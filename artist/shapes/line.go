package shapes

import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo/canvas"

// TODO: draw thick lines more efficiently

// ColorLine draws a line from one point to another with the specified weight
// and color.
func ColorLine (
	destination canvas.Canvas,
	color       color.RGBA,
	weight      int,
	min         image.Point,
	max         image.Point,
) (
	updatedRegion image.Rectangle,
) {
	
	updatedRegion = image.Rectangle { Min: min, Max: max }.Canon()
	updatedRegion.Max.X ++
	updatedRegion.Max.Y ++
	
	data, stride := destination.Buffer()
	bounds := destination.Bounds()
	context := linePlottingContext {
		dstData:   data,
		dstStride: stride,
		color:     color,
		weight:    weight,
		bounds:    bounds,
		min:       min,
		max:       max,
	}
	
	if abs(max.Y - min.Y) < abs(max.X - min.X) {
		if max.X < min.X { context.swap() }
		context.lineLow()
		
	} else {
		if max.Y < min.Y { context.swap() }
		context.lineHigh()
	}
	return
}

type linePlottingContext struct {
	dstData   []color.RGBA
	dstStride int
	color     color.RGBA
	weight    int
	bounds    image.Rectangle
	min       image.Point
	max       image.Point
}

func (context *linePlottingContext) swap () {
	temp := context.max
	context.max = context.min
	context.min = temp
}

func (context linePlottingContext) lineLow () {
	deltaX := context.max.X - context.min.X
	deltaY := context.max.Y - context.min.Y
	yi     := 1

	if deltaY < 0 {
		yi      = -1
		deltaY *= -1
	}

	D := (2 * deltaY) - deltaX
	point := context.min

	for ; point.X < context.max.X; point.X ++ {
		if !point.In(context.bounds) { break }
		context.plot(point)
		if D > 0 {
			D += 2 * (deltaY - deltaX)
			point.Y += yi
		} else {
			D += 2 * deltaY
		}
	}
}

func (context linePlottingContext) lineHigh () {
	deltaX := context.max.X - context.min.X
	deltaY := context.max.Y - context.min.Y
	xi     := 1

	if deltaX < 0 {
		xi      = -1
		deltaX *= -1
	}

	D := (2 * deltaX) - deltaY
	point := context.min

	for ; point.Y < context.max.Y; point.Y ++ {
		if !point.In(context.bounds) { break }
		context.plot(point)
		if D > 0 {
			point.X += xi
			D += 2 * (deltaX - deltaY)
		} else {
			D += 2 * deltaX
		}
	}
}

func abs (n int) int {
	if n < 0 { n *= -1}
	return n
}

func (context linePlottingContext) plot (center image.Point) {
	square :=
		image.Rect(0, 0, context.weight, context.weight).
		Sub(image.Pt(context.weight / 2, context.weight / 2)).
		Add(center).
		Intersect(context.bounds)

	for y := square.Min.Y; y < square.Min.Y; y ++ {
	for x := square.Min.X; x < square.Min.X; x ++ {
		context.dstData[x + y * context.dstStride] = context.color
	}}
}
