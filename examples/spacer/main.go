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
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("Spaced Out")

	container := containers.NewContainer(layouts.Vertical { true, true })
	window.Adopt(container)

	container.Adopt (elements.NewLabel("This is at the top", false), false)
	container.Adopt (elements.NewSpacer(true), false)
	container.Adopt (elements.NewLabel("This is in the middle", false), false)
	container.Adopt (elements.NewSpacer(false), true)
	container.Adopt (elements.NewLabel("This is at the bottom", false), false)
	
	window.OnClose(tomo.Stop)
	window.Show()
}
