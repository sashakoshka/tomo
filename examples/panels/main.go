package main

import "fmt"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/layouts"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(256, 256)
	window.SetTitle("Main")

	container := containers.NewContainer(layouts.Vertical { true, true })
	container.Adopt(elements.NewLabel("Main window", false), true)
	window.Adopt(container)
		
	window.OnClose(tomo.Stop)
	window.Show()

	createPanel(window, 0)
	createPanel(window, 1)
	createPanel(window, 2)
	createPanel(window, 3)
}

func createPanel (parent tomo.MainWindow, id int) {
	window, _ := parent.NewPanel(2, 2)
	title := fmt.Sprint("Panel #", id)
	window.SetTitle(title)
	container := containers.NewContainer(layouts.Vertical { true, true })
	container.Adopt(elements.NewLabel(title, false), true)
	window.Adopt(container)
	window.Show()
}
