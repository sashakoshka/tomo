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

	layout := layouts.NewVertical(true, true)
	window.Adopt(layout)

	layout.Adopt(basic.NewLabel("it is a label hehe"))
	layout.Adopt(basic.NewButton("yeah"), false)
	layout.Adopt(button := basic.NewButton("wow"), false)
	
	window.OnClose(tomo.Stop)
	window.Show()
}
