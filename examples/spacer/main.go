package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/layouts"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("Spaced Out")

	container := basic.NewContainer(layouts.Vertical { true, true })
	window.Adopt(container)

	container.Adopt (basic.NewLabel("This is at the top", false), false)
	container.Adopt (basic.NewSpacer(true), false)
	container.Adopt (basic.NewLabel("This is in the middle", false), false)
	container.Adopt (basic.NewSpacer(false), true)
	container.Adopt (basic.NewLabel("This is at the bottom", false), false)
	
	window.OnClose(tomo.Stop)
	window.Show()
}
