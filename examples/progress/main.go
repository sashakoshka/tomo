package main

import "time"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/popups"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(tomo.Bounds(0, 0, 0, 0))
	window.SetTitle("Approaching")
	container := elements.NewVBox(elements.SpaceBoth)
	window.Adopt(container)

	container.AdoptExpand(elements.NewLabel("Rapidly approaching your location..."))
	bar := elements.NewProgressBar(0)
	container.Adopt(bar)
	button := elements.NewButton("Stop")
	button.SetEnabled(false)
	container.Adopt(button)
	
	window.OnClose(tomo.Stop)
	window.Show()
	go fill(window, bar)
}

func fill (window tomo.Window, bar *elements.ProgressBar) {
	for progress := 0.0; progress < 1.0; progress += 0.01 {
		time.Sleep(time.Second / 24)
		tomo.Do (func () {
			bar.SetProgress(progress)
		})
	}
	tomo.Do (func () {
		popups.NewDialog (
			popups.DialogKindInfo,
			window,
			"I am here",
			"Don't look outside your window.")
	})
}
