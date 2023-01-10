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

	container := basic.NewContainer(layouts.Vertical { true, true })
	window.Adopt(container)

	label  := basic.NewLabel("it is a label hehe")
	button := basic.NewButton("press me")
	button.OnClick (func () {
		label.SetText (
			"woah, this button changes the label text! since the " +
			"size of this text box has changed, the window " +
			"should expand (unless you resized it already).")
	})
	
	container.Adopt(label, true)
	container.Adopt(basic.NewButton("yeah"), false)
	container.Adopt(button, false)
	
	window.OnClose(tomo.Stop)
	window.Show()
}
