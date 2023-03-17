package textdraw

import "unicode"
import "golang.org/x/image/font"
import "golang.org/x/image/math/fixed"

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

// RuneLayout contains layout information for a single rune relative to its
// word.
type RuneLayout struct {
	X     fixed.Int26_6
	Width fixed.Int26_6
	Rune  rune
}

// WordLayout contains layout information for a single word relative to its
// line.
type WordLayout struct {
	X          fixed.Int26_6
	Width      fixed.Int26_6
	SpaceAfter fixed.Int26_6
	Runes      []RuneLayout
}

// DoWord consumes exactly one word from the given string, and produces a word
// layout according to the given font. It returns the remaining text as well.
func DoWord (text []rune, face font.Face) (word WordLayout, remaining []rune) {
	remaining     = text
	gettingSpace := false
	x := fixed.Int26_6(0)
	lastRune     := rune(-1)
	for _, char := range text {
		// if we run into a line break, we must break out immediately
		// because it is not DoWord's job to handle that.
		if char == '\n' { break }
	
		// if we suddenly run into spaces, and then run into a word
		// again, we must break out immediately.
		if unicode.IsSpace(char) {
			gettingSpace = true
		} else if gettingSpace {
			break
		}

		// apply kerning
		if lastRune >= 0 { x += face.Kern(lastRune, char) }
		lastRune = char
		
		// consume and process the rune
		remaining = remaining[1:]
		_, advance, ok := face.GlyphBounds(char)
		if !ok { continue }
		word.Runes = append (word.Runes, RuneLayout {
			X:     x,
			Width: advance,
			Rune:  char,
		})

		// advance
		if gettingSpace {
			word.SpaceAfter += advance
		} else {
			word.Width += advance
		}
		x += advance
	}
	return
}

// LastRune returns the last rune in the word.
func (word WordLayout) LastRune () rune {
	if word.Runes == nil {
		return -1
	} else {
		return word.Runes[len(word.Runes) - 1].Rune
	}
}

// FirstRune returns the last rune in the word.
func (word WordLayout) FirstRune () rune {
	if word.Runes == nil {
		return -1
	} else {
		return word.Runes[0].Rune
	}
}

// LineLayout contains layout information for a single line.
type LineLayout struct {
	Y            fixed.Int26_6
	Width        fixed.Int26_6
	ContentWidth fixed.Int26_6
	SpaceAfter   fixed.Int26_6
	Words        []WordLayout
	BreakAfter bool
}

// DoLine consumes exactly one line from the given string, and produces a line
// layout according to the given font. It returns the remaining text as well. If
// maxWidth is greater than zero, this function will stop processing words once
// the limit is crossed. The word which would have crossed over the limit will
// not be processed.
func DoLine (text []rune, face font.Face, maxWidth fixed.Int26_6) (line LineLayout, remaining []rune) {
	remaining    = text
	x           := fixed.Int26_6(0)
	lastWord    := WordLayout { }
	isFirstWord := true
	for {
		// process one word
		word, remainingFromWord := DoWord(remaining, face)
		word.X = x
		x += word.Width

		// if we have gone over the maximum width, stop processing
		// words (if maxWidth is even specified)
		if !isFirstWord && maxWidth > 0 && x > maxWidth {
			break
		}

		x += word.SpaceAfter
		remaining = remainingFromWord

		// if the word actually has contents, add it
		if word.Runes != nil {
			lastWord   = word
			line.Words = append(line.Words, word)
		}

		// if we have hit the end of the line, stop processing words
		if len(remaining) == 0 { break }
		if remaining[0] == '\n' {
			line.BreakAfter = true
			remaining = remaining[1:]
			break
		}

		isFirstWord = false
	}

	// set the width of the line's content.
	line.ContentWidth = lastWord.X + lastWord.Width

	// set the line's width. this is subject to be overridden by the
	// TypeSetter to match the longest line.
	if maxWidth > 0 {
		line.Width = maxWidth
	} else {
		line.Width      = line.ContentWidth
		line.SpaceAfter = lastWord.SpaceAfter
	}
	return
}

// Align aligns the text in the line according to the specified alignment
// method.
func (line *LineLayout) Align (align Align) {
	if len(line.Words) == 0 { return }

	if align == AlignJustify {
		line.justify()
		return
	}

	leftOffset := -line.Words[0].X

	if align == AlignCenter {
		leftOffset += (line.Width - line.ContentWidth) / 2
	} else if align == AlignRight {
		leftOffset += line.Width - line.ContentWidth
	}

	for index := range line.Words {
		line.Words[index].X += leftOffset
	}
}

func (line *LineLayout) justify () {
	if len(line.Words) < 2 {
		line.Align(AlignLeft)
		return
	}

	// We are going to be moving the words, so we can't take SpaceAfter into
	// account.
	trueContentWidth := fixed.Int26_6(0)
	for _, word := range line.Words {
		trueContentWidth += word.Width
	}

	spaceCount   := fixed.Int26_6(len(line.Words) - 1)
	spacePerWord := (line.Width - trueContentWidth) / spaceCount
	x := fixed.Int26_6(0)
	for index, word := range line.Words {
		line.Words[index].X = x
		x += spacePerWord + word.Width
	}
}
