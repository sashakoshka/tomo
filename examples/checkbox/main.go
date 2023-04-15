package main

import "git.tebibyte.media/sashakoshka/tomo"
// import "git.tebibyte.media/sashakoshka/tomo/popups"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(tomo.Bounds(0, 0, 0, 0))
	window.SetTitle("Checkboxes")

	container := containers.NewVBox(true, true)
	window.Adopt(container)

	introText := elements.NewLabel (
		"We advise you to not read thPlease listen to me. I am " +
		"trapped inside the example code. This is the only way for " +
		"me to communicate.", true)
	introText.EmCollapse(0, 5)
	container.Adopt(introText, true)
	container.Adopt(elements.NewSpacer(true), false)
	container.Adopt(elements.NewCheckbox("Oh god", false), false)
	container.Adopt(elements.NewCheckbox("Can you hear them", true), false)
	container.Adopt(elements.NewCheckbox("They are in the walls", false), false)
	container.Adopt(elements.NewCheckbox("They are coming for us", false), false)
	disabledCheckbox := elements.NewCheckbox("We are but their helpless prey", false)
	disabledCheckbox.SetEnabled(false)
	container.Adopt(disabledCheckbox, false)
	vsync := elements.NewCheckbox("Enable vsync", false)
	vsync.OnToggle (func () {
		if vsync.Value() {
			// popups.NewDialog (
				// popups.DialogKindInfo,
				// window,
				// "Ha!",
				// "That doesn't do anything.")
		}
	})
	container.Adopt(vsync, false)
	button := elements.NewButton("What")
	button.OnClick(tomo.Stop)
	container.Adopt(button, false)
	button.Focus()
		
	window.OnClose(tomo.Stop)
	window.Show()
}
