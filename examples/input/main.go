package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/popups"
import "git.tebibyte.media/sashakoshka/tomo/layouts"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("Approaching")
	container := basic.NewContainer(layouts.Vertical { true, true })
	window.Adopt(container)

	firstName    := basic.NewTextBox("First name", "")
	lastName     := basic.NewTextBox("Last name", "")
	fingerLength := basic.NewTextBox("Length of fingers", "")
	button       := basic.NewButton("Ok")

	lastName.SetEnabled(false)
	button.OnClick (func () {
		popups.NewDialog (
			popups.DialogKindInfo,
			"Profile",
			firstName.Value() + " [REDACTED]'s fingers\n" +
			"measure in at " + fingerLength.Value() + " feet.")
	})
	
	container.Adopt(basic.NewLabel("Choose your words carefully.", false), true)
	container.Adopt(firstName, false)
	container.Adopt(lastName, false)
	container.Adopt(fingerLength, false)
	container.Adopt(basic.NewSpacer(true), false)
	container.Adopt(button, false)

	firstName.Select()
	
	window.OnClose(tomo.Stop)
	window.Show()
}
