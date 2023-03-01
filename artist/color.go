package artist

import "image/color"

// Hex creates a color.RGBA value from an RGBA integer value.
func Hex (color uint32) (c color.RGBA) {
	c.A = uint8(color)
	c.B = uint8(color >>  8)
	c.G = uint8(color >> 16)
	c.R = uint8(color >> 24)
	return
}
