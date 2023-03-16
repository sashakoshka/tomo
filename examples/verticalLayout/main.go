package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/layouts/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/testing"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("vertical stack")

	container := basicElements.NewContainer(basicLayouts.Vertical { true, true })
	window.Adopt(container)

	label    := basicElements.NewLabel("it is a label hehe", true)
	button   := basicElements.NewButton("drawing pad")
	okButton := basicElements.NewButton("OK")
	button.OnClick (func () {
		container.DisownAll()
		container.Adopt(basicElements.NewLabel("Draw here:", false), false)
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
