package config

// Config can return global configuration parameters.
type Config interface {
	// Padding returns the amount of internal padding elements should have.
	// An element's inner content (such as text) should be inset by this
	// amount, in addition to the inset returned by the pattern of its
	// background.
	Padding () int
	
	// Margin returns how much space should be put in between elements.
	Margin () int

	// HandleWidth returns how large grab handles should typically be. This
	// is important for accessibility reasons.
	HandleWidth () int

	// ScrollVelocity returns how many pixels should be scrolled every time
	// a scroll button is pressed.
	ScrollVelocity () int

	// ThemePath returns the directory path to the theme.
	ThemePath () string
}

// Default specifies default configuration values.
type Default struct { }

// Padding returns the default padding value.
func (Default) Padding () int {
	return 7
}

// Margin returns the default margin value.
func (Default) Margin () int {
	return 8
}

// HandleWidth returns the default handle width value.
func (Default) HandleWidth () int {
	return 16
}

// ScrollVelocity returns the default scroll velocity value.
func (Default) ScrollVelocity () int {
	return 16
}

// ThemePath returns the default theme path.
func (Default) ThemePath () (string) {
	return ""
}

// Wrapped wraps a configuration and uses Default if it is nil.
type Wrapped struct {
	Config
}

// Padding returns the amount of internal padding elements should have.
// An element's inner content (such as text) should be inset by this
// amount, in addition to the inset returned by the pattern of its
// background.
func (wrapped Wrapped) Padding () int {
	return wrapped.ensure().Padding()
}

// Margin returns how much space should be put in between elements.
func (wrapped Wrapped) Margin () int {
	return wrapped.ensure().Margin()
}

// HandleWidth returns how large grab handles should typically be. This
// is important for accessibility reasons.
func (wrapped Wrapped) HandleWidth () int {
	return wrapped.ensure().HandleWidth()
}

// ScrollVelocity returns how many pixels should be scrolled every time
// a scroll button is pressed.
func (wrapped Wrapped) ScrollVelocity () int {
	return wrapped.ensure().ScrollVelocity()
}

// ThemePath returns the directory path to the theme.
func (wrapped Wrapped) ThemePath () string {
	return wrapped.ensure().ThemePath()
}

func (wrapped Wrapped) ensure () (real Config) {
	real = wrapped.Config
	if real == nil { real = Default { } }
	return
}
