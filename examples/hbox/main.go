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

	container := elements.NewHBox(elements.SpaceBoth)
	window.Adopt(container)

	container.AdoptExpand(elements.NewLabelWrapped("this is sample text"))
	container.AdoptExpand(elements.NewLabelWrapped("this is sample text"))
	container.AdoptExpand(elements.NewLabelWrapped("this is sample text"))
	
	window.OnClose(tomo.Stop)
	window.Show()
}
