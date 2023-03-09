package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements/testing"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"
import _ "git.tebibyte.media/sashakoshka/ezprof/hook"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(480, 360)
	window.SetTitle("Draw Test")
	window.Adopt(testing.NewArtist())
	window.OnClose(tomo.Stop)
}
