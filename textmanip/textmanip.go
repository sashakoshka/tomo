package textmanip

import "unicode"

func WordToLeft (text []rune, cursor int) (length int) {
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

func WordToRight (text []rune, cursor int) (length int) {
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

func Backspace (text []rune, cursor int, word bool) (result []rune, moved int) {
	if cursor < 1         { return text, cursor }
	if cursor > len(text) { cursor = len(text) }

	moved = 1
	if word {
		moved = WordToLeft(text, cursor)
	}
	result = append(result, text[:cursor - moved]...)
	result = append(result, text[cursor:]...)
	moved = cursor - moved
	return
}

func Delete (text []rune, cursor int, word bool) (result []rune, moved int) {
	if cursor < 0         { return text, cursor }
	if cursor > len(text) { cursor = len(text) }

	moved = 1
	if word {
		moved = WordToRight(text, cursor)
	}
	result = append(result, text[:cursor]...)
	result = append(result, text[cursor + moved:]...)
	moved = cursor
	return
}

func Type (text []rune, cursor int, character rune) (result []rune, moved int) {
	if cursor < 0         { cursor = 0 }
	if cursor > len(text) { cursor = len(text) }
	result = append(result, text[:cursor]...)
	result = append(result, character)
	if cursor < len(text) {
		result = append(result, text[cursor:]...)
	}
	moved = cursor + 1
	return
}

func MoveLeft (text []rune, cursor int, word bool) (moved int) {
	if cursor < 1         { return cursor }
	if cursor > len(text) { cursor = len(text) }

	moved = 1
	if word {
		moved = WordToLeft(text, cursor)
	}
	moved = cursor - moved
	return
}

func MoveRight (text []rune, cursor int, word bool) (moved int) {
	if cursor < 0         { return cursor }
	if cursor > len(text) { cursor = len(text) }

	moved = 1
	if word {
		moved = WordToRight(text, cursor)
	}
	moved = cursor + moved
	return
}
