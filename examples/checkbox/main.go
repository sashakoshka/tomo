package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/popups"
import "git.tebibyte.media/sashakoshka/tomo/layouts/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("Checkboxes")

	container := containers.NewContainer(basicLayouts.Vertical { true, true })
	window.Adopt(container)

	introText := basicElements.NewLabel (
		"We advise you to not read thPlease listen to me. I am " +
		"trapped inside the example code. This is the only way for " +
		"me to communicate.", true)
	introText.EmCollapse(0, 5)
	container.Adopt(introText, true)
	container.Adopt(basicElements.NewSpacer(true), false)
	container.Adopt(basicElements.NewCheckbox("Oh god", false), false)
	container.Adopt(basicElements.NewCheckbox("Can you hear them", true), false)
	container.Adopt(basicElements.NewCheckbox("They are in the walls", false), false)
	container.Adopt(basicElements.NewCheckbox("They are coming for us", false), false)
	disabledCheckbox := basicElements.NewCheckbox("We are but their helpless prey", false)
	disabledCheckbox.SetEnabled(false)
	container.Adopt(disabledCheckbox, false)
	vsync := basicElements.NewCheckbox("Enable vsync", false)
	vsync.OnToggle (func () {
		if vsync.Value() {
			popups.NewDialog (
				popups.DialogKindInfo,
				"Ha!",
				"That doesn't do anything.")
		}
	})
	container.Adopt(vsync, false)
	button := basicElements.NewButton("What")
	button.OnClick(tomo.Stop)
	container.Adopt(button, false)
	button.Focus()
		
	window.OnClose(tomo.Stop)
	window.Show()
}
