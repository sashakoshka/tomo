package textmanip

import "unicode"

// Dot represents a cursor or text selection. It has a start and end position,
// referring to where the user began and ended the selection respectively.
type Dot struct { Start, End int }

// EmptyDot returns a zero-width dot at the specified position.
func EmptyDot (position int) Dot {
	return Dot { position, position }
}

// Canon places the lesser value at the start, and the greater value at the end.
// Note that a canonized dot does not in all cases correspond directly to the
// original, because there is a semantic value to the start and end positions.
func (dot Dot) Canon () Dot {
	if dot.Start > dot.End {
		return Dot { dot.End, dot.Start }
	} else {
		return dot
	}
}

// Empty returns whether or not the 
func (dot Dot) Empty () bool {
	return dot.Start == dot.End
}

// Add shifts the dot to the right by the specified amount.
func (dot Dot) Add (delta int) Dot {
	return Dot {
		dot.Start + delta,
		dot.End   + delta,
	}
}

// Sub shifts the dot to the left by the specified amount.
func (dot Dot) Sub (delta int) Dot {
	return Dot {
		dot.Start - delta,
		dot.End   - delta,
	}
}

// Constrain constrains the dot's start and end from zero to length (inclusive).
func (dot Dot) Constrain (length int) Dot {
	if dot.Start < 0      { dot.Start = 0 }
	if dot.Start > length { dot.Start = length}
	if dot.End < 0        { dot.End = 0 }
	if dot.End > length   { dot.End = length}
	return dot
}

// Width returns how many runes the dot spans.
func (dot Dot) Width () int {
	dot = dot.Canon()
	return dot.End - dot.Start
}

// Slice returns the subset of text that the dot covers.
func (dot Dot) Slice (text []rune) []rune {
	dot = dot.Canon().Constrain(len(text))
	return text[dot.Start:dot.End]
}

// WordToLeft returns how far away to the left the next word boundary is from a
// given position.
func WordToLeft (text []rune, position int) (length int) {
	if position < 1 { return }
	if position > len(text) { position = len(text) }

	index := position - 1
	for index >= 0 && unicode.IsSpace(text[index]) {
		length ++
		index --
	}
	for index >= 0 && !unicode.IsSpace(text[index]) {
		length ++
		index --
	}
	return
}

// WordToRight returns how far away to the right the next word boundary is from
// a given position.
func WordToRight (text []rune, position int) (length int) {
	if position < 0 { return }
	if position > len(text) { position = len(text) }

	index := position
	for index < len(text) && unicode.IsSpace(text[index]) {
		length ++
		index ++
	}
	for index < len(text) && !unicode.IsSpace(text[index]) {
		length ++
		index ++
	}
	return
}

// WordAround returns a dot that surrounds the word at the specified position.
func WordAround (text []rune, position int) (around Dot) {
	return Dot {
		position - WordToLeft(text, position),
		position + WordToRight(text, position),
	}
}

// Backspace deletes the rune to the left of the dot. If word is true, it
// deletes up until the next word boundary on the left. If the dot is non-empty,
// it deletes the text inside of the dot.
func Backspace (text []rune, dot Dot, word bool) (result []rune, moved Dot) {
	dot = dot.Constrain(len(text))
	if dot.Empty() {
		distance := 1
		if word {
			distance = WordToLeft(text, dot.End)
		}
		result = append (
			result,
			text[:dot.Sub(distance).Constrain(len(text)).End]...)
		result = append(result, text[dot.End:]...)
		moved = EmptyDot(dot.Sub(distance).Start)
		return
	} else {
		return Delete(text, dot, word)
	}
}

// Delete deletes the rune to the right of the dot. If word is true, it deletes
// up until the next word boundary on the right. If the dot is non-empty, it
// deletes the text inside of the dot.
func Delete (text []rune, dot Dot, word bool) (result []rune, moved Dot) {
	dot = dot.Constrain(len(text))
	if dot.Empty() {
		distance := 1
		if word {
			distance = WordToRight(text, dot.End)
		}
		result = append(result, text[:dot.End]...)
		result = append (
			result,
			text[dot.Add(distance).Constrain(len(text)).End:]...)
		moved = dot
		return	
	} else {
		dot = dot.Canon()
		result = append(result, text[:dot.Start]...)
		result = append(result, text[dot.End:]...)
		moved = EmptyDot(dot.Start)
		return
	}
}

// Lift removes the section of text inside of the dot, and returns a copy of it.
func Lift (text []rune, dot Dot) (result []rune, moved Dot, lifted []rune) {
	dot = dot.Constrain(len(text))
	if dot.Empty() {
		moved = dot
		return
	}

	dot = dot.Canon()
	lifted = make([]rune, dot.Width())
	copy(lifted, dot.Slice(text))
	result = append(result, text[:dot.Start]...)
	result = append(result, text[dot.End:]...)
	moved = EmptyDot(dot.Start)
	return
}

// Type inserts one of more runes into the text at the dot position. If the dot
// is non-empty, it replaces the text inside of the dot with the new runes.
func Type (text []rune, dot Dot, characters ...rune) (result []rune, moved Dot) {
	dot = dot.Constrain(len(text))
	if dot.Empty() {
		result = append(result, text[:dot.End]...)
		result = append(result, characters...)
		if dot.End < len(text) {
			result = append(result, text[dot.End:]...)
		}
		moved = EmptyDot(dot.Add(len(characters)).End)
		return
	} else {
		dot = dot.Canon()
		result = append(result, text[:dot.Start]...)
		result = append(result, characters...)
		result = append(result, text[dot.End:]...)
		moved = EmptyDot(dot.Add(len(characters)).Start)
		return
	}
}

// MoveLeft moves the dot left one rune. If word is true, it moves the dot to
// the next word boundary on the left.
func MoveLeft (text []rune, dot Dot, word bool) (moved Dot) {
	dot = dot.Canon().Constrain(len(text))
	distance := 0
	if dot.Empty() {
		distance = 1
	}
	if word {
		distance = WordToLeft(text, dot.Start)
	}
	moved = EmptyDot(dot.Sub(distance).Start)
	return
}

// MoveRight moves the dot right one rune. If word is true, it moves the dot to
// the next word boundary on the right.
func MoveRight (text []rune, dot Dot, word bool) (moved Dot) {
	dot = dot.Canon().Constrain(len(text))
	distance := 0
	if dot.Empty() {
		distance = 1
	}
	if word {
		distance = WordToRight(text, dot.End)
	}
	moved = EmptyDot(dot.Add(distance).End)
	return
}

// SelectLeft moves the end of the dot left one rune. If word is true, it moves
// the end of the dot to the next word boundary on the left.
func SelectLeft (text []rune, dot Dot, word bool) (moved Dot) {
	dot = dot.Constrain(len(text))
	distance := 1
	if word {
		distance = WordToLeft(text, dot.End)
	}
	dot.End -= distance
	return dot
}

// SelectRight moves the end of the dot right one rune. If word is true, it
// moves the end of the dot to the next word boundary on the right.
func SelectRight (text []rune, dot Dot, word bool) (moved Dot) {
	dot = dot.Constrain(len(text))
	distance := 1
	if word {
		distance = WordToRight(text, dot.End)
	}
	dot.End += distance
	return dot
}
