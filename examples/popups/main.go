package main

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
	window.SetTitle("Dialog Boxes")

	container := elements.NewVBox(elements.SpaceBoth)
	window.Adopt(container)

	container.AdoptExpand(elements.NewLabel("Try out different dialogs:"))

	infoButton := elements.NewButton("popups.DialogKindInfo")
	infoButton.OnClick (func () {
		popups.NewDialog (
			popups.DialogKindInfo,
			window,
			"Information",
			"You are wacky")
	})
	container.Adopt(infoButton)
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
	container.Adopt(questionButton)
	
	warningButton := elements.NewButton("popups.DialogKindWarning")
	warningButton.OnClick (func () {
		popups.NewDialog (
			popups.DialogKindWarning,
			window,
			"Warning",
			"They are fast approaching.")
	})
	container.Adopt(warningButton)
	
	errorButton := elements.NewButton("popups.DialogKindError")
	errorButton.OnClick (func () {
		popups.NewDialog (
			popups.DialogKindError,
			window,
			"Error",
			"There is nowhere left to go.")
	})
	container.Adopt(errorButton)

	menuButton := elements.NewButton("menu")
	menuButton.OnClick (func () {
		// TODO: make a better way to get the bounds of something
		menu, err := window.NewMenu (
			tomo.Bounds(0, 0, 64, 64).
			Add(menuButton.Entity().Bounds().Min))
		if err != nil { println(err.Error()) }
		menu.Adopt(elements.NewLabelWrapped("I'm a shy window..."))
		menu.Show()
	})
	container.Adopt(menuButton)

	cancelButton := elements.NewButton("No thank you.")
	cancelButton.OnClick(nasin.Stop)
	container.Adopt(cancelButton)
		
	window.OnClose(nasin.Stop)
	window.Show()
	return nil
}
