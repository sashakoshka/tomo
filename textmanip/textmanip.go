package textmanip

import "unicode"

type Dot struct { Start, End int }

func EmptyDot (position int) Dot {
	return Dot { position, position }
}

func (dot Dot) Canon () Dot {
	if dot.Start > dot.End {
		return Dot { dot.End, dot.Start }
	} else {
		return dot
	}
}

func (dot Dot) Empty () bool {
	return dot.Start == dot.End
}

func (dot Dot) Add (delta int) Dot {
	return Dot {
		dot.Start + delta,
		dot.End   + delta,
	}
}

func (dot Dot) Sub (delta int) Dot {
	return Dot {
		dot.Start - delta,
		dot.End   - delta,
	}
}

func WordToLeft (text []rune, dot Dot) (length int) {
	cursor := dot.End
	if cursor < 1 { return }
	if cursor > len(text) { cursor = len(text) }

	index := cursor - 1
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

func WordToRight (text []rune, dot Dot) (length int) {
	cursor := dot.End
	if cursor < 0 { return }
	if cursor > len(text) { cursor = len(text) }

	index := cursor
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

func Backspace (text []rune, dot Dot, word bool) (result []rune, moved Dot) {
	if dot.Empty() {
		cursor := dot.End
		if cursor < 1         { return text, dot }
		if cursor > len(text) { cursor = len(text) }

		distance := 1
		if word {
			distance = WordToLeft(text, dot)
		}
		result = append(result, text[:cursor - distance]...)
		result = append(result, text[cursor:]...)
		moved = EmptyDot(cursor - distance)
	} else {
		return Delete(text, dot, word)
	}

	return
}

func Delete (text []rune, dot Dot, word bool) (result []rune, moved Dot) {
	if dot.Empty() {
		cursor := dot.End
		if cursor < 0         { return text, dot }
		if cursor > len(text) { cursor = len(text) }

		distance := 1
		if word {
			distance = WordToRight(text, dot)
		}
		result = append(result, text[:cursor]...)
		result = append(result, text[cursor + distance:]...)
		moved = dot
		return	
	} else {
		result = append(result, text[:dot.Start]...)
		result = append(result, text[dot.End:]...)
		moved = EmptyDot(dot.Start)
		return
	}
}

func Type (text []rune, dot Dot, character rune) (result []rune, moved Dot) {
	if dot.Empty() {
		cursor := dot.End
		if cursor < 0         { cursor = 0 }
		if cursor > len(text) { cursor = len(text) }
		result = append(result, text[:cursor]...)
		result = append(result, character)
		if cursor < len(text) {
			result = append(result, text[cursor:]...)
		}
		moved = EmptyDot(cursor + 1)
		return
	} else {
		result = append(result, text[:dot.Start]...)
		result = append(result, character)
		result = append(result, text[dot.End:]...)
		moved = EmptyDot(dot.Start)
		return
	}
}

func MoveLeft (text []rune, dot Dot, word bool) (moved Dot) {
	cursor := dot.Start

	if cursor < 1         { return EmptyDot(cursor) }
	if cursor > len(text) { cursor = len(text) }

	distance := 1
	if word {
		distance = WordToLeft(text, dot)
	}
	moved = EmptyDot(cursor - distance)
	return
}

func MoveRight (text []rune, dot Dot, word bool) (moved Dot) {
	cursor := dot.End
	
	if cursor < 0         { return EmptyDot(cursor) }
	if cursor > len(text) { cursor = len(text) }

	distance := 1
	if word {
		distance = WordToRight(text, dot)
	}
	moved = EmptyDot(cursor + distance)
	return
}
