package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/popups"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(tomo.Bounds(0, 0, 0, 0))
	window.SetTitle("Checkboxes")

	introText := elements.NewLabelWrapped (
		"We advise you to not read thPlease listen to me. I am " +
		"trapped inside the example code. This is the only way for " +
		"me to communicate.")
	introText.EmCollapse(0, 5)
	
	disabledCheckbox := elements.NewCheckbox("We are but their helpless prey", false)
	disabledCheckbox.SetEnabled(false)
	
	vsync := elements.NewCheckbox("Enable vsync", false)
	vsync.OnToggle (func () {
		if vsync.Value() {
			popups.NewDialog (
				popups.DialogKindInfo,
				window,
				"Ha!",
				"That doesn't do anything.")
		}
	})
	
	button := elements.NewButton("What")
	button.OnClick(tomo.Stop)
	
	box := elements.NewVBox(elements.SpaceBoth)
	box.AdoptExpand(introText)
	box.Adopt (
		elements.NewLine(),
		elements.NewCheckbox("Oh god", false),
		elements.NewCheckbox("Can you hear them", true),
		elements.NewCheckbox("They are in the walls", false),
		elements.NewCheckbox("They are coming for us", false),
		disabledCheckbox,
		vsync, button)
	window.Adopt(box)
		
	button.Focus()
	window.OnClose(tomo.Stop)
	window.Show()
}
