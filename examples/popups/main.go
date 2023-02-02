package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/popups"
import "git.tebibyte.media/sashakoshka/tomo/layouts/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("Dialog Boxes")

	container := basicElements.NewContainer(basicLayouts.Vertical { true, true })
	window.Adopt(container)

	container.Adopt(basicElements.NewLabel("Try out different dialogs:", false), true)

	infoButton := basicElements.NewButton("popups.DialogKindInfo")
	infoButton.OnClick (func () {
		popups.NewDialog (
			popups.DialogKindInfo,
			"Information",
			"You are wacky")
	})
	container.Adopt(infoButton, false)
	infoButton.Focus()
	
	questionButton := basicElements.NewButton("popups.DialogKindQuestion")
	questionButton.OnClick (func () {
		popups.NewDialog (
			popups.DialogKindQuestion,
			"The Big Question",
			"Are you real?",
			popups.Button { "Yes",      func () { } },
			popups.Button { "No",       func () { } },
			popups.Button { "Not sure", func () { } })
	})
	container.Adopt(questionButton, false)
	
	warningButton := basicElements.NewButton("popups.DialogKindWarning")
	warningButton.OnClick (func () {
		popups.NewDialog (
			popups.DialogKindQuestion,
			"Warning",
			"They are fast approaching.")
	})
	container.Adopt(warningButton, false)
	
	errorButton := basicElements.NewButton("popups.DialogKindError")
	errorButton.OnClick (func () {
		popups.NewDialog (
			popups.DialogKindQuestion,
			"Error",
			"There is nowhere left to go.")
	})
	container.Adopt(errorButton, false)

	cancelButton := basicElements.NewButton("No thank you.")
	cancelButton.OnClick(tomo.Stop)
	container.Adopt(cancelButton, false)
		
	window.OnClose(tomo.Stop)
	window.Show()
}
