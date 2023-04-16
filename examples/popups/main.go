package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/popups"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"

func main () {
	tomo.Run(run)
}

func run () {
	window, err := tomo.NewWindow(tomo.Bounds(0, 0, 0, 0))
	if err != nil { panic(err.Error()) }
	window.SetTitle("Dialog Boxes")

	container := containers.NewVBox(true, true)
	window.Adopt(container)

	container.Adopt(elements.NewLabel("Try out different dialogs:", false), true)

	infoButton := elements.NewButton("popups.DialogKindInfo")
	infoButton.OnClick (func () {
		popups.NewDialog (
			popups.DialogKindInfo,
			window,
			"Information",
			"You are wacky")
	})
	container.Adopt(infoButton, false)
	infoButton.Focus()
	
	questionButton := elements.NewButton("popups.DialogKindQuestion")
	questionButton.OnClick (func () {
		popups.NewDialog (
			popups.DialogKindQuestion,
			window,
			"The Big Question",
			"Are you real?",
			popups.Button { "Yes",      func () { } },
			popups.Button { "No",       func () { } },
			popups.Button { "Not sure", func () { } })
	})
	container.Adopt(questionButton, false)
	
	warningButton := elements.NewButton("popups.DialogKindWarning")
	warningButton.OnClick (func () {
		popups.NewDialog (
			popups.DialogKindWarning,
			window,
			"Warning",
			"They are fast approaching.")
	})
	container.Adopt(warningButton, false)
	
	errorButton := elements.NewButton("popups.DialogKindError")
	errorButton.OnClick (func () {
		popups.NewDialog (
			popups.DialogKindError,
			window,
			"Error",
			"There is nowhere left to go.")
	})
	container.Adopt(errorButton, false)

	menuButton := elements.NewButton("menu")
	menuButton.OnClick (func () {
		// TODO: make a better way to get the bounds of something
		menu, err := window.NewMenu (
			tomo.Bounds(0, 0, 64, 64).
			Add(menuButton.Entity().Bounds().Min))
		if err != nil { println(err.Error()) }
		menu.Adopt(elements.NewLabel("I'm a shy window...", true))
		menu.Show()
	})
	container.Adopt(menuButton, false)

	cancelButton := elements.NewButton("No thank you.")
	cancelButton.OnClick(tomo.Stop)
	container.Adopt(cancelButton, false)
		
	window.OnClose(tomo.Stop)
	window.Show()
}
