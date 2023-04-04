package main

import "fmt"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/layouts"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("Table")

	container := containers.NewContainer(layouts.Vertical { true, true })
	table := containers.NewTableContainer(7, 7, true, true)
	scroller := containers.NewScrollContainer(true, true)

	index := 0
	for row := 0; row < 7; row ++ {
	for column := 0; column < 7; column ++ {
		if index % 2 == 0 {
			label := elements.NewLabel (
				fmt.Sprintf("%d, %d", row, column),
				false)
			table.Set(row, column, label)
		}
		index ++
	}}
	table.Set(2, 1, elements.NewButton("Oh hi mars!"))

	statusLabel := elements.NewLabel("Selected: none", false)
	table.Collapse(128, 128)
	table.OnSelect (func () {
		column, row := table.Selected()
		statusLabel.SetText (
			fmt.Sprintf("Selected: %d, %d",
			column, row))
	})

	scroller.Adopt(table)
	container.Adopt(scroller, true)
	container.Adopt(statusLabel, false)
	window.Adopt(container)
	
	window.OnClose(tomo.Stop)
	window.Show()
}
