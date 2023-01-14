package artist

import "image"
import "image/color"

// Pattern is capable of generating a pattern pixel by pixel.
type Pattern interface {
	// AtWhen returns the color of the pixel located at (x, y) relative to
	// the origin point of the pattern (0, 0), when the pattern has the
	// specified width and height. Patterns may ignore the width and height
	// parameters, but it may be useful for some patterns such as gradients.
	AtWhen (x, y, width, height int) (color.RGBA)
}

// Texture is a struct that allows an image to be converted into a tiling
// texture pattern.
type Texture struct {
	data []color.RGBA
	width, height int
}

// NewTexture converts an image into a texture.
func NewTexture (source image.Image) (texture Texture) {
	bounds := source.Bounds()
	texture.width  = bounds.Dx()
	texture.height = bounds.Dy()
	texture.data   = make([]color.RGBA, texture.width * texture.height)

	index := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y ++ {
	for x := bounds.Min.X; x < bounds.Max.X; x ++ {
		r, g, b, a := source.At(x, y).RGBA()
		texture.data[index] = color.RGBA {
			uint8(r >> 8),
			uint8(g >> 8),
			uint8(b >> 8),
			uint8(a >> 8),
		}
		index ++
	}}
	return
}

// AtWhen returns the color at the specified x and y coordinates, wrapped to the
// image's width. the width and height are ignored.
func (texture Texture) AtWhen (x, y, width, height int) (pixel color.RGBA) {
	x %= texture.width
	y %= texture.height
	if x < 0 { x += texture.width  }
	if y < 0 { y += texture.height }
	return texture.data[x + y * texture.width]
}
