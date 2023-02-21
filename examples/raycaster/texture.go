package main

import "io"
import "image"
import "image/color"

type Textures []Texture

type Texture struct {
	Data   []color.RGBA
	Stride int
}

func (texture Textures) At (wall int, offset Vector) color.RGBA {
	wall --
	if wall < 0 || wall >= len(texture) { return color.RGBA { } }
	image := texture[wall]
	
	xOffset := int(offset.X * float64(image.Stride))
	yOffset := int(offset.Y * float64(len(image.Data) / image.Stride))
	
	index := xOffset + yOffset * image.Stride
	if index <  0               { return color.RGBA { } }
	if index >= len(image.Data) { return color.RGBA { } }
	return image.Data[index]
}

func TextureFrom (source io.Reader) (texture Texture, err error) {
	sourceImage, _, err := image.Decode(source)
	if err != nil { return }
	bounds := sourceImage.Bounds()
	texture.Stride = bounds.Dx()
	texture.Data = make([]color.RGBA, bounds.Dx() * bounds.Dy())

	index := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y ++ {
	for x := bounds.Min.X; x < bounds.Max.X; x ++ {
		r, g, b, a := sourceImage.At(x, y).RGBA()
		texture.Data[index] = color.RGBA {
			R: uint8(r >> 8),
			G: uint8(g >> 8),
			B: uint8(b >> 8),
			A: uint8(a >> 8),
		}
		index ++
	}}
	return texture, nil
}
