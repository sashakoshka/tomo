package main

import "math"
import "image"

type World struct {
	Data   []int
	Stride int
}

func (world World) At (position image.Point) int {
	if position.X < 0 { return 0 }
	if position.Y < 0 { return 0 }
	if position.X >= world.Stride { return 0 }
	index := position.X + position.Y * world.Stride
	if index >= len(world.Data) { return 0 }
	return world.Data[index]
}

type Vector struct {
	X, Y float64
}

func (vector Vector) Point () (image.Point) {
	return image.Pt(int(vector.X), int(vector.Y))
}

func (vector Vector) Add (other Vector) Vector {
	return Vector {
		vector.X + other.X,
		vector.Y + other.Y,
	}
}

func (vector Vector) Sub (other Vector) Vector {
	return Vector {
		vector.X - other.X,
		vector.Y - other.Y,
	}
}

func (vector Vector) Mul (by float64) Vector {
	return Vector {
		vector.X * by,
		vector.Y * by,
	}
}

func (vector Vector) Hypot () float64 {
	return math.Hypot(vector.X, vector.Y)
}

type Camera struct {
	Vector
	Angle float64
	Fov   float64
}

func (camera *Camera) Rotate (by float64) {
	camera.Angle += by
	if camera.Angle < 0 { camera.Angle += math.Pi * 2 }
	if camera.Angle > math.Pi * 2  { camera.Angle = 0 }
}

func (camera *Camera) Walk (by float64) {
	delta := camera.Delta()
	camera.X += delta.X * by
	camera.Y += delta.Y * by
}

func (camera *Camera) Strafe (by float64) {
	delta := camera.OffsetDelta()
	camera.X += delta.X * by
	camera.Y += delta.Y * by
}

func (camera *Camera) Delta () Vector {
	return Vector {
		math.Cos(camera.Angle),
		math.Sin(camera.Angle),
	}
}

func (camera *Camera) OffsetDelta () Vector {
	offset := math.Pi / 2
	return Vector {
		math.Cos(camera.Angle + offset),
		math.Sin(camera.Angle + offset),
	}
}

type Ray struct {
	Vector
	Angle  float64
}

func (ray *Ray) Cast (
	world World,
	max int,
) (
	distance float64,
	hit Vector,
	wall int,
	horizontal bool,
) {
	// return ray.castV(world, max)
	cellAt := world.At(ray.Point())
	if cellAt > 0 {
		return 0, Vector { }, cellAt, false
	}
	hDistance, hPos, hWall := ray.castH(world, max)
	vDistance, vPos, vWall := ray.castV(world, max)
	if hDistance < vDistance {
		return hDistance, hPos, hWall, true
	} else {
		return vDistance, vPos, vWall, false
	}
}

func (ray *Ray) castH (world World, max int) (distance float64, hit Vector, wall int) {
	var position Vector
	var delta    Vector
	var offset   Vector
	ray.Angle = math.Mod(ray.Angle, math.Pi * 2)
	if ray.Angle < 0 {
		ray.Angle += math.Pi * 2
	}
	tan := math.Tan(math.Pi - ray.Angle)
	if ray.Angle > math.Pi {
		// facing up
		position.Y = math.Floor(ray.Y)
		delta.Y  = -1
		offset.Y = -1
	} else if ray.Angle < math.Pi {
		// facing down
		position.Y = math.Floor(ray.Y) + 1
		delta.Y = 1
	} else {
		// facing straight left or right
		return float64(max), Vector { }, 0
	}
	position.X = ray.X + (ray.Y - position.Y) / tan
	delta.X    = -delta.Y / tan

	// cast da ray
	steps := 0
	for {
		cell := world.At(position.Add(offset).Point())
		if cell > 0 || steps > max { break }
		position = position.Add(delta)
		steps ++
	}

	return position.Sub(ray.Vector).Hypot(),
		position,
		world.At(position.Add(offset).Point())
}

func (ray *Ray) castV (world World, max int) (distance float64, hit Vector, wall int) {
	var position Vector
	var delta    Vector
	var offset   Vector
	tan := math.Tan(math.Pi - ray.Angle)
	offsetAngle := math.Mod(ray.Angle + math.Pi / 2, math.Pi * 2)
	if offsetAngle > math.Pi {
		// facing left
		position.X = math.Floor(ray.X)
		delta.X  = -1
		offset.X = -1
	} else if offsetAngle < math.Pi {
		// facing right
		position.X = math.Floor(ray.X) + 1
		delta.X = 1
	} else {
		// facing straight left or right
		return float64(max), Vector { }, 0
	}
	position.Y = ray.Y + (ray.X - position.X) * tan
	delta.Y    = -delta.X * tan

	// cast da ray
	steps := 0
	for {
		cell := world.At(position.Add(offset).Point())
		if cell > 0 || steps > max { break }
		position = position.Add(delta)
		steps ++
	}

	return position.Sub(ray.Vector).Hypot(),
		position,
		world.At(position.Add(offset).Point())
}
