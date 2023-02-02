package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/layouts/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(360, 2)
	window.SetTitle("horizontal stack")

	container := basicElements.NewContainer(basicLayouts.Horizontal { true, true })
	window.Adopt(container)

	container.Adopt(basicElements.NewLabel("this is sample text", true), true)
	container.Adopt(basicElements.NewLabel("this is sample text", true), true)
	container.Adopt(basicElements.NewLabel("this is sample text", true), true)
	
	window.OnClose(tomo.Stop)
	window.Show()
}
