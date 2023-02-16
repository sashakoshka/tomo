package textdraw

import "image"
import "golang.org/x/image/font"
import "golang.org/x/image/math/fixed"

// TypeSetter manages several lines of text, and can perform layout operations
// on them. It automatically avoids performing redundant work. It has no
// constructor and its zero value can be used safely.
type TypeSetter struct {
	lines []LineLayout
	text []rune
	
	layoutClean bool
	alignClean  bool
	
	align Align
	face  font.Face
	maxWidth  int
	maxHeight int

	layoutBounds      image.Rectangle
	layoutBoundsSpace image.Rectangle
}

func (setter *TypeSetter) needLayout () {
	if setter.layoutClean { return }
	setter.layoutClean = true
	setter.alignClean  = false

	// we need to have a font and some text to do anything
	setter.lines = nil
	setter.layoutBounds      = image.Rectangle { }
	setter.layoutBoundsSpace = image.Rectangle { }
	if len(setter.text) == 0 { return }
	if setter.face  == nil { return }

	horizontalExtent      := fixed.Int26_6(0)
	horizontalExtentSpace := fixed.Int26_6(0)

	lastLine  := LineLayout { }
	metrics   := setter.face.Metrics()
	remaining := setter.text
	y         := fixed.Int26_6(0)
	maxY      := fixed.I(setter.maxHeight) + metrics.Height
	for len(remaining) > 0 && (y < maxY || setter.maxHeight == 0) {
		// process one line
		line, remainingFromLine := DoLine (
			remaining, setter.face, fixed.I(setter.maxWidth))
		remaining = remainingFromLine

		// add the line
		line.Y = y
		y += metrics.Height
		if line.Width > horizontalExtent {
			horizontalExtent = line.Width
		}
		lineWidthSpace := line.Width + line.SpaceAfter
		if lineWidthSpace > horizontalExtentSpace {
			horizontalExtentSpace = lineWidthSpace
		}
		setter.lines = append(setter.lines, line)
		lastLine = line
	}

	// add a null onto the end because the very end of the text should have
	// a valid layout position
	lastWord := &lastLine.Words[len(lastLine.Words) - 1]
	lastWord.Runes = append (lastWord.Runes, RuneLayout {
		X: lastWord.Width + lastWord.SpaceAfter,
		Rune: 0,
	})

	// set all line widths to horizontalExtent if we don't have a specified
	// maximum width
	if setter.maxWidth == 0 {
		for index := range setter.lines {
			setter.lines[index].Width = horizontalExtent
		}
		setter.layoutBounds.Max.X      = horizontalExtent.Round()
		setter.layoutBoundsSpace.Max.X = horizontalExtentSpace.Round()
	} else {
		setter.layoutBounds.Max.X      = setter.maxWidth
		setter.layoutBoundsSpace.Max.X = setter.maxWidth
	}

	y -= metrics.Height
	if setter.maxHeight == 0 {
		setter.layoutBounds.Min.Y = -metrics.Ascent.Round()
		setter.layoutBounds.Max.Y =
			y.Round() +
			metrics.Descent.Round()
	} else {
		setter.layoutBounds.Min.Y = -metrics.Ascent.Round()
		setter.layoutBounds.Max.Y =
			setter.maxHeight -
			metrics.Ascent.Round()
	}
	setter.layoutBoundsSpace.Min.Y = setter.layoutBounds.Min.Y
	setter.layoutBoundsSpace.Max.Y = setter.layoutBounds.Max.Y
}

func (setter *TypeSetter) needAlignedLayout () {
	if setter.alignClean && setter.layoutClean { return }
	setter.needLayout()
	setter.alignClean = true

	for index := range setter.lines {
		setter.lines[index].Align(setter.align)
	}
}

// SetAlign sets the alignment method of the typesetter.
func (setter *TypeSetter) SetAlign (align Align) {
	if setter.align == align { return }
	setter.alignClean = false
	setter.align = align
}

// SetText sets the text content of the typesetter.
func (setter *TypeSetter) SetText (text []rune) {
	setter.layoutClean = false
	setter.alignClean  = false
	setter.text = text
}

// SetFace sets the font face of the typesetter.
func (setter *TypeSetter) SetFace (face font.Face) {
	if setter.face == face { return }
	setter.layoutClean = false
	setter.alignClean  = false
	setter.face = face
}

// SetMaxWidth sets the maximum width of the typesetter. If the maximum width
// is greater than zero, the text will wrap to that width. If the maximum width
// is zero, the text will not wrap and instead extend as far as it needs to.
func (setter *TypeSetter) SetMaxWidth (width int) {
	if setter.maxWidth == width { return }
	setter.layoutClean = false
	setter.alignClean  = false
	setter.maxWidth = width
}

// SetMaxHeight sets the maximum height of the typesetter. If the maximum height
// is greater than zero, no lines will be laid out past that point. If the
// maximum height is zero, the text's maximum height will not be constrained.
func (setter *TypeSetter) SetMaxHeight (heignt int) {
	if setter.maxHeight == heignt { return }
	setter.layoutClean = false
	setter.alignClean  = false
	setter.maxHeight = heignt
}

// Em returns the width of one emspace according to the typesetter's font, which
// is the width of the capital letter 'M'.
func (setter *TypeSetter) Em () (width fixed.Int26_6) {
	if setter.face == nil { return 0 }
	width, _ = setter.face.GlyphAdvance('M')
	return
}

// LineHeight returns the height of one line according to the typesetter's font.
func (setter *TypeSetter) LineHeight () fixed.Int26_6 {
	if setter.face == nil { return 0 }
	return setter.face.Metrics().Height
}

// MaxWidth returns the maximum width of the typesetter as set by SetMaxWidth.
func (setter *TypeSetter) MaxWidth () int {
	return setter.maxWidth
}

// MaxHeight returns the maximum height of the typesetter as set by
// SetMaxHeight.
func (setter *TypeSetter) MaxHeight () int {
	return setter.maxHeight
}

// Face returns the TypeSetter's font face as set by SetFace.
func (setter *TypeSetter) Face () font.Face {
	return setter.face
}

// Length returns the amount of runes in the typesetter.
func (setter *TypeSetter) Length () int {
	return len(setter.text)
}

// RuneIterator is a function that can iterate accross a typesetter's runes.
type RuneIterator func (
	index    int,
	char     rune,
	position fixed.Point26_6,
) (
	keepGoing bool,
)

// For calls the specified iterator for every rune in the typesetter. If the
// iterator returns false, the loop will immediately stop.
func (setter *TypeSetter) For (iterator RuneIterator) {
	setter.needAlignedLayout()

	index := 0
	for _, line := range setter.lines {
		for _, word := range line.Words {
		for _, char := range word.Runes {
			keepGoing := iterator (index, char.Rune, fixed.Point26_6 {
				X: word.X + char.X,
				Y: line.Y,
			})
			if !keepGoing { return }
			index ++
		}}
		if line.BreakAfter { index ++ }
	}
}

// AtPosition returns the index of the rune at the specified position.
func (setter *TypeSetter) AtPosition (position fixed.Point26_6) (index int) {
	println("XXX", position.Y.Round())
	setter.needAlignedLayout()
	
	if setter.lines == nil { return }
	if setter.face  == nil { return }

	// find the first line who's bottom bound is greater than position.Y. if
	// we haven't found it, then dont set the line variable (defaults to the
	// last line)
	metrics := setter.face.Metrics()
	line := setter.lines[len(setter.lines) - 1]
	lineSize := 0
	for _, curLine := range setter.lines {
		for _, curWord := range curLine.Words {
			lineSize += len(curWord.Runes)
		}
		if curLine.BreakAfter { lineSize ++ }
		index += lineSize
		
		if curLine.Y + metrics.Descent > position.Y {
			line = curLine
			break
		}
	}
	index -= lineSize

	if line.Words == nil { return }

	// find the first rune who's right bound is greater than position.X.
	for _, curWord := range line.Words {
		for _, curChar := range curWord.Runes {
			x := curWord.X + curChar.X + curChar.Width
			println(index, x.Round(), position.X.Round())
			if x > position.X { goto foundRune }
			index ++
		}
	}
	foundRune:
	return
}

// PositionAt returns the position of the rune at the specified index.
func (setter *TypeSetter) PositionAt (index int) (position fixed.Point26_6) {
	setter.needAlignedLayout()
	
	setter.For (func (i int, r rune, p fixed.Point26_6) bool {
		position = p
		return i < index
	})
	return
}

// LayoutBounds returns the semantic bounding box of the text. The origin point
// (0, 0) of the rectangle corresponds to the origin of the first line's
// baseline.
func (setter *TypeSetter) LayoutBounds () (image.Rectangle) {
	setter.needLayout()
	return setter.layoutBounds
	
}

// LayoutBoundsSpace is like LayoutBounds, but it also takes into account the
// trailing whitespace at the end of each line (if it exists).
func (setter *TypeSetter) LayoutBoundsSpace () (image.Rectangle) {
	setter.needLayout()
	return setter.layoutBoundsSpace
}

// ReccomendedHeightFor returns the reccomended max height if the text were to
// have its maximum width set to the given width. This does not alter the
// typesetter's state.
func (setter *TypeSetter) ReccomendedHeightFor (width int) (height int) {
	setter.needLayout()
	
	if setter.lines == nil { return }
	if setter.face  == nil { return }

	metrics := setter.face.Metrics()
	dot := fixed.Point26_6 { 0, metrics.Height }
	firstWord := true
	for _, line := range setter.lines {
		for _, word := range line.Words {
			if word.Width + dot.X > fixed.I(width) && !firstWord {
				dot.Y += metrics.Height
				dot.X = 0
				firstWord = true
			}
			dot.X += word.Width + word.SpaceAfter
			firstWord = false
		}
		if line.BreakAfter {
			dot.Y += metrics.Height
			dot.X = 0
			firstWord = true
		}
	}

	return dot.Y.Round()
}
