package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(tomo.Bounds(0, 0, 360, 0))
	window.SetTitle("horizontal stack")

	container := elements.NewHBox(true, true)
	window.Adopt(container)

	container.Adopt(elements.NewLabel("this is sample text", true), true)
	container.Adopt(elements.NewLabel("this is sample text", true), true)
	container.Adopt(elements.NewLabel("this is sample text", true), true)
	
	window.OnClose(tomo.Stop)
	window.Show()
}
