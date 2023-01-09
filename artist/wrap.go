package artist

import "git.tebibyte.media/sashakoshka/tomo"

import "image"
import "image/draw"
import "image/color"

// WrappedImage wraps an image.Image and allows it to satisfy tomo.Image.
type WrappedImage struct { Underlying image.Image }

// WrapImage wraps a generic image.Image and allows it to satisfy tomo.Image.
// Do not use this function to wrap images that already satisfy tomo.Image,
// because the resulting wrapped image will be rather slow in comparison.
func WrapImage (underlying image.Image) (wrapped tomo.Image) {
	wrapped = WrappedImage { Underlying: underlying }
	return
}

func (wrapped WrappedImage) Bounds () (bounds image.Rectangle) {
	bounds = wrapped.Underlying.Bounds()
	return
}

func (wrapped WrappedImage) ColorModel () (model color.Model) {
	model = wrapped.Underlying.ColorModel()
	return
}

func (wrapped WrappedImage) At (x, y int) (pixel color.Color) {
	pixel = wrapped.Underlying.At(x, y)
	return
}

func (wrapped WrappedImage) RGBAAt (x, y int) (pixel color.RGBA) {
	r, g, b, a := wrapped.Underlying.At(x, y).RGBA()
	pixel.R = uint8(r >> 8)
	pixel.G = uint8(g >> 8)
	pixel.B = uint8(b >> 8)
	pixel.A = uint8(a >> 8)
	return
}

// WrappedCanvas wraps a draw.Image and allows it to satisfy tomo.Canvas.
type WrappedCanvas struct { Underlying draw.Image }

// WrapCanvas wraps a generic draw.Image and allows it to satisfy tomo.Canvas.
// Do not use this function to wrap images that already satisfy tomo.Canvas,
// because the resulting wrapped image will be rather slow in comparison.
func WrapCanvas (underlying draw.Image) (wrapped tomo.Canvas) {
	wrapped = WrappedCanvas { Underlying: underlying }
	return
}

func (wrapped WrappedCanvas) Bounds () (bounds image.Rectangle) {
	bounds = wrapped.Underlying.Bounds()
	return
}

func (wrapped WrappedCanvas) ColorModel () (model color.Model) {
	model = wrapped.Underlying.ColorModel()
	return
}

func (wrapped WrappedCanvas) At (x, y int) (pixel color.Color) {
	pixel = wrapped.Underlying.At(x, y)
	return
}

func (wrapped WrappedCanvas) RGBAAt (x, y int) (pixel color.RGBA) {
	r, g, b, a := wrapped.Underlying.At(x, y).RGBA()
	pixel.R = uint8(r >> 8)
	pixel.G = uint8(g >> 8)
	pixel.B = uint8(b >> 8)
	pixel.A = uint8(a >> 8)
	return
}

func (wrapped WrappedCanvas) Set (x, y int, pixel color.Color) {
	wrapped.Underlying.Set(x, y, pixel)
}

func (wrapped WrappedCanvas) SetRGBA (x, y int, pixel color.RGBA) {
	wrapped.Underlying.Set(x, y, pixel)
}

// ToRGBA clones an existing image.Image into an image.RGBA struct, which
// directly satisfies tomo.Image. This is useful for things like icons and
// textures.
func ToRGBA (input image.Image) (output *image.RGBA) {
	bounds := input.Bounds()
	output = image.NewRGBA(bounds)
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y ++ {
	for x := bounds.Min.X; x < bounds.Max.X; x ++ {
		output.Set(x, y, input.At(x, y))
	}}
	return
}
