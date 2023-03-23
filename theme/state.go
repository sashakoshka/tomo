package theme

// Case sepecifies what kind of element is using a pattern. It contains a
// namespace parameter, an element parameter, and an optional component trail.
// All parameter values should be written in camel case. Themes can change their
// styling based on the case for fine-grained control over the look and feel of
// specific elements.
type Case struct {
	// Namespace refers to the package that the element comes from. This is
	// so different element packages can have elements with the same name
	// while still allowing themes to differentiate between them.
	Namespace string

	// Element refers to the name of the element. This should (generally) be
	// the type name of the element. For example: Button, Input, Container,
	// etc.
	Element string

	// Component specifies the specific part of the element that is being
	// referred to. This parameter is entirely optional.
	Component string
}
 
// C can be used as shorthand to generate a case struct. The component parameter
// may be left out of this argument list for brevity. Arguments passed after
// component will be ignored.
func C (namespace, element string, component ...string) Case {
	if component == nil { component = []string { "" } }
	return Case {
		Namespace: namespace,
		Element:   element,
		Component: component[0],
	}
}

// Match determines if a case matches the specified parameters. A blank string
// will act as a wildcard.
func (c Case) Match (namespace, element, component string) bool {
	if namespace == "" { namespace = c.Namespace }
	if element   == "" { element   = c.Element   }
	if component == "" { component = c.Component }

	return  namespace == c.Namespace &&
		element   == c.Element   &&
		component == c.Component
}

// State lists parameters which can change the appearance of some patterns and
// colors. For example, passing a State with Selected set to true may result in
// a pattern that has a colored border within it.
type State struct {
	// On should be set to true if the element that is using this pattern is
	// in some sort of selected or "on" state, such as if a checkbox is
	// checked, a file is selected, or a switch is toggled on. This is only
	// necessary if the element in question is capable of being toggled or
	// selected.
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

// FontStyle specifies stylistic alterations to a font face.
type FontStyle int; const (
	FontStyleRegular    FontStyle = 0
	FontStyleBold       FontStyle = 1
	FontStyleItalic     FontStyle = 2
	FontStyleBoldItalic FontStyle = 1 | 2
)

// FontSize specifies the general size of a font face in a semantic way.
type FontSize int; const (
	// FontSizeNormal is the default font size that should be used for most
	// things.
	FontSizeNormal FontSize = iota

	// FontSizeLarge is a larger font size suitable for things like section
	// headings.
	FontSizeLarge

	// FontSizeHuge is a very large font size suitable for things like
	// titles, wizard step names, digital clocks, etc.
	FontSizeHuge

	// FontSizeSmall is a smaller font size. Try not to use this unless it
	// makes a lot of sense to do so, because it can negatively impact
	// accessibility. It is useful for things like copyright notices at the
	// bottom of some window that the average user doesn't actually care
	// about.
	FontSizeSmall
)
