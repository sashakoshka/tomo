package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/layouts/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/fun"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("Piano")
	container := basicElements.NewContainer(basicLayouts.Vertical { true, true })
	window.Adopt(container)

	label := basicElements.NewLabel("Play a song!", false)
	container.Adopt(label, false)
	piano := fun.NewPiano(3, 5)
	container.Adopt(piano, true)
	
	window.OnClose(tomo.Stop)
	window.Show()
}
