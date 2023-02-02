package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/layouts/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("Spaced Out")

	container := basicElements.NewContainer(basicLayouts.Vertical { true, true })
	window.Adopt(container)

	container.Adopt (basicElements.NewLabel("This is at the top", false), false)
	container.Adopt (basicElements.NewSpacer(true), false)
	container.Adopt (basicElements.NewLabel("This is in the middle", false), false)
	container.Adopt (basicElements.NewSpacer(false), true)
	container.Adopt (basicElements.NewLabel("This is at the bottom", false), false)
	
	window.OnClose(tomo.Stop)
	window.Show()
}
