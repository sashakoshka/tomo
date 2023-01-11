package artist

// import "fmt"
import "image"
import "unicode"
import "image/draw"
import "golang.org/x/image/font"
import "golang.org/x/image/math/fixed"
import "git.tebibyte.media/sashakoshka/tomo"

type characterLayout struct {
	x         int
	character rune
}

type wordLayout struct {
	position   image.Point
	width      int
	text       []characterLayout
}

// Align specifies a text alignment method.
type Align int

const (
	// AlignLeft aligns the start of each line to the beginning point
	// of each dot.
	AlignLeft Align = iota
	AlignRight
	AlignCenter
	AlignJustify
)

// TextDrawer is a struct that is capable of efficient rendering of wrapped
// text, and calculating text bounds. It avoids doing redundant work
// automatically.
type TextDrawer struct {
	text   string
	runes  []rune
	face   font.Face
	width  int
	height int
	align  Align
	wrap   bool
	cut    bool

	layout       []wordLayout
	layoutClean  bool
	layoutBounds image.Rectangle
}

// SetText sets the text of the text drawer.
func (drawer *TextDrawer) SetText (text string) {
	if drawer.text == text { return }
	drawer.text  = text
	drawer.runes = []rune(text)
	drawer.layoutClean = false
}

// SetFace sets the font face of the text drawer.
func (drawer *TextDrawer) SetFace (face font.Face) {
	if drawer.face == face { return }
	drawer.face = face
	drawer.layoutClean = false
}

// SetMaxWidth sets a maximum width for the text drawer, and recalculates the
// layout if needed. If zero is given, there will be no width limit and the text
// will not wrap.
func (drawer *TextDrawer) SetMaxWidth (width int) {
	if drawer.width == width { return }
	drawer.width = width
	drawer.wrap = width != 0
	drawer.layoutClean = false
}

// SetMaxHeight sets a maximum height for the text drawer. Lines that are
// entirely below this height will not be drawn, and lines that are on the cusp
// of this maximum height will be clipped at the point that they cross it.
func (drawer *TextDrawer) SetMaxHeight (height int) {
	if drawer.height == height { return }
	drawer.height = height
	drawer.cut = height != 0
	drawer.layoutClean = false
}

// SetAlignment specifies how the drawer should align its text. For this to have
// an effect, a maximum width must have been set.
func (drawer *TextDrawer) SetAlignment (align Align) {
	if drawer.align == align { return }
	drawer.align = align
	drawer.layoutClean = false
}

// Draw draws the drawer's text onto the specified canvas at the given offset.
func (drawer *TextDrawer) Draw (
	destination tomo.Canvas,
	source      tomo.Image,
	offset      image.Point,
) (
	updatedRegion image.Rectangle,
) {
	if !drawer.layoutClean { drawer.recalculate() }
	for _, word := range drawer.layout {
	for _, character := range word.text {
		destinationRectangle,
		mask, maskPoint, _, ok := drawer.face.Glyph (
			fixed.P (
				offset.X + word.position.X + character.x,
				offset.Y + word.position.Y),
			character.character)
		if !ok { continue }

		// FIXME: clip destination rectangle if we are on the cusp of
		// the maximum height.

		draw.DrawMask (
			destination,
			destinationRectangle,
			source, image.Point { },
			mask, maskPoint,
			draw.Over)

		updatedRegion = updatedRegion.Union(destinationRectangle)
	}}
	return
}

// LayoutBounds returns a semantic bounding box for text to be used to determine
// an offset for drawing. If a maximum width or height has been set, those will
// be used as the width and height of the bounds respectively. The origin point
// (0, 0) of the returned bounds will be equivalent to the baseline at the start
// of the first line. As such, the minimum of the bounds will be negative.
func (drawer *TextDrawer) LayoutBounds () (bounds image.Rectangle) {
	if !drawer.layoutClean { drawer.recalculate() }
	bounds = drawer.layoutBounds
	return
}

// Em returns the width of an emspace.
func (drawer *TextDrawer) Em () (width fixed.Int26_6) {
	if drawer.face == nil { return }
	width, _ = drawer.face.GlyphAdvance('M')
	return
}

// LineHeight returns the height of one line.
func (drawer *TextDrawer) LineHeight () (height fixed.Int26_6) {
	if drawer.face == nil { return }
	metrics := drawer.face.Metrics()
	height = metrics.Height
	return
}

func (drawer *TextDrawer) recalculate () {
	drawer.layoutClean = true
	drawer.layout = nil
	drawer.layoutBounds = image.Rectangle { }
	if drawer.runes == nil { return }
	if drawer.face  == nil { return }

	metrics := drawer.face.Metrics()
	dot := fixed.Point26_6 { 0, 0 }
	index := 0
	horizontalExtent := 0

	previousCharacter := rune(-1)
	for index < len(drawer.runes) {
		word := wordLayout { }
		word.position.X = dot.X.Round()
		word.position.Y = dot.Y.Round()

		// process a word
		currentCharacterX := fixed.Int26_6(0)
		wordWidth         := fixed.Int26_6(0)
		for index < len(drawer.runes) && !unicode.IsSpace(drawer.runes[index]) {
			character := drawer.runes[index]
			_, advance, ok := drawer.face.GlyphBounds(character)
			index ++
			if !ok { continue }

			word.text = append(word.text, characterLayout {
				x: currentCharacterX.Round(),
				character: character,
			})
			
			dot.X             += advance
			wordWidth         += advance
			currentCharacterX += advance
			if dot.X.Round () > horizontalExtent {
				horizontalExtent = dot.X.Round()
			}
			if previousCharacter >= 0 {
				dot.X += drawer.face.Kern (
					previousCharacter,
					character)
			}
			previousCharacter = character
		}
		word.width = wordWidth.Round()

		// detect if the word that was just processed goes out of
		// bounds, and if it does, wrap it
		if drawer.wrap &&
			word.width + word.position.X > drawer.width &&
			word.position.X > 0 {
			
			word.position.Y += metrics.Height.Round()
			word.position.X = 0
			dot.Y += metrics.Height
			dot.X = wordWidth
		}

		// add the word to the layout
		drawer.layout = append(drawer.layout, word)

		// skip over whitespace, going onto a new line if there is a
		// newline character
		for index < len(drawer.runes) && unicode.IsSpace(drawer.runes[index]) {
			character := drawer.runes[index]
			if character == '\n' {
				dot.Y += metrics.Height
				dot.X = 0
				previousCharacter = character
				index ++
			} else {
				_, advance, ok := drawer.face.GlyphBounds(character)
				index ++
				if !ok { continue }
				
				dot.X += advance
				if previousCharacter >= 0 {
					dot.X += drawer.face.Kern (
						previousCharacter,
						character)
				}
				previousCharacter = character
			}
		}

		// if there is a set maximum height, and we have crossed it,
		// stop processing more words. and remove any words that have
		// also crossed the line.
		if
			drawer.cut &&
			(dot.Y - metrics.Ascent - metrics.Descent).Round() >
			drawer.height {

			for
				index := len(drawer.layout) - 1;
				index >= 0; index -- {

				if drawer.layout[index].position.Y < dot.Y.Round() {
					break
				}
				drawer.layout = drawer.layout[:index]
			}
			break
		}
	}

	if drawer.wrap {
		drawer.layoutBounds.Max.X = drawer.width
	} else {
		drawer.layoutBounds.Max.X = horizontalExtent
	}

	if drawer.cut {
		drawer.layoutBounds.Min.Y = 0 - metrics.Ascent.Round()
		drawer.layoutBounds.Max.Y = drawer.height - metrics.Ascent.Round()
	} else {
		drawer.layoutBounds.Min.Y = 0 - metrics.Ascent.Round()
		drawer.layoutBounds.Max.Y = dot.Y.Round() + metrics.Descent.Round()
	}
	
	// TODO:
	// for each line, calculate the bounds as if the words are left aligned,
	// and then at the end of the process go through each line and re-align
	// everything. this will make the process far simpler.
}
