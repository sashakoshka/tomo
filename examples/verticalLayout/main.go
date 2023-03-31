package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/layouts"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import "git.tebibyte.media/sashakoshka/tomo/elements/testing"
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(128, 128)
	window.SetTitle("vertical stack")

	container := containers.NewContainer(layouts.Vertical { true, true })
	window.Adopt(container)

	label    := elements.NewLabel("it is a label hehe", true)
	button   := elements.NewButton("drawing pad")
	okButton := elements.NewButton("OK")
	button.OnClick (func () {
		container.DisownAll()
		container.Adopt(elements.NewLabel("Draw here:", false), false)
		container.Adopt(testing.NewMouse(), true)
		container.Adopt(okButton, false)
		okButton.Focus()
	})
	okButton.OnClick(tomo.Stop)
	
	container.Adopt(label, true)
	container.Adopt(button, false)
	container.Adopt(okButton, false)
	okButton.Focus()
	
	window.OnClose(tomo.Stop)
	window.Show()
}
