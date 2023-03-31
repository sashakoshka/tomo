package tomo

import "time"

// Config can return global configuration parameters.
type Config interface {
	// ScrollVelocity returns how many pixels should be scrolled every time
	// a scroll button is pressed.
	ScrollVelocity () int

	// DoubleClickDelay returns the maximum delay between two clicks for
	// them to be registered as a double click.
	DoubleClickDelay () time.Duration
}
