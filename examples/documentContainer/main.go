package main

import "os"
import "image"
import _ "image/png"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(383, 360)
	window.SetTitle("Scroll")
	
	file, err := os.Open("assets/banner.png")
	if err != nil { panic(err.Error()); return  }
	logo, _, err := image.Decode(file)
	file.Close()
	if err != nil { panic(err.Error()); return  }

	scrollContainer := containers.NewScrollContainer(false, true)
	document := containers.NewDocumentContainer()

	document.Adopt (elements.NewLabel (
		"A document container is a vertically stacked container " +
		"capable of properly laying out flexible elements such as " +
		"text-wrapped labels. You can also include normal elements " +
		"like:", true), true)
	document.Adopt (elements.NewButton (
		"Buttons,"), true)
	document.Adopt (elements.NewCheckbox (
		"Checkboxes,", true), true)
	document.Adopt(elements.NewTextBox("", "And text boxes."), true)
	document.Adopt (elements.NewSpacer(true), true)
	document.Adopt (elements.NewLabel (
		"Document containers are meant to be placed inside of a " +
		"ScrollContainer, like this one.", true), true)
	document.Adopt (elements.NewLabel (
		"You could use document containers to do things like display various " +
		"forms of hypertext (like HTML, gemtext, markdown, etc.), " +
		"lay out a settings menu with descriptive label text between " +
		"control groups like in iOS, or list comment or chat histories.",
		true), true)
	document.Adopt(elements.NewImage(logo), true)
	document.Adopt (elements.NewLabel (
		"You can also choose whether each element is on its own line " +
		"(sort of like an HTML/CSS block element) or on a line with " +
		"other adjacent elements (like an HTML/CSS inline element).",
		true), true)
	document.Adopt(elements.NewButton("Just"), false)
	document.Adopt(elements.NewButton("like"), false)
	document.Adopt(elements.NewButton("this."), false)
	document.Adopt (elements.NewLabel (
		"Oh, you're a switch? Then name all of these switches:",
		true), true)
	for i := 0; i < 30; i ++ {
		document.Adopt(elements.NewSwitch("", false), false)
	}

	scrollContainer.Adopt(document)
	window.Adopt(scrollContainer)
	window.OnClose(tomo.Stop)
	window.Show()
}
