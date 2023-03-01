package config

// Config can return global configuration parameters.
type Config interface {
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
