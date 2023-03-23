package artist

import "image"

// Side represents one side of a rectangle.
type Side int; const (
	SideTop Side = iota
	SideRight
	SideBottom
	SideLeft
)

// Inset represents an inset amount for all four sides of a rectangle. The top
// side is at index zero, the right at index one, the bottom at index two, and
// the left at index three. These values may be negative.
type Inset [4]int

// I allows you to create an inset in a CSS-ish way:
//
//   - One argument: all sides are set to this value
//   - Two arguments: the top and bottom sides are set to the first value, and
//     the left and right sides are set to the second value.
//   - Three arguments: the top side is set by the first value, the left and
//     right sides are set by the second vaue, and the bottom side is set by the
//     third value.
//   - Four arguments: each value corresponds to a side.
//
// This function will panic if an argument count that isn't one of these is
// given.
func I (sides ...int) Inset {
	switch len(sides) {
	case 1: return Inset { sides[0], sides[0], sides[0], sides[0] }
	case 2: return Inset { sides[0], sides[1], sides[0], sides[1] }
	case 3: return Inset { sides[0], sides[1], sides[2], sides[1] }
	case 4: return Inset { sides[0], sides[1], sides[2], sides[3] }
	default: panic("I: illegal argument count.")
	}
}

// Apply returns the given rectangle, shrunk on all four sides by the given
// inset. If a measurment of the inset is negative, that side will instead be
// expanded outward. If the rectangle's dimensions cannot be reduced any
// further, an empty rectangle near its center will be returned.
func (inset Inset) Apply (bigger image.Rectangle) (smaller image.Rectangle) {
	smaller = bigger
	if smaller.Dx() < inset[3] + inset[1] {
		smaller.Min.X = (smaller.Min.X + smaller.Max.X) / 2
		smaller.Max.X = smaller.Min.X
	} else {
		smaller.Min.X += inset[3]
		smaller.Max.X -= inset[1]
	}

	if smaller.Dy() < inset[0] + inset[2] {
		smaller.Min.Y = (smaller.Min.Y + smaller.Max.Y) / 2
		smaller.Max.Y = smaller.Min.Y
	} else {
		smaller.Min.Y += inset[0]
		smaller.Max.Y -= inset[2]
	}
	return
}

// Inverse returns a negated version of the inset.
func (inset Inset) Inverse () (prime Inset) {
	return Inset {
		inset[0] * -1,
		inset[1] * -1,
		inset[2] * -1,
		inset[3] * -1,
	}
}

// Horizontal returns the sum of SideRight and SideLeft.
func (inset Inset) Horizontal () int {
	return inset[SideRight] + inset[SideLeft]
}

// Vertical returns the sum of SideTop and SideBottom.
func (inset Inset) Vertical () int {
	return inset[SideTop] + inset[SideBottom]
}
