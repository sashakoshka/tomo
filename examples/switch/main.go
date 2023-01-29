package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/layouts"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("Switches")

	container := basic.NewContainer(layouts.Vertical { true, true })
	window.Adopt(container)

	container.Adopt(basic.NewSwitch("hahahah", false), false)
	container.Adopt(basic.NewSwitch("hehehehheheh", false), false)
	container.Adopt(basic.NewSwitch("you can flick da swicth", false), false)
		
	window.OnClose(tomo.Stop)
	window.Show()
}
