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

	container := elements.NewVBox(elements.SpaceBoth)

	label    := elements.NewLabelWrapped("it is a label hehe")
	button   := elements.NewButton("drawing pad")
	okButton := elements.NewButton("OK")
	button.OnClick (func () {
		container.DisownAll()
		container.Adopt(elements.NewLabel("Draw here (not really):"))
		container.AdoptExpand(testing.NewMouse())
		container.Adopt(okButton)
		okButton.Focus()
	})
	okButton.OnClick(tomo.Stop)

	container.AdoptExpand(label)
	container.Adopt(button, okButton)
	window.Adopt(container)
	
	okButton.Focus()
	window.OnClose(tomo.Stop)
	window.Show()
}
