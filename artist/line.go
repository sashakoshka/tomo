package artist

import "image"
import "git.tebibyte.media/sashakoshka/tomo"

func Line (
	destination tomo.Canvas,
	source tomo.Image,
	weight int,
	min image.Point,
	max image.Point,
) (
	updatedRegion image.Rectangle,
) {
	// TODO: respect weight
	
	updatedRegion = image.Rectangle { Min: min, Max: max }.Canon()
	updatedRegion.Max.X ++
	updatedRegion.Max.Y ++
	
	if abs(max.Y - min.Y) <
		abs(max.X - min.X) {
		
		if max.X < min.X {
			temp := min
			min = max
			max = temp
		}
		lineLow(destination, source, weight, min, max)
	} else {
	
		if max.Y < min.Y {
			temp := min
			min = max
			max = temp
		}
		lineHigh(destination, source, weight, min, max)
	}
	return
}

func lineLow (
	destination tomo.Canvas,
	source tomo.Image,
	weight int,
	min image.Point,
	max image.Point,
) {
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
		destination.SetRGBA(x, y, source.RGBAAt(x, y))
		if D > 0 {
			y += yi
			D += 2 * (deltaY - deltaX)
		} else {
			D += 2 * deltaY
		}
	}
}

func lineHigh (
	destination tomo.Canvas,
	source tomo.Image,
	weight int,
	min image.Point,
	max image.Point,
) {
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
		destination.SetRGBA(x, y, source.RGBAAt(x, y))
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