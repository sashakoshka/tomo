package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/layouts"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("vertical stack")

	container := basic.NewContainer(layouts.Horizontal { true, true })
	window.Adopt(container)

	container.Adopt(basic.NewTest(), true)
	container.Adopt(basic.NewLabel("<- left\nright ->", false), false)
	container.Adopt(basic.NewTest(), true)
	
	window.OnClose(tomo.Stop)
	window.Show()
}
