package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(tomo.Bounds(0, 0, 0, 0))
	window.SetTitle("Spaced Out")

	container := elements.NewVBox(true, true)
	window.Adopt(container)

	container.Adopt (elements.NewLabel("This is at the top", false), false)
	container.Adopt (elements.NewSpacer(true), false)
	container.Adopt (elements.NewLabel("This is in the middle", false), false)
	container.Adopt (elements.NewSpacer(false), true)
	container.Adopt (elements.NewLabel("This is at the bottom", false), false)
	
	window.OnClose(tomo.Stop)
	window.Show()
}
