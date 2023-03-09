package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements/testing"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"
import "git.tebibyte.media/sashakoshka/ezprof/ez"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(480, 360)
	window.Adopt(testing.NewArtist())
	window.OnClose(tomo.Stop)
	window.Show()
	ez.Prof()
}
