package tomo

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
