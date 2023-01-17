package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/layouts"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(256, 2)
	window.SetTitle("horizontal stack")

	container := basic.NewContainer(layouts.Horizontal { true, true })
	window.Adopt(container)

	container.Adopt(basic.NewLabel("this is sample text", true), true)
	container.Adopt(basic.NewLabel("this is sample text", true), true)
	container.Adopt(basic.NewLabel("this is sample text", true), true)
	
	window.OnClose(tomo.Stop)
	window.Show()
}
