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
	window.SetTitle("Switches")

	container := containers.NewContainer(layouts.Vertical { true, true })
	window.Adopt(container)

	container.Adopt(elements.NewSwitch("hahahah", false), false)
	container.Adopt(elements.NewSwitch("hehehehheheh", false), false)
	container.Adopt(elements.NewSwitch("you can flick da swicth", false), false)
		
	window.OnClose(tomo.Stop)
	window.Show()
}
