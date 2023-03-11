package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(480, 360)
	window.SetTitle("Scroll")

	scrollContainer := basicElements.NewScrollContainer(false, true)
	document := basicElements.NewDocumentContainer()

	document.Adopt (basicElements.NewLabel (
		"A document container is a vertically stacked container " +
		"capable of properly laying out flexible elements such as " +
		"text-wrapped labels.", true))
	document.Adopt (basicElements.NewButton (
		"You can also include normal elements like buttons,"))
	document.Adopt (basicElements.NewButton (
		"You can also include normal elements like buttons,"))
	document.Adopt (basicElements.NewButton (
		"You can also include normal elements like buttons,"))
	document.Adopt (basicElements.NewButton (
		"You can also include normal elements like buttons,"))
	document.Adopt (basicElements.NewButton (
		"You can also include normal elements like buttons,"))
	document.Adopt (basicElements.NewButton (
		"You can also include normal elements like buttons,"))
	document.Adopt (basicElements.NewCheckbox (
		"checkboxes,", true))
	document.Adopt(basicElements.NewTextBox("", "And text boxes."))
	document.Adopt (basicElements.NewSpacer(true))
	document.Adopt (basicElements.NewLabel (
		"Document containers are meant to be placed inside of a " +
		"ScrollContainer, like this one.", true))
	document.Adopt (basicElements.NewLabel (
		"You could use document containers to do things like display various " +
		"forms of hypertext (like HTML, gemtext, markdown, etc.), " +
		"lay out a settings menu with descriptive label text between " +
		"control groups like in iOS, or list comment or chat histories.", true))

	scrollContainer.Adopt(document)
	window.Adopt(scrollContainer)
	window.OnClose(tomo.Stop)
	window.Show()
}
