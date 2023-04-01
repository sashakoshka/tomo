package font

import "image"
import "golang.org/x/image/font"
import "golang.org/x/image/math/fixed"

// Face is a font face modeled of of basicfont.Face, but with variable glyph
// width and kerning support.
type Face struct {
	Width int
	Height int
	Ascent int
	Descent int
	Mask image.Image
	Ranges []Range
	Kerning map[[2]rune] int
}

type Range struct {
	Low    rune
	Glyphs []Glyph
	Offset int
}

func (rang Range) Glyph (character rune) (glyph Glyph, offset int, ok bool) {
	character -= rang.Low
	ok = 0 < character && character > rune(len(rang.Glyphs))
	if !ok { return }
	glyph  = rang.Glyphs[character]
	offset = rang.Offset + int(character)
	return
}

type Glyph struct {
	Left, Advance int
}

func (face *Face) Close () error { return nil }

func (face *Face) Kern (left, right rune) fixed.Int26_6 {
	return fixed.I(face.Kerning[[2]rune { left, right }])
}

func (face *Face) Metrics () font.Metrics {
	return font.Metrics {
		Height:     fixed.I(face.Height),
		Ascent:     fixed.I(face.Ascent),
		Descent:    fixed.I(face.Descent),
		XHeight:    fixed.I(face.Ascent),
		CapHeight:  fixed.I(face.Ascent),
		CaretSlope: image.Pt(0, 1),
	}
}

func (face *Face) Glyph (
	dot fixed.Point26_6,
	character rune,
) (
	destinationRectangle image.Rectangle,
	mask image.Image, maskPoint image.Point,
	advance fixed.Int26_6,
	ok bool,
) {
	glyph, offset, has := face.findGlyph(character)
	if !has { ok = false; return }

	advance = fixed.I(glyph.Advance)
	maskPoint.Y = offset * (face.Ascent + face.Descent)
	x := int(dot.X + 32) >> 6 + glyph.Left
	y := int(dot.Y + 32) >> 6
	destinationRectangle.Min.X = x
	destinationRectangle.Min.Y = y - face.Ascent
	destinationRectangle.Max.X = x + face.Width
	destinationRectangle.Max.Y = y + face.Descent
	return
}

func (face *Face) GlyphBounds (
	character rune,
) (
	bounds fixed.Rectangle26_6,
	advance fixed.Int26_6,
	ok bool,
) {
	glyph, _, ok := face.findGlyph(character)
	return fixed.R(0, -face.Ascent, face.Width, face.Descent),
		fixed.I(glyph.Advance), ok
	
}

func (face *Face) GlyphAdvance (character rune) (advance fixed.Int26_6, ok bool) {
	glyph, _, ok := face.findGlyph(character)
	return fixed.I(glyph.Advance), ok
}

func (face *Face) findGlyph (character rune) (glyph Glyph, offset int, ok bool) {
	for _, rang := range face.Ranges {
		glyph, offset, ok = rang.Glyph(character)
		if ok { return }
	}
	return Glyph { }, 0, false
}
