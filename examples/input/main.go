package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/popups"
import "git.tebibyte.media/sashakoshka/tomo/layouts"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(tomo.Bounds(0, 0, 0, 0))
	window.SetTitle("Enter Details")
	container := containers.NewContainer(layouts.Vertical { true, true })
	window.Adopt(container)

	// create inputs
	firstName    := elements.NewTextBox("First name", "")
	lastName     := elements.NewTextBox("Last name", "")
	fingerLength := elements.NewTextBox("Length of fingers", "")
	button       := elements.NewButton("Ok")

	button.SetEnabled(false)
	button.OnClick (func () {
		// create a dialog displaying the results
		popups.NewDialog (
			popups.DialogKindInfo,
			window,
			"Profile",
			firstName.Value() + " " + lastName.Value() +
			"'s fingers\nmeasure in at " + fingerLength.Value() +
			" feet.")
	})

	// enable the Ok button if all three inputs have text in them
	check := func () {
		button.SetEnabled (
			firstName.Filled() &&
			lastName.Filled() &&
			fingerLength.Filled())
	}
	firstName.OnChange(check)
	lastName.OnChange(check)
	fingerLength.OnChange(check)

	// add elements to container
	container.Adopt(elements.NewLabel("Choose your words carefully.", false), true)
	container.Adopt(firstName, false)
	container.Adopt(lastName, false)
	container.Adopt(fingerLength, false)
	container.Adopt(elements.NewSpacer(true), false)
	container.Adopt(button, false)
	
	window.OnClose(tomo.Stop)
	window.Show()
}
