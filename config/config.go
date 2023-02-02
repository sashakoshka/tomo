package config

// Padding returns the amount of internal padding elements should have. An
// element's inner content (such as text) should be inset by this amount,
// in addition to the inset returned by the pattern of its background. When
// using the aforementioned inset values to calculate the element's minimum size
// or the position and alignment of its content, all parameters in the
// PatternState should be unset except for Case.
func Padding () int {
	return 7
}

// Margin returns how much space should be put in between elements.
func Margin () int {
	return 8
}

// HandleWidth returns how large grab handles should typically be. This is
// important for accessibility reasons.
func HandleWidth () int {
	return 16
}
