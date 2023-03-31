package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
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
	container.Adopt(icons(theme.IconHome, theme.IconRepositories), true)
	container.Adopt(icons(theme.IconFile, theme.IconCD), true)
	container.Adopt(icons(theme.IconOpen, theme.IconRemoveBookmark), true)

	closeButton := elements.NewButton("Ok")
	closeButton.SetIcon(theme.IconYes)
	closeButton.ShowText(false)
	closeButton.OnClick(tomo.Stop)
	container.Adopt(closeButton, false)
	
	window.OnClose(tomo.Stop)
	window.Show()
}

func icons (min, max theme.Icon) (container *containers.Container) {
	container = containers.NewContainer(layouts.Horizontal { true, false })
	for index := min; index <= max; index ++ {
		container.Adopt(elements.NewIcon(index, theme.IconSizeSmall), true)
	}
	return
}
