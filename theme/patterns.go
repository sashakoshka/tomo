package theme

import "image"
import "git.tebibyte.media/sashakoshka/tomo/artist"

// PatternState lists parameters which can change the appearance of some
// patterns. For example, passing a PatternState with Selected set to true may
// result in a pattern that has a colored border within it.
type PatternState struct {
	// On should be set to true if the element that is using this pattern is
	// in some sort of "on" state, such as if a checkbox is checked or a
	// switch is toggled on. This is only necessary if the element in
	// question is capable of being toggled.
	On bool

	// Selected should be set to true if the element that is using this
	// pattern is currently selected.
	Selected bool

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
	if smaller.Dx() < inset[3] + inset[0] {
		smaller.Min.X = (smaller.Min.X + smaller.Max.X) / 2
		smaller.Max.X = smaller.Min.X
	} else {
		smaller.Min.X += inset[3]
		smaller.Max.X -= inset[1]
	}

	if smaller.Dy() < inset[1] + inset[2] {
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
		if state.Selected {
			return selectedInputPattern, Inset { 1, 1, 1, 1 }
		} else {
			return inputPattern, Inset { 1, 1, 1, 1 }
		}
	}
}

// TODO: for list and item patterns, have all that bizarre padding/2 information
// in the insets.

// ListPattern returns a background pattern for a list of things.
func ListPattern (state PatternState) (pattern artist.Pattern, inset Inset) {
	if state.Selected {
		return selectedListPattern, Inset { }
	} else {
		return listPattern, Inset { }
	}
}

// ItemPattern returns a background pattern for a list item.
func ItemPattern (state PatternState) (pattern artist.Pattern, inset Inset) {
	if state.On {
		return selectedListEntryPattern, Inset { 1, 1, 1, 1 }
	} else {
		return listEntryPattern, Inset { 1, 1, 1, 1 }
	}
}

// ButtonPattern returns a pattern to be displayed on buttons.
func ButtonPattern (state PatternState) (pattern artist.Pattern, inset Inset) {
	if state.Disabled {
		return disabledButtonPattern, Inset { 1, 1, 1, 1 }
	} else {
		if state.Pressed {
			if state.Selected {
				return pressedSelectedButtonPattern, Inset {
					2, 0, 0, 2 }
			} else {
				return pressedButtonPattern, Inset { 2, 0, 0, 2 }
			}
		} else {
			if state.Selected {
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
		if state.Selected {
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
	if state.Selected {
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
