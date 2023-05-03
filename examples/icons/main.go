package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/nasin"
import "git.tebibyte.media/sashakoshka/tomo/elements"

func main () {
	nasin.Run(Application { })
}

type Application struct { }

func (Application) Init () error {
	window, err := nasin.NewWindow(tomo.Bounds(0, 0, 360, 0))
	if err != nil { return err }
	window.SetTitle("Icons")

	container := elements.NewVBox(elements.SpaceBoth)
	window.Adopt(container)

	container.Adopt (
		elements.NewLabel("Just some of the wonderful icons we have:"),
		elements.NewLine())
	container.AdoptExpand (
		icons(tomo.IconHome,   tomo.IconHistory),
		icons(tomo.IconFile,   tomo.IconNetwork),
		icons(tomo.IconOpen,   tomo.IconRemoveFavorite),
		icons(tomo.IconCursor, tomo.IconDistort))

	closeButton := elements.NewButton("Yes verynice")
	closeButton.SetIcon(tomo.IconYes)
	closeButton.OnClick(window.Close)
	container.Adopt(closeButton)
	
	window.OnClose(nasin.Stop)
	window.Show()
	return nil
}

func icons (min, max tomo.Icon) (container *elements.Box) {
	container = elements.NewHBox(elements.SpaceMargin)
	for index := min; index <= max; index ++ {
		container.AdoptExpand(elements.NewIcon(index, tomo.IconSizeSmall))
	}
	return
}
