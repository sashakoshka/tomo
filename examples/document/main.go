package main

import "os"
import "image"
import _ "image/png"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/nasin"
import "git.tebibyte.media/sashakoshka/tomo/elements"

func main () {
	nasin.Run(Application { })
}

type Application struct { }

func (Application) Init () error {
	window, err := nasin.NewWindow(tomo.Bounds(0, 0, 383, 360))
	if err != nil { return err }
	window.SetTitle("Document Container")
	
	file, err := os.Open("assets/banner.png")
	if err != nil { return err }
	logo, _, err := image.Decode(file)
	file.Close()
	if err != nil { return err }

	document := elements.NewDocument()
	document.Adopt (
		elements.NewLabelWrapped (
			"A document container is a vertically stacked container " +
			"capable of properly laying out flexible elements such as " +
			"text-wrapped labels. You can also include normal elements " +
			"like:"),
		elements.NewButton("Buttons,"),
		elements.NewCheckbox("Checkboxes,", true),
		elements.NewTextBox("", "And text boxes."),
		elements.NewLine(),
		elements.NewLabelWrapped (
			"Document containers are meant to be placed inside of a " +
			"ScrollContainer, like this one."),
		elements.NewLabelWrapped (
			"You could use document containers to do things like display various " +
			"forms of hypertext (like HTML, gemtext, markdown, etc.), " +
			"lay out a settings menu with descriptive label text between " +
			"control groups like in iOS, or list comment or chat histories."),
		elements.NewImage(logo),
		elements.NewLabelWrapped (
			"You can also choose whether each element is on its own line " +
			"(sort of like an HTML/CSS block element) or on a line with " +
			"other adjacent elements (like an HTML/CSS inline element)."))
	document.AdoptInline (
		elements.NewButton("Just"),
		elements.NewButton("like"),
		elements.NewButton("this."))
	document.Adopt (elements.NewLabelWrapped (
		"Oh, you're a switch? Then name all of these switches:"))
	for i := 0; i < 30; i ++ {
		document.AdoptInline(elements.NewSwitch("", false))
	}

	window.Adopt(elements.NewScroll(elements.ScrollVertical, document))
	window.OnClose(nasin.Stop)
	window.Show()
	return nil
}
