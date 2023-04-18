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

	container := elements.NewVBox (
		elements.SpaceBoth,
		elements.NewLabel("This is at the top"),
		elements.NewLine(),
		elements.NewLabel("This is in the middle"))
	container.AdoptExpand(elements.NewSpacer())
	container.Adopt(elements.NewLabel("This is at the bottom"))
	
	window.Adopt(container)
	window.OnClose(tomo.Stop)
	window.Show()
}
