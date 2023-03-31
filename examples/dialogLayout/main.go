package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/layouts"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("dialog")

	container := containers.NewContainer(layouts.Dialog { true, true })
	window.Adopt(container)

	container.Adopt(elements.NewLabel("you will explode", false), true)
	cancel := elements.NewButton("Cancel")
	cancel.SetEnabled(false)
	container.Adopt(cancel, false)
	okButton := elements.NewButton("OK")
	container.Adopt(okButton, false)
	okButton.Focus()
		
	window.OnClose(tomo.Stop)
	window.Show()
}
