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
	window.SetTitle("Switches")

	container := basicElements.NewContainer(basicLayouts.Vertical { true, true })
	window.Adopt(container)

	container.Adopt(basicElements.NewSwitch("hahahah", false), false)
	container.Adopt(basicElements.NewSwitch("hehehehheheh", false), false)
	container.Adopt(basicElements.NewSwitch("you can flick da swicth", false), false)
		
	window.OnClose(tomo.Stop)
	window.Show()
}
