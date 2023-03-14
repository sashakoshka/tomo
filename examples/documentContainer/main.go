package main

import "os"
import "image"
import _ "image/png"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/layouts/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"

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

	scrollContainer := basicElements.NewScrollContainer(false, true)
	document := basicElements.NewDocumentContainer()

	document.Adopt (basicElements.NewLabel (
		"A document container is a vertically stacked container " +
		"capable of properly laying out flexible elements such as " +
		"text-wrapped labels. You can also include normal elements " +
		"like:", true))
	document.Adopt (basicElements.NewButton (
		"Buttons,"))
	document.Adopt (basicElements.NewCheckbox (
		"Checkboxes,", true))
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
	document.Adopt(basicElements.NewImage(logo))
	document.Adopt (basicElements.NewLabel (
		"Oh, you're a switch? Then name all of these switches:", true))
	for i := 0; i < 3; i ++ {
		switchContainer := basicElements.NewContainer (basicLayouts.Horizontal {
			Gap: true,
		})
		for i := 0; i < 10; i ++ {
			switchContainer.Adopt(basicElements.NewSwitch("", false), true)
		}
		document.Adopt(switchContainer)
	}

	scrollContainer.Adopt(document)
	window.Adopt(scrollContainer)
	window.OnClose(tomo.Stop)
	window.Show()
}
