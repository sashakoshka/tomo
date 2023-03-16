package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/layouts/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("dialog")

	container := basicElements.NewContainer(basicLayouts.Dialog { true, true })
	window.Adopt(container)

	container.Adopt(basicElements.NewLabel("you will explode", true), true)
	cancel := basicElements.NewButton("Cancel")
	cancel.SetEnabled(false)
	container.Adopt(cancel, false)
	okButton := basicElements.NewButton("OK")
	container.Adopt(okButton, false)
	okButton.Focus()
		
	window.OnClose(tomo.Stop)
	window.Show()
}
