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

func (dot Dot) Constrain (length int) Dot {
	if dot.Start < 0      { dot.Start = 0 }
	if dot.Start > length { dot.Start = length}
	if dot.End < 0        { dot.End = 0 }
	if dot.End > length   { dot.End = length}
	return dot
}

func (dot Dot) Width () int {
	dot = dot.Canon()
	return dot.End - dot.Start
}

func (dot Dot) Slice (text []rune) []rune {
	return text[dot.Start:dot.End]
}

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

func WordAround (text []rune, position int) (around Dot) {
	return Dot {
		WordToLeft(text, position),
		WordToRight(text, position),
	}
}

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

func Type (text []rune, dot Dot, character rune) (result []rune, moved Dot) {
	dot = dot.Constrain(len(text))
	if dot.Empty() {
		result = append(result, text[:dot.End]...)
		result = append(result, character)
		if dot.End < len(text) {
			result = append(result, text[dot.End:]...)
		}
		moved = EmptyDot(dot.Add(1).End)
		return
	} else {
		dot = dot.Canon()
		result = append(result, text[:dot.Start]...)
		result = append(result, character)
		result = append(result, text[dot.End:]...)
		moved = EmptyDot(dot.Add(1).Start)
		return
	}
}

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

func SelectLeft (text []rune, dot Dot, word bool) (moved Dot) {
	dot = dot.Constrain(len(text))
	distance := 1
	if word {
		distance = WordToLeft(text, dot.End)
	}
	dot.End -= distance
	return dot
}

func SelectRight (text []rune, dot Dot, word bool) (moved Dot) {
	dot = dot.Constrain(len(text))
	distance := 1
	if word {
		distance = WordToRight(text, dot.End)
	}
	dot.End += distance
	return dot
}
