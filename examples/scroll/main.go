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
	window.SetTitle("Scroll")
	container := basic.NewContainer(layouts.Vertical { true, true })
	window.Adopt(container)

	container.Adopt(basic.NewLabel("look at this non sense", false), false)

	textBox := basic.NewTextBox("", "sample text sample text")
	scrollContainer := basic.NewScrollContainer(true, true)
	scrollContainer.Adopt(textBox)
	container.Adopt(scrollContainer, true)
	
	window.OnClose(tomo.Stop)
	window.Show()
}
