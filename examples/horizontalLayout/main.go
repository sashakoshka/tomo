package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/layouts"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(360, 2)
	window.SetTitle("horizontal stack")

	container := containers.NewContainer(layouts.Horizontal { true, true })
	window.Adopt(container)

	container.Adopt(elements.NewLabel("this is sample text", true), true)
	container.Adopt(elements.NewLabel("this is sample text", true), true)
	container.Adopt(elements.NewLabel("this is sample text", true), true)
	
	window.OnClose(tomo.Stop)
	window.Show()
}
