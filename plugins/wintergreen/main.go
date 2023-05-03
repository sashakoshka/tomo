// Plugin wintergreen provides a calm, bluish green theme.
package main

import "tomo"
import "tomo/plugins/wintergreen/wintergreen"

func Expects () tomo.Version {
	return tomo.Version { 0, 0, 0 }
}

func Name () string {
	return "Wintergreen"
}

func Description () string {
	return "A calm, bluish green theme."
}

func NewTheme () (tomo.Theme) {
	return wintergreen.Theme { }
}
