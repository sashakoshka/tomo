package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import "git.tebibyte.media/sashakoshka/tomo/elements/testing"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(tomo.Bounds(0, 0, 128, 128))
	window.SetTitle("vertical stack")

	container := elements.NewVBox(true, true)

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
	
	window.Adopt(container)
	window.OnClose(tomo.Stop)
	window.Show()
}
