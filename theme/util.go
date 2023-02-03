package theme

import "image/color"
import "git.tebibyte.media/sashakoshka/tomo/artist"

func hex (color uint32) (c color.RGBA) {
	c.A = uint8(color)
	c.B = uint8(color >>  8)
	c.G = uint8(color >> 16)
	c.R = uint8(color >> 24)
	return
}

func uhex (color uint32) (pattern artist.Pattern) {
	return artist.NewUniform(hex(color))
}
