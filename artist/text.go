package artist

// import "fmt"
import "image"
import "unicode"
import "image/draw"
import "golang.org/x/image/font"
import "golang.org/x/image/math/fixed"
import "git.tebibyte.media/sashakoshka/tomo/canvas"

type characterLayout struct {
	x         int
	character rune
}

type wordLayout struct {
	position    image.Point
	width       int
	spaceAfter  int
	breaksAfter int
	text        []characterLayout
	whitespace  []characterLayout
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
func (drawer *TextDrawer) SetText (runes []rune) {
	// if drawer.runes == runes { return }
	drawer.runes = runes
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
	destination canvas.Canvas,
	source      Pattern,
	offset      image.Point,
) (
	updatedRegion image.Rectangle,
) {
	wrappedSource := WrappedPattern {
		Pattern: source,
		Width:  0,
		Height: 0, // TODO: choose a better width and height
	}

	if !drawer.layoutClean { drawer.recalculate() }
	// TODO: reimplement a version of draw mask that takes in a pattern and
	// only draws to a tomo.Canvas.
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
			wrappedSource, image.Point { },
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

// ReccomendedHeightFor returns the reccomended max height if the text were to
// have its maximum width set to the given width. This does not alter the
// drawer's state.
func (drawer *TextDrawer) ReccomendedHeightFor (width int) (height int) {
	if drawer.face == nil { return }
	if !drawer.layoutClean { drawer.recalculate() }
	metrics := drawer.face.Metrics()
	dot := fixed.Point26_6 { 0, metrics.Height }
	for _, word := range drawer.layout {
		if word.width + dot.X.Round() > width {
			dot.Y += metrics.Height
			dot.X = 0
		}
		dot.X += fixed.I(word.width + word.spaceAfter)
		if word.breaksAfter > 0 {
			dot.Y += fixed.I(word.breaksAfter).Mul(metrics.Height)
			dot.X = 0
		}
	}

	return dot.Y.Round()
}

// PositionOf returns the position of the character at the specified index
// relative to the baseline.
func (drawer *TextDrawer) PositionOf (index int) (position image.Point) {
	if !drawer.layoutClean { drawer.recalculate() }
	index ++
	for _, word := range drawer.layout {
		position = word.position
		for _, character := range word.text {
			index --
			position.X = word.position.X + character.x
			if index < 1 { return }
		}
		for _, character := range word.whitespace {
			index --
			position.X = word.position.X + character.x
			if index < 1 { return }
		}
	}
	return
}

// Length returns the amount of runes in the drawer's text.
func (drawer *TextDrawer) Length () (length int) {
	return len(drawer.runes)
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
	horizontalExtent  := 0
	currentCharacterX := fixed.Int26_6(0)

	previousCharacter := rune(-1)
	for index < len(drawer.runes) {
		word := wordLayout { }
		word.position.X = dot.X.Round()
		word.position.Y = dot.Y.Round()

		// process a word
		currentCharacterX  = 0
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

		// process whitespace, going onto a new line if there is a
		// newline character
		spaceWidth := fixed.Int26_6(0)
		for index < len(drawer.runes) && unicode.IsSpace(drawer.runes[index]) {
			character := drawer.runes[index]
			_, advance, ok := drawer.face.GlyphBounds(character)
			index ++
			if !ok { continue }
			word.whitespace = append(word.whitespace, characterLayout {
				x: currentCharacterX.Round(),
				character: character,
			})
			spaceWidth        += advance
			currentCharacterX += advance
			
			if character == '\n' {
				dot.Y += metrics.Height
				dot.X = 0
				word.breaksAfter ++
				break
			} else {
				dot.X += advance
				if previousCharacter >= 0 {
					dot.X += drawer.face.Kern (
						previousCharacter,
						character)
				}
			}
			previousCharacter = character
		}
		word.spaceAfter = spaceWidth.Round()

		// add the word to the layout
		drawer.layout = append(drawer.layout, word)

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

	// add a little null to the last character
	if len(drawer.layout) > 0 {
		lastWord := &drawer.layout[len(drawer.layout) - 1]
		lastWord.whitespace = append (
			lastWord.whitespace,
			characterLayout {
				x: currentCharacterX.Round(),
			})
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
