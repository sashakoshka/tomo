package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(128, 128)
	window.SetTitle("hellorld!")
	window.Adopt(basic.NewTest())
	window.OnClose(tomo.Stop)
	window.Show()
}
