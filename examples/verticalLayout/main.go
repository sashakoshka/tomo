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

	label    := basic.NewLabel("it is a label hehe", false)
	button   := basic.NewButton("drawing pad")
	okButton := basic.NewButton("OK")
	button.OnClick (func () {
		container.DisownAll()
		container.Adopt(basic.NewLabel("Draw here:", false), false)
		container.Adopt(basic.NewTest(), true)
		container.Adopt(okButton, false)
		okButton.Select()
	})
	okButton.OnClick(tomo.Stop)
	
	container.Adopt(label, true)
	container.Adopt(button, false)
	container.Adopt(okButton, false)
	okButton.Select()
	
	window.OnClose(tomo.Stop)
	window.Show()
}
