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
	window.SetTitle("Enter Details")
	container := elements.NewVBox(elements.SpaceBoth)
	window.Adopt(container)

	// create inputs
	firstName    := elements.NewTextBox("First name", "")
	lastName     := elements.NewTextBox("Last name", "")
	fingerLength := elements.NewTextBox("Length of fingers", "")
	purpose      := elements.NewComboBox (
		"",
		"Gaslight",
		"Gatekeep",
		"Girlboss")
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
			firstName.Filled()    &&
			lastName.Filled()     &&
			fingerLength.Filled() &&
			purpose.Filled())
	}
	firstName.OnChange(check)
	lastName.OnChange(check)
	fingerLength.OnChange(check)
	purpose.OnChange(check)

	// add elements to container
	container.AdoptExpand(elements.NewLabel("Choose your words carefully."))
	container.Adopt (
		firstName, lastName,
		fingerLength,
		elements.NewLabel("Purpose:"),
		purpose,
		elements.NewLine(), button)
	window.OnClose(nasin.Stop)
	window.Show()
	return nil
}
