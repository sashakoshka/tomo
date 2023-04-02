package theme

import "image"
import "bytes"
import _ "embed"
import _ "image/png"
import "image/color"
import "golang.org/x/image/font"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/data"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import defaultfont "git.tebibyte.media/sashakoshka/tomo/default/font"
import "git.tebibyte.media/sashakoshka/tomo/artist/patterns"

//go:embed assets/wintergreen.png
var defaultAtlasBytes []byte
var defaultAtlas      canvas.Canvas
var defaultTextures   [14][9]artist.Pattern
//go:embed assets/wintergreen-icons-small.png
var defaultIconsSmallAtlasBytes []byte
var defaultIconsSmall [640]binaryIcon
//go:embed assets/wintergreen-icons-large.png
var defaultIconsLargeAtlasBytes []byte
var defaultIconsLarge [640]binaryIcon

func atlasCell (col, row int, border artist.Inset) {
	bounds := image.Rect(0, 0, 16, 16).Add(image.Pt(col, row).Mul(16))
	defaultTextures[col][row] = patterns.Border {
		Canvas: canvas.Cut(defaultAtlas, bounds),
		Inset:  border,
	}
}

func atlasCol (col int, border artist.Inset) {
	for index, _ := range defaultTextures[col] {
		atlasCell(col, index, border)
	}
}

type binaryIcon struct {
	data   []bool
	stride int
}

func (icon binaryIcon) Draw (destination canvas.Canvas, color color.RGBA, at image.Point) {
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
	defaultAtlas = canvas.FromImage(defaultAtlasImage)

	// PatternDead
	atlasCol(0, artist.Inset { })
	// PatternRaised
	atlasCol(1, artist.Inset { 6, 6, 6, 6 })
	// PatternSunken
	atlasCol(2, artist.Inset { 4, 4, 4, 4 })
	// PatternPinboard
	atlasCol(3, artist.Inset { 2, 2, 2, 2 })
	// PatternButton
	atlasCol(4, artist.Inset { 6, 6, 6, 6 })
	// PatternInput
	atlasCol(5, artist.Inset { 4, 4, 4, 4 })
	// PatternGutter
	atlasCol(6, artist.Inset { 7, 7, 7, 7 })
	// PatternHandle
	atlasCol(7, artist.Inset { 3, 3, 3, 3 })
	// PatternLine
	atlasCol(8, artist.Inset { 1, 1, 1, 1 })
	// PatternMercury
	atlasCol(13, artist.Inset { 2, 2, 2, 2 })

	// PatternButton: basic.checkbox
	atlasCol(9, artist.Inset { 3, 3, 3, 3 })
	// PatternRaised: basic.listEntry
	atlasCol(10, artist.Inset { 3, 3, 3, 3 })
	// PatternRaised: fun.flatKey
	atlasCol(11, artist.Inset { 3, 3, 5, 3 })
	// PatternRaised: fun.sharpKey
	atlasCol(12, artist.Inset { 3, 3, 4, 3 })

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
	switch style {
	case tomo.FontStyleBold:
		return defaultfont.FaceBold
	case tomo.FontStyleItalic:
		return defaultfont.FaceItalic
	case tomo.FontStyleBoldItalic:
		return defaultfont.FaceBoldItalic
	default:
		return defaultfont.FaceRegular
	}
}

// Icon returns an icon from the default set corresponding to the given name.
func (Default) Icon (id tomo.Icon, size tomo.IconSize, c tomo.Case) artist.Icon {
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
func (Default) MimeIcon (data.Mime, tomo.IconSize, tomo.Case) artist.Icon {
	// TODO
	return nil
}

// Pattern returns a pattern from the default theme corresponding to the given
// pattern ID.
func (Default) Pattern (id tomo.Pattern, state tomo.State, c tomo.Case) artist.Pattern {
	offset := 0; switch {
	case state.Disabled:                 offset = 1
	case state.Pressed && state.On:      offset = 4
	case state.Focused && state.On:      offset = 6
	case state.On:                       offset = 2
	case state.Pressed:                  offset = 3
	case state.Focused:                  offset = 5
	}

	switch id {
	case tomo.PatternBackground: return patterns.Uhex(0xaaaaaaFF)
	case tomo.PatternDead:       return defaultTextures[0][offset]
	case tomo.PatternRaised:
		if c.Match("tomo", "listEntry", "") {
			return defaultTextures[10][offset]
		} else {
			return defaultTextures[1][offset]
		}
	case tomo.PatternSunken:   return defaultTextures[2][offset]
	case tomo.PatternPinboard: return defaultTextures[3][offset]
	case tomo.PatternButton:
		switch {
		case c.Match("tomo", "checkbox", ""):  
			return defaultTextures[9][offset]
		case c.Match("tomo", "piano", "flatKey"):
			return defaultTextures[11][offset]
		case c.Match("tomo", "piano", "sharpKey"):
			return defaultTextures[12][offset]
		default:
			return defaultTextures[4][offset]
		}
	case tomo.PatternInput:   return defaultTextures[5][offset]
	case tomo.PatternGutter:  return defaultTextures[6][offset]
	case tomo.PatternHandle:  return defaultTextures[7][offset]
	case tomo.PatternLine:    return defaultTextures[8][offset]
	case tomo.PatternMercury: return defaultTextures[13][offset]
	default:             return patterns.Uhex(0xFF00FFFF)
	}
}

func (Default) Color (id tomo.Color, state tomo.State, c tomo.Case) color.RGBA {
	if state.Disabled { return artist.Hex(0x444444FF) }
	
	return artist.Hex (map[tomo.Color] uint32 {
		tomo.ColorBlack:        0x000000FF,
		tomo.ColorRed:          0x880000FF,
		tomo.ColorGreen:        0x008800FF,
		tomo.ColorYellow:       0x888800FF,
		tomo.ColorBlue:         0x000088FF,
		tomo.ColorPurple:       0x880088FF,
		tomo.ColorCyan:         0x008888FF,
		tomo.ColorWhite:        0x888888FF,
		tomo.ColorBrightBlack:  0x888888FF,
		tomo.ColorBrightRed:    0xFF0000FF,
		tomo.ColorBrightGreen:  0x00FF00FF,
		tomo.ColorBrightYellow: 0xFFFF00FF,
		tomo.ColorBrightBlue:   0x0000FFFF,
		tomo.ColorBrightPurple: 0xFF00FFFF,
		tomo.ColorBrightCyan:   0x00FFFFFF,
		tomo.ColorBrightWhite:  0xFFFFFFFF,
	
		tomo.ColorForeground: 0x000000FF,
		tomo.ColorBackground: 0xAAAAAAFF,
		tomo.ColorAccent:     0x408090FF,
	} [id])
}

// Padding returns the default padding value for the given pattern.
func (Default) Padding (id tomo.Pattern, c tomo.Case) artist.Inset {
	switch id {
	case tomo.PatternRaised:
		if c.Match("tomo", "listEntry", "") {
			return artist.I(4, 8)
		} else {
			return artist.I(8)
		}
	case tomo.PatternSunken:
		if c.Match("tomo", "list", "") {
			return artist.I(4, 0, 3)
		} else if c.Match("basic", "progressBar", "") {
			return artist.I(2, 1, 1, 2)
		} else {
			return artist.I(8)
		}
	case tomo.PatternPinboard:
		if c.Match("tomo", "piano", "") {
			return artist.I(2)
		} else {
			return artist.I(8)
		}
	case tomo.PatternGutter:     return artist.I(0)
	case tomo.PatternLine:       return artist.I(1)
	case tomo.PatternMercury:    return artist.I(5)
	default:                return artist.I(8)
	}
}

// Margin returns the default margin value for the given pattern.
func (Default) Margin (id tomo.Pattern, c tomo.Case) image.Point {
	return image.Pt(8, 8)
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
