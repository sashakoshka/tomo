package main

import "os"
import "image"
import _ "image/png"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/layouts"
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
		"like:", true))
	document.Adopt (elements.NewButton (
		"Buttons,"))
	document.Adopt (elements.NewCheckbox (
		"Checkboxes,", true))
	document.Adopt(elements.NewTextBox("", "And text boxes."))
	document.Adopt (elements.NewSpacer(true))
	document.Adopt (elements.NewLabel (
		"Document containers are meant to be placed inside of a " +
		"ScrollContainer, like this one.", true))
	document.Adopt (elements.NewLabel (
		"You could use document containers to do things like display various " +
		"forms of hypertext (like HTML, gemtext, markdown, etc.), " +
		"lay out a settings menu with descriptive label text between " +
		"control groups like in iOS, or list comment or chat histories.", true))
	document.Adopt(elements.NewImage(logo))
	document.Adopt (elements.NewLabel (
		"Oh, you're a switch? Then name all of these switches:", true))
	for i := 0; i < 3; i ++ {
		switchContainer := containers.NewContainer (layouts.Horizontal {
			Gap: true,
		})
		for i := 0; i < 10; i ++ {
			switchContainer.Adopt(elements.NewSwitch("", false), true)
		}
		document.Adopt(switchContainer)
	}

	scrollContainer.Adopt(document)
	window.Adopt(scrollContainer)
	window.OnClose(tomo.Stop)
	window.Show()
}
