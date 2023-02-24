package shapes

import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo/canvas"

// TODO: draw thick lines more efficiently

// Line draws a line from one point to another with the specified weight and
// pattern.
func Line (
	destination canvas.Canvas,
	source      canvas.Canvas,
	weight int,
	min image.Point,
	max image.Point,
) (
	updatedRegion image.Rectangle,
) {
	
	updatedRegion = image.Rectangle { Min: min, Max: max }.Canon()
	updatedRegion.Max.X ++
	updatedRegion.Max.Y ++
	width  := updatedRegion.Dx()
	height := updatedRegion.Dy()
	
	if abs(max.Y - min.Y) <
		abs(max.X - min.X) {
		
		if max.X < min.X {
			temp := min
			min = max
			max = temp
		}
		lineLow(destination, source, weight, min, max, width, height)
	} else {
	
		if max.Y < min.Y {
			temp := min
			min = max
			max = temp
		}
		lineHigh(destination, source, weight, min, max, width, height)
	}
	return
}

func lineLow (
	destination canvas.Canvas,
	source Pattern,
	weight int,
	min image.Point,
	max image.Point,
	width, height int,
) {
	data, stride := destination.Buffer()
	bounds := destination.Bounds()

	deltaX := max.X - min.X
	deltaY := max.Y - min.Y
	yi     := 1

	if deltaY < 0 {
		yi      = -1
		deltaY *= -1
	}

	D := (2 * deltaY) - deltaX
	y := min.Y

	for x := min.X; x < max.X; x ++ {
		if !(image.Point { x, y }).In(bounds) { break }
		squareAround(data, stride, source, x, y, width, height, weight)
		// data[x + y * stride] = source.AtWhen(x, y, width, height)
		if D > 0 {
			y += yi
			D += 2 * (deltaY - deltaX)
		} else {
			D += 2 * deltaY
		}
	}
}

func lineHigh (
	destination canvas.Canvas,
	source Pattern,
	weight int,
	min image.Point,
	max image.Point,
	width, height int,
) {
	data, stride := destination.Buffer()
	bounds := destination.Bounds()

	deltaX := max.X - min.X
	deltaY := max.Y - min.Y
	xi     := 1

	if deltaX < 0 {
		xi      = -1
		deltaX *= -1
	}

	D := (2 * deltaX) - deltaY
	x := min.X

	for y := min.Y; y < max.Y; y ++ {
		if !(image.Point { x, y }).In(bounds) { break }
		squareAround(data, stride, source, x, y, width, height, weight)
		// data[x + y * stride] = source.AtWhen(x, y, width, height)
		if D > 0 {
			x += xi
			D += 2 * (deltaX - deltaY)
		} else {
			D += 2 * deltaX
		}
	}
}

func abs (in int) (out int) {
	if in < 0 { in *= -1}
	out = in
	return
}

// TODO: this method of doing things sucks and can cause a segfault. we should
// not be doing it this way
func squareAround (
	data   []color.RGBA,
	stride int,
	source Pattern,
	x, y, patternWidth, patternHeight, diameter int,
) {
	minY := y - diameter + 1
	minX := x - diameter + 1
	maxY := y + diameter
	maxX := x + diameter
	for y = minY; y < maxY; y ++ {
	for x = minX; x < maxX; x ++ {
		data[x + y * stride] =
			source.AtWhen(x, y, patternWidth, patternHeight)
	}}
}
