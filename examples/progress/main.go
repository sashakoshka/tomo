package main

import "time"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/popups"
import "git.tebibyte.media/sashakoshka/tomo/layouts"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("Approaching")
	container := basic.NewContainer(layouts.Vertical { true, true })
	window.Adopt(container)

	container.Adopt (basic.NewLabel (
		"Rapidly approaching your location...", false), false)
	bar := basic.NewProgressBar(0)
	container.Adopt(bar, false)
	button := basic.NewButton("Stop")
	button.SetEnabled(false)
	container.Adopt(button, false)
	
	window.OnClose(tomo.Stop)
	window.Show()
	go fill(bar)
}

func fill (bar *basic.ProgressBar) {
	for progress := 0.0; progress < 1.0; progress += 0.01 {
		time.Sleep(time.Second / 24)
		tomo.Do (func () {
			bar.SetProgress(progress)
		})
	}
	tomo.Do (func () {
		popups.NewDialog (
			popups.DialogKindInfo,
			"I am here",
			"Don't look outside your window.")
	})
}
