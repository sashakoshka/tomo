package config

import "git.tebibyte.media/sashakoshka/tomo"

// Default specifies default configuration values.
type Default struct { }


// ScrollVelocity returns the default scroll velocity value.
func (Default) ScrollVelocity () int {
	return 16
}

// Wrapped wraps a configuration and uses Default if it is nil.
type Wrapped struct {
	tomo.Config
}

// ScrollVelocity returns how many pixels should be scrolled every time
// a scroll button is pressed.
func (wrapped Wrapped) ScrollVelocity () int {
	return wrapped.ensure().ScrollVelocity()
}

func (wrapped Wrapped) ensure () (real tomo.Config) {
	real = wrapped.Config
	if real == nil { real = Default { } }
	return
}
