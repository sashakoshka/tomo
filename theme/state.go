package theme

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
