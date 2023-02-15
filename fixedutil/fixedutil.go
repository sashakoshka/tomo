// Package fixedutil contains functions that make working with fixed precision
// values easier.
package fixedutil

import "image"
import "golang.org/x/image/math/fixed"

// Pt creates a fixed point from a regular point.
func Pt (point image.Point) fixed.Point26_6 {
	return fixed.P(point.X, point.Y)
}

// RoundPt rounds a fixed point into a regular point.
func RoundPt (point fixed.Point26_6) image.Point {
	return image.Pt(point.X.Round(), point.Y.Round())
}

// FloorPt creates a regular point from the floor of a fixed point.
func FloorPt (point fixed.Point26_6) image.Point {
	return image.Pt(point.X.Floor(),point.Y.Floor())
}

// CeilPt creates a regular point from the ceiling of a fixed point.
func CeilPt (point fixed.Point26_6) image.Point {
	return image.Pt(point.X.Ceil(),point.Y.Ceil())
}
