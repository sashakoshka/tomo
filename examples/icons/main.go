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
	window, _ := tomo.NewWindow(360, 2)
	window.SetTitle("Icons")

	container := containers.NewContainer(layouts.Vertical { true, true })
	window.Adopt(container)

	container.Adopt(elements.NewLabel("Just some of the wonderful icons we have:", false), false)
	container.Adopt(elements.NewSpacer(true), false)
	container.Adopt(icons(tomo.IconHome, tomo.IconHistory), true)
	container.Adopt(icons(tomo.IconFile, tomo.IconNetwork), true)
	container.Adopt(icons(tomo.IconOpen, tomo.IconRemoveFavorite), true)
	container.Adopt(icons(tomo.IconCursor, tomo.IconDistort), true)

	closeButton := elements.NewButton("Ok")
	closeButton.SetIcon(tomo.IconYes)
	closeButton.ShowText(false)
	closeButton.OnClick(tomo.Stop)
	container.Adopt(closeButton, false)
	
	window.OnClose(tomo.Stop)
	window.Show()
}

func icons (min, max tomo.Icon) (container *containers.Container) {
	container = containers.NewContainer(layouts.Horizontal { true, false })
	for index := min; index <= max; index ++ {
		container.Adopt(elements.NewIcon(index, tomo.IconSizeSmall), true)
	}
	return
}
