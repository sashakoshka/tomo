package artist

import "image/color"

// Noisy is a pattern that randomly interpolates between two patterns in a
// deterministic fashion.
type Noisy struct {
	Low   Pattern
	High  Pattern
	Seed  uint32
	Harsh bool
}

// AtWhen satisfies the pattern interface.
func (pattern Noisy) AtWhen (x, y, width, height int) (c color.RGBA) {
	// FIXME: this will occasionally generate "clumps"
	special := uint32(x + y * 348905)
	special += (pattern.Seed + 1) * 15485863
	random := (special * special * special % 2038074743)
	fac := float64(random) / 2038074743.0

	if pattern.Harsh {
		if fac > 0.5 {
			return pattern.High.AtWhen(x, y, width, height)
		} else {
			return pattern.Low.AtWhen(x, y, width, height)
		}
	} else {
		return LerpRGBA (
			pattern.Low.AtWhen(x, y, width, height),
			pattern.High.AtWhen(x, y, width, height), fac)
	}
}
