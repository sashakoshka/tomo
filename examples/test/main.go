package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements/testing"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(128, 128)
	window.SetTitle("hellorld!")
	window.Adopt(testing.NewMouse())
	window.OnClose(tomo.Stop)
	window.Show()
}
