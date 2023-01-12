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
	window.SetTitle("dialog")

	container := basic.NewContainer(layouts.Dialog { true, true })
	window.Adopt(container)

	container.Adopt(basic.NewLabel("you will explode", false), true)
	cancel := basic.NewButton("Cancel")
	cancel.SetEnabled(false)
	container.Adopt(cancel, false)
	okButton := basic.NewButton("OK")
	container.Adopt(okButton, false)
	okButton.Select()
		
	window.OnClose(tomo.Stop)
	window.Show()
}
