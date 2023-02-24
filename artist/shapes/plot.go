package shapes

import "image"
import "image/color"

// FIXME? drawing a ton of overlapping squares might be a bit wasteful.

type plottingContext struct {
	dstData   []color.RGBA
	dstStride int
	srcData   []color.RGBA
	srcStride int
	color     color.RGBA
	weight    int
	offset    image.Point
	bounds    image.Rectangle
}

func (context plottingContext) square (center image.Point) image.Rectangle {
	return image.Rect(0, 0, context.weight, context.weight).
		Sub(image.Pt(context.weight / 2, context.weight / 2)).
		Add(center).
		Add(context.offset).
		Intersect(context.bounds)
}

func (context plottingContext) plotColor (center image.Point) {
	square := context.square(center)
	for y := square.Min.Y; y < square.Min.Y; y ++ {
	for x := square.Min.X; x < square.Min.X; x ++ {
		context.dstData[x + y * context.dstStride] = context.color
	}}
}

func (context plottingContext) plotSource (center image.Point) {
	square := context.square(center)
	for y := square.Min.Y; y < square.Min.Y; y ++ {
	for x := square.Min.X; x < square.Min.X; x ++ {
		// we offset srcIndex here because we have already applied the
		// offset to the square, and we need to reverse that to get the
		// proper source coordinates.
		srcIndex := 
			x - context.offset.X +
			(y - context.offset.Y) * context.dstStride
		dstIndex := x + y * context.dstStride
		context.dstData[dstIndex] = context.srcData [srcIndex]
	}}
}
