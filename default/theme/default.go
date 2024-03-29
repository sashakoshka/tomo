package theme

import "image"
import "bytes"
import _ "embed"
import _ "image/png"
import "image/color"
import "golang.org/x/image/font"
import "golang.org/x/image/font/basicfont"
import "tomo"
import "tomo/data"
import "art"
import "art/artutil"
import "art/patterns"

//go:embed assets/default.png
var defaultAtlasBytes []byte
var defaultAtlas      art.Canvas
var defaultTextures   [7][7]art.Pattern
//go:embed assets/wintergreen-icons-small.png
var defaultIconsSmallAtlasBytes []byte
var defaultIconsSmall [640]binaryIcon
//go:embed assets/wintergreen-icons-large.png
var defaultIconsLargeAtlasBytes []byte
var defaultIconsLarge [640]binaryIcon

func atlasCell (col, row int, border art.Inset) {
	bounds := image.Rect(0, 0, 8, 8).Add(image.Pt(col, row).Mul(8))
	defaultTextures[col][row] = patterns.Border {
		Canvas: art.Cut(defaultAtlas, bounds),
		Inset:  border,
	}
}

func atlasCol (col int, border art.Inset) {
	for index, _ := range defaultTextures[col] {
		atlasCell(col, index, border)
	}
}

type binaryIcon struct {
	data   []bool
	stride int
}

func (icon binaryIcon) Draw (destination art.Canvas, color color.RGBA, at image.Point) {
	bounds := icon.Bounds().Add(at).Intersect(destination.Bounds())
	point := image.Point { }
	data, stride := destination.Buffer()

	for point.Y = bounds.Min.Y; point.Y < bounds.Max.Y; point.Y ++ {
	for point.X = bounds.Min.X; point.X < bounds.Max.X; point.X ++ {
		srcPoint := point.Sub(at)
		srcIndex := srcPoint.X + srcPoint.Y * icon.stride
		dstIndex := point.X + point.Y * stride
		if icon.data[srcIndex] {
			data[dstIndex] = color
		}
	}}
}

func (icon binaryIcon) Bounds () image.Rectangle {
	return image.Rect(0, 0, icon.stride, len(icon.data) / icon.stride)
}

func binaryIconFrom (source image.Image, clip image.Rectangle) (icon binaryIcon) {
	bounds := source.Bounds().Intersect(clip)
	if bounds.Empty() { return }
	
	icon.stride = bounds.Dx()
	icon.data   = make([]bool, bounds.Dx() * bounds.Dy())

	point    := image.Point { }
	dstIndex := 0
	for point.Y = bounds.Min.Y; point.Y < bounds.Max.Y; point.Y ++ {
	for point.X = bounds.Min.X; point.X < bounds.Max.X; point.X ++ {
		r, g, b, a := source.At(point.X, point.Y).RGBA()
		if a > 0x8000 && (r + g + b) / 3 < 0x8000 {
			icon.data[dstIndex] = true
		}
		dstIndex ++
	}}
	return
}

func init () {
	defaultAtlasImage, _, _ := image.Decode(bytes.NewReader(defaultAtlasBytes))
	defaultAtlas = art.FromImage(defaultAtlasImage)

	atlasCol(0, art.I(0))
	atlasCol(1, art.I(3))
	atlasCol(2, art.I(1))
	atlasCol(3, art.I(1))
	atlasCol(4, art.I(1))
	atlasCol(5, art.I(3))
	atlasCol(6, art.I(1))

	// set up small icons
	defaultIconsSmallAtlasImage, _, _ := image.Decode (
		bytes.NewReader(defaultIconsSmallAtlasBytes))
	point     := image.Point { }
	iconIndex := 0
	for point.Y = 0; point.Y < 20; point.Y ++ {
	for point.X = 0; point.X < 32; point.X ++ {
		defaultIconsSmall[iconIndex] = binaryIconFrom (
			defaultIconsSmallAtlasImage,
			image.Rect(0, 0, 16, 16).Add(point.Mul(16)))
		iconIndex ++
	}}

	// set up large icons
	defaultIconsLargeAtlasImage, _, _ := image.Decode (
		bytes.NewReader(defaultIconsLargeAtlasBytes))
	point     = image.Point { }
	iconIndex = 0
	for point.Y = 0; point.Y < 8; point.Y ++ {
	for point.X = 0; point.X < 32; point.X ++ {
		defaultIconsLarge[iconIndex] = binaryIconFrom (
			defaultIconsLargeAtlasImage,
			image.Rect(0, 0, 32, 32).Add(point.Mul(32)))
		iconIndex ++
	}}
	iconIndex = 384
	for point.Y = 8; point.Y < 12; point.Y ++ {
	for point.X = 0; point.X < 32; point.X ++ {
		defaultIconsLarge[iconIndex] = binaryIconFrom (
			defaultIconsLargeAtlasImage,
			image.Rect(0, 0, 32, 32).Add(point.Mul(32)))
		iconIndex ++
	}}
}

// Default is the default theme.
type Default struct { }

// FontFace returns the default font face.
func (Default) FontFace (style tomo.FontStyle, size tomo.FontSize, c tomo.Case) font.Face {
	return basicfont.Face7x13
}

// Icon returns an icon from the default set corresponding to the given name.
func (Default) Icon (id tomo.Icon, size tomo.IconSize, c tomo.Case) art.Icon {
	if size == tomo.IconSizeLarge {
		if id < 0 || int(id) >= len(defaultIconsLarge) {
			return nil
		} else {
			return defaultIconsLarge[id]
		}
	} else {
		if id < 0 || int(id) >= len(defaultIconsSmall) {
			return nil
		} else {
			return defaultIconsSmall[id]
		}
	}
}

// MimeIcon returns an icon from the default set corresponding to the given mime.
// type.
func (Default) MimeIcon (data.Mime, tomo.IconSize, tomo.Case) art.Icon {
	// TODO
	return nil
}

// Pattern returns a pattern from the default theme corresponding to the given
// pattern ID.
func (Default) Pattern (id tomo.Pattern, state tomo.State, c tomo.Case) art.Pattern {
	offset := 0; switch {
	case state.Disabled:            offset = 1
	case state.Pressed && state.On: offset = 4
	case state.Focused && state.On: offset = 6
	case state.On:                  offset = 2
	case state.Pressed:             offset = 3
	case state.Focused:             offset = 5
	}

	switch id {
	case tomo.PatternBackground: return patterns.Uhex(0xaaaaaaFF)
	case tomo.PatternDead:       return defaultTextures[0][offset]
	case tomo.PatternRaised:     return defaultTextures[1][offset]
	case tomo.PatternSunken:     return defaultTextures[2][offset]
	case tomo.PatternPinboard:   return defaultTextures[3][offset]
	case tomo.PatternButton:     return defaultTextures[1][offset]
	case tomo.PatternInput:      return defaultTextures[2][offset]
	case tomo.PatternGutter:     return defaultTextures[2][offset]
	case tomo.PatternHandle:     return defaultTextures[3][offset]
	case tomo.PatternLine:       return defaultTextures[0][offset]
	case tomo.PatternMercury:    return defaultTextures[4][offset]
	case tomo.PatternTableHead:  return defaultTextures[5][offset]
	case tomo.PatternTableCell:  return defaultTextures[5][offset]
	case tomo.PatternLamp:       return defaultTextures[6][offset]
	default:                     return patterns.Uhex(0xFF00FFFF)
	}
}

func (Default) Color (id tomo.Color, state tomo.State, c tomo.Case) color.RGBA {
	if state.Disabled { return artutil.Hex(0x444444FF) }
	
	return artutil.Hex (map[tomo.Color] uint32 {
		tomo.ColorBlack:        0x272d24FF,
		tomo.ColorRed:          0x8c4230FF,
		tomo.ColorGreen:        0x69905fFF,
		tomo.ColorYellow:       0x9a973dFF,
		tomo.ColorBlue:         0x3d808fFF,
		tomo.ColorPurple:       0x8c608bFF,
		tomo.ColorCyan:         0x3d8f84FF,
		tomo.ColorWhite:        0xaea894FF,
		tomo.ColorBrightBlack:  0x4f5142FF,
		tomo.ColorBrightRed:    0xbd6f59FF,
		tomo.ColorBrightGreen:  0x8dad84FF,
		tomo.ColorBrightYellow: 0xe2c558FF,
		tomo.ColorBrightBlue:   0x77b1beFF,
		tomo.ColorBrightPurple: 0xc991c8FF,
		tomo.ColorBrightCyan:   0x74c7b7FF,
		tomo.ColorBrightWhite:  0xcfd7d2FF,
	
		tomo.ColorForeground: 0x000000FF,
		tomo.ColorMidground:  0x656565FF,
		tomo.ColorBackground: 0xAAAAAAFF,
		tomo.ColorShadow:     0x000000FF,
		tomo.ColorShine:      0xFFFFFFFF,
		tomo.ColorAccent:     0xff3300FF,
	} [id])
}

// Padding returns the default padding value for the given pattern.
func (Default) Padding (id tomo.Pattern, c tomo.Case) art.Inset {
	switch id {
	case tomo.PatternGutter: return art.I(0)
	case tomo.PatternLine:   return art.I(1)
	default:                 return art.I(6)
	}
}

// Margin returns the default margin value for the given pattern.
func (Default) Margin (id tomo.Pattern, c tomo.Case) image.Point {
	return image.Pt(6, 6)
}

// Hints returns rendering optimization hints for a particular pattern.
// These are optional, but following them may result in improved
// performance.
func (Default) Hints (pattern tomo.Pattern, c tomo.Case) (hints tomo.Hints) {
	return
}

// Sink returns the default sink vector for the given pattern.
func (Default) Sink (pattern tomo.Pattern, c tomo.Case) image.Point {
	return image.Point { 1, 1 }
}
