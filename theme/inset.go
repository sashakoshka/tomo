package theme

import "image"

// Inset represents an inset amount for all four sides of a rectangle. The top
// side is at index zero, the right at index one, the bottom at index two, and
// the left at index three. These values may be negative.
type Inset [4]int

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
