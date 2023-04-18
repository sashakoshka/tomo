package main

import "fmt"
import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(tomo.Bounds(200, 200, 256, 256))
	window.SetTitle("Main")

	container := elements.NewVBox (
		elements.SpaceBoth,
		elements.NewLabel("Main window"))
	window.Adopt(container)
		
	window.OnClose(tomo.Stop)
	window.Show()

	createPanel(window, 0, tomo.Bounds(-64, 20,  0, 0))
	createPanel(window, 1, tomo.Bounds(200, 20,  0, 0))
	createPanel(window, 2, tomo.Bounds(-64, 180, 0, 0))
	createPanel(window, 3, tomo.Bounds(200, 180, 0, 0))
}

func createPanel (parent tomo.MainWindow, id int, bounds image.Rectangle) {
	window, _ := parent.NewPanel(bounds)
	title := fmt.Sprint("Panel #", id)
	window.SetTitle(title)
	container := elements.NewVBox (
		elements.SpaceBoth,
		elements.NewLabel(title))
	window.Adopt(container)
	window.Show()
}
