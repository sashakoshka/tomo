package main

import "fmt"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import "git.tebibyte.media/sashakoshka/tomo/layouts/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("Main")

	container := containers.NewContainer(basicLayouts.Vertical { true, true })
	container.Adopt(basicElements.NewLabel("Main window", false), true)
	window.Adopt(container)
		
	window.OnClose(tomo.Stop)
	window.Show()

	createPanel(window, 0)
	// createPanel(window, 1)
	// createPanel(window, 2)
	// createPanel(window, 3)
}

func createPanel (parent elements.MainWindow, id int) {
	window, _ := parent.NewPanel(2, 2)
	title := fmt.Sprint("Panel #", id)
	window.SetTitle(title)
	container := containers.NewContainer(basicLayouts.Vertical { true, true })
	container.Adopt(basicElements.NewLabel(title, false), true)
	window.Adopt(container)
	window.Show()
}
