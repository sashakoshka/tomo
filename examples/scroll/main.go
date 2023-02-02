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
	window.SetTitle("Scroll")
	container := basicElements.NewContainer(basicLayouts.Vertical { true, true })
	window.Adopt(container)

	container.Adopt(basicElements.NewLabel("look at this non sense", false), false)

	textBox := basicElements.NewTextBox("", "sample text sample text")
	scrollContainer := basicElements.NewScrollContainer(true, false)
	scrollContainer.Adopt(textBox)
	container.Adopt(scrollContainer, true)
	
	window.OnClose(tomo.Stop)
	window.Show()
}
