package artist

import "image/color"

// StripeDirection specifies the direction of stripes.
type StripeDirection int

const (
	StripeDirectionVertical StripeDirection = iota
	StripeDirectionDiagonalRight
	StripeDirectionHorizontal
	StripeDirectionDiagonalLeft
)

// Striped is a pattern that produces stripes of two alternating colors.
type Striped struct {
	First     Pattern
	Second    Pattern
	Direction StripeDirection
	Weight    int
}

// AtWhen satisfies the Pattern interface.
func (pattern Striped) AtWhen (x, y, width, height int) (c color.RGBA) {
	position := 0
	switch pattern.Direction {
	case StripeDirectionVertical:
		position = x
	case StripeDirectionDiagonalRight:
		position = x + y
	case StripeDirectionHorizontal:
		position = y
	case StripeDirectionDiagonalLeft:
		position = x - y
	}

	position %= pattern.Weight * 2
	if position < 0 {
		position += pattern.Weight * 2
	}
	
	if position < pattern.Weight {
		return pattern.First.AtWhen(x, y, width, height)
	} else {
		return pattern.Second.AtWhen(x, y, width, height)
	}
}
