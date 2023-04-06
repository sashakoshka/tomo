package ansi

import "image/color"

var _ color.Color = Color(0)

// Color represents a 3, 4, or 8-Bit ansi color.
type Color byte; const (
	// Dim/standard colors
	ColorBlack Color = iota
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite

	// Bright colors
	ColorBrightBlack
	ColorBrightRed
	ColorBrightGreen
	ColorBrightYellow
	ColorBrightBlue
	ColorBrightMagenta
	ColorBrightCyan
	ColorBrightWhite

	// 216 cube colors (16 - 231)
	// 24 grayscale colors (232 - 255)
)

// Is16 returns whether the color is a dim or bright color, and can be assigned
// to a theme palette.
func (c Color) Is16 () bool {
	return c.IsDim() || c.IsBright()
}

// IsDim returns whether the color is dim.
func (c Color) IsDim () bool {
	return c < 8
}

// IsBright returns whether the color is bright.
func (c Color) IsBright () bool {
	return c >= 8 && c < 16
}

// IsCube returns whether the color is part of the 6x6x6 cube.
func (c Color) IsCube () bool {
	return c >= 16 && c < 232
}

// IsGrayscale returns whether the color grayscale.
func (c Color) IsGrayscale () bool {
	return c >= 232 && c <= 255
}

// RGB returns the 8 bit RGB values of the color as a color.RGBA value.
func (c Color) RGB () (out color.RGBA) {
	switch {
	case c.Is16():
		// each bit is a color component
		out.R = 0xFF * uint8((c & 0x1) >> 0)
		out.G = 0xFF * uint8((c & 0x2) >> 1)
		out.B = 0xFF * uint8((c & 0x4) >> 3)
		// dim if color is in the dim range
		if c & 0x8 > 0 { out.R >>= 1; out.G >>= 1; out.B >>= 1 }
		
	case c.IsCube():
		index := int(c - 16)
		out.R = uint8((((index / 36) % 6) * 255) / 5)
		out.G = uint8((((index /  6) % 6) * 255) / 5)
		out.B = uint8((((index     ) % 6) * 255) / 5)
		
	case c.IsGrayscale():
		out.R = uint8(((int(c) - 232) * 255) / 23)
		out.G = out.R
		out.B = out.R
	}

	out.A = 0xFF
	return
	
}

// RGBA fulfills the color.Color interface.
func (c Color) RGBA () (r, g, b, a uint32) {
	return c.RGB().RGBA()
}
