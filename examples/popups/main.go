package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/popups"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/layouts"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("Dialog Boxes")

	container := basic.NewContainer(layouts.Vertical { true, true })
	window.Adopt(container)

	container.Adopt(basic.NewLabel("Try out different dialogs:", false), true)

	infoButton := basic.NewButton("popups.DialogKindInfo")
	infoButton.OnClick (func () {
		popups.NewDialog (
			popups.DialogKindInfo,
			"Information",
			"You are wacky")
	})
	container.Adopt(infoButton, false)
	infoButton.Select()
	
	questionButton := basic.NewButton("popups.DialogKindQuestion")
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
	
	warningButton := basic.NewButton("popups.DialogKindWarning")
	warningButton.OnClick (func () {
		popups.NewDialog (
			popups.DialogKindQuestion,
			"Warning",
			"They are fast approaching.")
	})
	container.Adopt(warningButton, false)
	
	errorButton := basic.NewButton("popups.DialogKindError")
	errorButton.OnClick (func () {
		popups.NewDialog (
			popups.DialogKindQuestion,
			"Error",
			"There is nowhere left to go.")
	})
	container.Adopt(errorButton, false)

	cancelButton := basic.NewButton("No thank you.")
	cancelButton.OnClick(tomo.Stop)
	container.Adopt(cancelButton, false)
		
	window.OnClose(tomo.Stop)
	window.Show()
}
