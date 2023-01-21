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
	First     Stroke
	Second    Stroke
	Direction StripeDirection
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

	phase := pattern.First.Weight + pattern.Second.Weight
	position %= phase
	if position < 0 {
		position += phase
	}
	
	if position < pattern.First.Weight {
		return pattern.First.Pattern.AtWhen(x, y, width, height)
	} else {
		return pattern.Second.Pattern.AtWhen(x, y, width, height)
	}
}
