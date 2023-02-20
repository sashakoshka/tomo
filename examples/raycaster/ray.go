package main

import "math"
import "image"

type World interface {
	At (image.Point) int
}

type DefaultWorld struct {
	Data   []int
	Stride int
}

func (world DefaultWorld) At (position image.Point) int {
	if position.X < 0 { return 0 }
	if position.Y < 0 { return 0 }
	if position.X >= world.Stride { return 0 }
	index := position.X + position.Y * world.Stride
	if index >= len(world.Data) { return 0 }
	return world.Data[index]
}

type Camera struct {
	X, Y  float64
	Angle float64
	Fov   float64
}

func (camera *Camera) Point () (image.Point) {
	return image.Pt(int(camera.X), int(camera.Y))
}

func (camera *Camera) Rotate (by float64) {
	camera.Angle += by
	if camera.Angle < 0 { camera.Angle += math.Pi * 2 }
	if camera.Angle > math.Pi * 2  { camera.Angle = 0 }
}

func (camera *Camera) Walk (by float64) {
	dx, dy := camera.Delta()
	camera.X += dx * by
	camera.Y += dy * by
}

func (camera *Camera) Strafe (by float64) {
	dx, dy := camera.OffsetDelta()
	camera.X += dx * by
	camera.Y += dy * by
}

func (camera *Camera) Delta () (x float64, y float64) {
	return math.Cos(camera.Angle), math.Sin(camera.Angle)
}

func (camera *Camera) OffsetDelta () (x float64, y float64) {
	offset := math.Pi / 2
	return math.Cos(camera.Angle + offset), math.Sin(camera.Angle + offset)
}

type Ray struct {
	X, Y   float64
	Angle  float64
	Precision int
}

func (ray *Ray) Cast (world World, max int) (distance float64) {
	precision := 64

	dX := math.Cos(ray.Angle) / float64(precision)
	dY := math.Sin(ray.Angle) / float64(precision)
	origX, origY := ray.X, ray.Y

	wall  := 0
	depth := 0
	for wall == 0 && depth < max * precision {
		ray.X += dX
		ray.Y += dY
		wall = world.At(ray.Point())
		depth ++
	}

	distanceX := origX - ray.X
	distanceY := origY - ray.Y
	return math.Sqrt(distanceX * distanceX + distanceY * distanceY)
}

func (ray *Ray) Point () (image.Point) {
	return image.Pt(int(ray.X), int(ray.Y))
}
