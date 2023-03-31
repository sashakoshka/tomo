package tomo

// Config can return global configuration parameters.
type Config interface {
	// ScrollVelocity returns how many pixels should be scrolled every time
	// a scroll button is pressed.
	ScrollVelocity () int
}
