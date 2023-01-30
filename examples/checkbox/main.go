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
	window.SetTitle("Checkboxes")

	container := basic.NewContainer(layouts.Vertical { true, true })
	window.Adopt(container)

	container.Adopt (basic.NewLabel (
		"We advise you to not read thPlease listen to me. I am " +
		"trapped inside the example code. This is the only way for " +
		"me to communicate.", true), true)
	container.Adopt(basic.NewSpacer(true), false)
	container.Adopt(basic.NewCheckbox("Oh god", false), false)
	container.Adopt(basic.NewCheckbox("Can you hear them", true), false)
	container.Adopt(basic.NewCheckbox("They are in the walls", false), false)
	container.Adopt(basic.NewCheckbox("They are coming for us", false), false)
	disabledCheckbox := basic.NewCheckbox("We are but their helpless prey", false)
	disabledCheckbox.SetEnabled(false)
	container.Adopt(disabledCheckbox, false)
	vsync := basic.NewCheckbox("Enable vsync", false)
	vsync.OnToggle (func () {
		if vsync.Value() {
			popups.NewDialog (
				popups.DialogKindInfo,
				"Ha!",
				"That doesn't do anything.")
		}
	})
	container.Adopt(vsync, false)
	button := basic.NewButton("What")
	button.OnClick(tomo.Stop)
	container.Adopt(button, false)
	button.Focus()
		
	window.OnClose(tomo.Stop)
	window.Show()
}
