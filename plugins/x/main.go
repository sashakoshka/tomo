// Plugin x provides the X11 backend as a plugin.
package main

import "tomo"
import "tomo/plugins/x/x"

func Expects () tomo.Version {
	return tomo.Version { 0, 0, 0 }
}

func Name () string {
	return "X"
}

func Description () string {
	return "Provides an X11 backend."
}

func NewBackend () (tomo.Backend, error) {
	return x.NewBackend()
}
