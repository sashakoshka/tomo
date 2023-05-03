package config

import "time"
import "tomo"

// Default specifies default configuration values.
type Default struct { }


// ScrollVelocity returns the default scroll velocity value.
func (Default) ScrollVelocity () int {
	return 16
}

// DoubleClickDelay returns the default double click delay.
func (Default) DoubleClickDelay () time.Duration {
	return time.Second / 2
}

// Wrapped wraps a configuration and uses Default if it is nil.
type Wrapped struct {
	tomo.Config
}

// ScrollVelocity returns how many pixels should be scrolled every time a scroll
// button is pressed.
func (wrapped Wrapped) ScrollVelocity () int {
	return wrapped.ensure().ScrollVelocity()
}

// DoubleClickDelay returns the maximum delay between two clicks for them to be
// registered as a double click.
func (wrapped Wrapped) DoubleClickDelay () time.Duration {
	return wrapped.ensure().DoubleClickDelay()
}

func (wrapped Wrapped) ensure () (real tomo.Config) {
	real = wrapped.Config
	if real == nil { real = Default { } }
	return
}
