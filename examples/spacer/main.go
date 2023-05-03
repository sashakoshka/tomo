package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/nasin"
import "git.tebibyte.media/sashakoshka/tomo/elements"

func main () {
	nasin.Run(Application { })
}

type Application struct { }

func (Application) Init () error {
	window, err := nasin.NewWindow(tomo.Bounds(0, 0, 0, 0))
	if err != nil { return err }
	window.SetTitle("Spaced Out")

	container := elements.NewVBox (
		elements.SpaceBoth,
		elements.NewLabel("This is at the top"),
		elements.NewLine(),
		elements.NewLabel("This is in the middle"))
	container.AdoptExpand(elements.NewSpacer())
	container.Adopt(elements.NewLabel("This is at the bottom"))
	
	window.Adopt(container)
	window.OnClose(nasin.Stop)
	window.Show()
	return nil
}
