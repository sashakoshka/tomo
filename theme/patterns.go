package theme

import "image"
import "git.tebibyte.media/sashakoshka/tomo/artist"

// Case sepecifies what kind of element is using a pattern. It contains a
// namespace parameter and an element parameter. The element parameter does not
// necissarily need to match an element name, but if it can, it should. Both
// parameters should be written in camel case. Themes can change their styling
// based on this parameter for fine-grained control over the look and feel of
// specific elements.
type Case struct { Namespace, Element string } 
 
// C can be used as shorthand to generate a case struct as used in PatternState.
func C (namespace, element string) (c Case) {
	return Case {
		Namespace: namespace,
		Element: element,
	}
}

// PatternState lists parameters which can change the appearance of some
// patterns. For example, passing a PatternState with Selected set to true may
// result in a pattern that has a colored border within it.
type PatternState struct {
	Case

	// On should be set to true if the element that is using this pattern is
	// in some sort of "on" state, such as if a checkbox is checked or a
	// switch is toggled on. This is only necessary if the element in
	// question is capable of being toggled.
	On bool

	// Focused should be set to true if the element that is using this
	// pattern is currently focused.
	Focused bool

	// Pressed should be set to true if the element that is using this
	// pattern is being pressed down by the mouse. This is only necessary if
	// the element in question processes mouse button events.
	Pressed bool

	// Disabled should be set to true if the element that is using this
	// pattern is locked and cannot be interacted with. Disabled variations
	// of patterns are typically flattened and greyed-out.
	Disabled bool

	// Invalid should be set to true if th element that is using this
	// pattern wants to warn the user of an invalid interaction or data
	// entry. Invalid variations typically have some sort of reddish tint
	// or outline.
	Invalid bool
}

// Inset represents an inset amount for all four sides of a rectangle. The top
// side is at index zero, the right at index one, the bottom at index two, and
// the left at index three. These values may be negative.
type Inset [4]int

// Apply returns the given rectangle, shrunk on all four sides by the given
// inset. If a measurment of the inset is negative, that side will instead be
// expanded outward. If the rectangle's dimensions cannot be reduced any
// further, an empty rectangle near its center will be returned.
func (inset Inset) Apply (bigger image.Rectangle) (smaller image.Rectangle) {
	smaller = bigger
	if smaller.Dx() < inset[3] + inset[1] {
		smaller.Min.X = (smaller.Min.X + smaller.Max.X) / 2
		smaller.Max.X = smaller.Min.X
	} else {
		smaller.Min.X += inset[3]
		smaller.Max.X -= inset[1]
	}

	if smaller.Dy() < inset[0] + inset[2] {
		smaller.Min.Y = (smaller.Min.Y + smaller.Max.Y) / 2
		smaller.Max.Y = smaller.Min.Y
	} else {
		smaller.Min.Y += inset[0]
		smaller.Max.Y -= inset[2]
	}
	return
}

// Inverse returns a negated version of the inset.
func (inset Inset) Inverse () (prime Inset) {
	return Inset {
		inset[0] * -1,
		inset[1] * -1,
		inset[2] * -1,
		inset[3] * -1,
	}
}

// AccentPattern returns the accent pattern, which is usually just a solid
// color.
func AccentPattern (state PatternState) (pattern artist.Pattern, inset Inset) {
	return accentPattern, Inset { }
}

// BackgroundPattern returns the main background pattern.
func BackgroundPattern (state PatternState) (pattern artist.Pattern, inset Inset) {
	return backgroundPattern, Inset { }
}

// DeadPattern returns a pattern that can be used to mark an area or gap that
// serves no purpose, but still needs aesthetic structure.
func DeadPattern (state PatternState) (pattern artist.Pattern, inset Inset) {
	return deadPattern, Inset { }
}

// ForegroundPattern returns the color text should be.
func ForegroundPattern (state PatternState) (pattern artist.Pattern, inset Inset) {
	if state.Disabled {
		return weakForegroundPattern, Inset { }
	} else {
		return foregroundPattern, Inset { }
	}
}

// InputPattern returns a background pattern for any input field that can be
// edited by typing with the keyboard.
func InputPattern (state PatternState) (pattern artist.Pattern, inset Inset) {
	if state.Disabled {
		return disabledInputPattern, Inset { 1, 1, 1, 1 }
	} else {
		if state.Focused {
			return selectedInputPattern, Inset { 1, 1, 1, 1 }
		} else {
			return inputPattern, Inset { 1, 1, 1, 1 }
		}
	}
}

// ListPattern returns a background pattern for a list of things.
func ListPattern (state PatternState) (pattern artist.Pattern, inset Inset) {
	if state.Focused {
		pattern = focusedListPattern
		inset = Inset { 2, 1, 2, 1 }
	} else {
		pattern = listPattern
		inset = Inset { 2, 1, 1, 1 }
	}
	return
}

// ItemPattern returns a background pattern for a list item.
func ItemPattern (state PatternState) (pattern artist.Pattern, inset Inset) {
	if state.Focused {
		if state.On {
			pattern = focusedOnListEntryPattern
		} else {
			pattern = focusedListEntryPattern
		}
	} else {
		if state.On {
			pattern = onListEntryPattern
		} else {
			pattern = listEntryPattern
		}
	}
	inset = Inset { 4, 6, 4, 6 }
	return
}

// ButtonPattern returns a pattern to be displayed on buttons.
func ButtonPattern (state PatternState) (pattern artist.Pattern, inset Inset) {
	if state.Disabled {
		return disabledButtonPattern, Inset { 1, 1, 1, 1 }
	} else {
		if state.Pressed {
			if state.Focused {
				return pressedSelectedButtonPattern, Inset {
					2, 0, 0, 2 }
			} else {
				return pressedButtonPattern, Inset { 2, 0, 0, 2 }
			}
		} else {
			if state.Focused {
				return selectedButtonPattern, Inset { 1, 1, 1, 1 }
			} else {
				return buttonPattern, Inset { 1, 1, 1, 1 }
			}
		}
	}
}

// GutterPattern returns a pattern to be used to mark a track along which
// something slides.
func GutterPattern (state PatternState) (pattern artist.Pattern, inset Inset) {
	if state.Disabled {
		return disabledScrollGutterPattern, Inset { 0, 0, 0, 0 }
	} else {
		return scrollGutterPattern, Inset { 0, 0, 0, 0 }
	}
}

// HandlePattern returns a pattern to be displayed on a grab handle that slides
// along a gutter.
func HandlePattern (state PatternState) (pattern artist.Pattern, inset Inset) {
	if state.Disabled {
		return disabledScrollBarPattern, Inset { 1, 1, 1, 1 }
	} else {
		if state.Focused {
			if state.Pressed {
				return pressedSelectedScrollBarPattern, Inset { 1, 1, 1, 1 }
			} else {
				return selectedScrollBarPattern, Inset { 1, 1, 1, 1 }
			}
		} else {
			if state.Pressed {
				return pressedScrollBarPattern, Inset { 1, 1, 1, 1 }
			} else {
				return scrollBarPattern, Inset { 1, 1, 1, 1 }
			}
		}
	}
}

// SunkenPattern returns a general purpose pattern that is sunken/engraved into
// the background.
func SunkenPattern (state PatternState) (pattern artist.Pattern, inset Inset) {
	return sunkenPattern, Inset { 1, 1, 1, 1 }
}

// RaisedPattern returns a general purpose pattern that is raised up out of the
// background.
func RaisedPattern (state PatternState) (pattern artist.Pattern, inset Inset) {
	if state.Focused {
		return selectedRaisedPattern, Inset { 1, 1, 1, 1 }
	} else {
		return raisedPattern, Inset { 1, 1, 1, 1 }
	}
}

// PinboardPattern returns a textured backdrop pattern. Anything drawn within it
// should have its own background pattern.
func PinboardPattern (state PatternState) (pattern artist.Pattern, inset Inset) {
	return texturedSunkenPattern, Inset { 1, 1, 1, 1 }
}
