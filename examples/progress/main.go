package main

import "time"
import "tomo"
import "tomo/nasin"
import "tomo/popups"
import "tomo/elements"

func main () {
	nasin.Run(Application { })
}

type Application struct { }

func (Application) Init () error {
	window, err := nasin.NewWindow(tomo.Bounds(0, 0, 0, 0))
	if err != nil { return err }
	window.SetTitle("Approaching")
	container := elements.NewVBox(elements.SpaceBoth)
	window.Adopt(container)

	container.AdoptExpand(elements.NewLabel("Rapidly approaching your location..."))
	bar := elements.NewProgressBar(0)
	container.Adopt(bar)
	button := elements.NewButton("Stop")
	button.SetEnabled(false)
	container.Adopt(button)
	
	window.OnClose(nasin.Stop)
	window.Show()
	go fill(window, bar)
	return nil
}

func fill (window tomo.Window, bar *elements.ProgressBar) {
	for progress := 0.0; progress < 1.0; progress += 0.01 {
		time.Sleep(time.Second / 24)
		nasin.Do (func () {
			bar.SetProgress(progress)
		})
	}
	nasin.Do (func () {
		popups.NewDialog (
			popups.DialogKindInfo,
			window,
			"I am here",
			"Don't look outside your window.")
	})
}
