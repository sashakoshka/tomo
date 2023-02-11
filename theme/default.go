package theme

import "image"
import "golang.org/x/image/font"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/defaultfont"

// Default is the default theme.
type Default struct { }

// FontFace returns the default font face.
func (Default) FontFace (style FontStyle, size FontSize, c Case) font.Face {
	switch style {
	case FontStyleBold:
		return defaultfont.FaceBold
	case FontStyleItalic:
		return defaultfont.FaceItalic
	case FontStyleBoldItalic:
		return defaultfont.FaceBoldItalic
	default:
		return defaultfont.FaceRegular
	}
}

// Icon returns an icon from the default set corresponding to the given name.
func (Default) Icon (string, Case) artist.Pattern {
	// TODO
	return uhex(0)
}

// Pattern returns a pattern from the default theme corresponding to the given
// pattern ID.
func (Default) Pattern (
	pattern Pattern,
	c Case,
	state PatternState,
) artist.Pattern {
	switch pattern {
	case PatternAccent:
		return accentPattern
	case PatternBackground:
		return backgroundPattern
	case PatternForeground:
		if state.Disabled || c == C("basic", "spacer") {
			return weakForegroundPattern
		} else {
			return foregroundPattern
		}
	case PatternDead:
		return deadPattern
	case PatternRaised:
		if c == C("basic", "listEntry") {
			if state.Focused {
				if state.On {
					return focusedOnListEntryPattern
				} else {
					return focusedListEntryPattern
				}
			} else {
				if state.On {
					return onListEntryPattern
				} else {
					return listEntryPattern
				}
			}
		} else {
			if state.Focused {
				return selectedRaisedPattern
			} else {
				return raisedPattern
			}
		}
	case PatternSunken:
		if c == C("basic", "list") {
			if state.Focused {
				return focusedListPattern
			} else {
				return listPattern
			}
		} else if c == C("basic", "textBox") {
			if state.Disabled {
				return disabledInputPattern
			} else {
				if state.Focused {
					return selectedInputPattern
				} else {
					return inputPattern
				}
			}
		} else {
			if state.Focused {
				return focusedSunkenPattern
			} else {
				return sunkenPattern
			}
		}
	case PatternPinboard:
		if state.Focused {
			return focusedTexturedSunkenPattern
		} else {
			return texturedSunkenPattern
		}
	case PatternButton:
		if state.Disabled {
			return disabledButtonPattern
		} else {
			if c == C("fun", "sharpKey") {
				if state.Pressed {
					return pressedDarkButtonPattern
				} else {
					return darkButtonPattern
				}
			} else if c == C("fun", "flatKey") {
				if state.Pressed {
					return pressedButtonPattern
				} else {
					return buttonPattern
				}	
			} else {
				if state.Pressed || state.On && c == C("basic", "checkbox") {
					if state.Focused {
						return pressedSelectedButtonPattern
					} else {
						return pressedButtonPattern
					}
				} else {
					if state.Focused {
						return selectedButtonPattern
					} else {
						return buttonPattern
					}
				}
			}
		}
	case PatternInput:
		if state.Disabled {
			return disabledInputPattern
		} else {
			if state.Focused {
				return selectedInputPattern
			} else {
				return inputPattern
			}
		}
	case PatternGutter:
		if state.Disabled {
			return disabledScrollGutterPattern
		} else {
			return scrollGutterPattern
		}
	case PatternHandle:
		if state.Disabled {
			return disabledScrollBarPattern
		} else {
			if state.Focused {
				if state.Pressed {
					return pressedSelectedScrollBarPattern
				} else {
					return selectedScrollBarPattern
				}
			} else {
				if state.Pressed {
					return pressedScrollBarPattern
				} else {
					return scrollBarPattern
				}
			}
		}
	default:
		return uhex(0)
	}
}

// Inset returns the default inset value for the given pattern.
func (Default) Inset (pattern Pattern, c Case) Inset {
	switch pattern {
	case PatternRaised:
		if c == C("basic", "listEntry") {
			return Inset { 4, 6, 4, 6 }
		} else {
			return Inset { 2, 2, 2, 2 }
		}
		
	case PatternSunken:
		if c == C("basic", "list") {
			return Inset { 2, 1, 2, 1 }
		} else if c == C("basic", "progressBar") {
			return Inset { 2, 1, 1, 2 }
		} else {
			return Inset { 2, 2, 2, 2 }
		}

	case PatternPinboard:
		return Inset { 2, 2, 2, 2 }
	
	case PatternInput, PatternButton, PatternHandle:
		return Inset { 2, 2, 2, 2}

	default: return Inset { }
	}
}

// Sink returns the default sink vector for the given pattern.
func (Default) Sink (pattern Pattern, c Case) image.Point {
	return image.Point { 1, 1 }
}
