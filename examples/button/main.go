package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("example button")
	button := basic.NewButton("hello tomo!")
	button.OnClick (func () {
		// when we set the button's text to something longer, the window
		// will automatically resize to accomodate it.
		button.SetText("you clicked me.\nwow, there are two lines!")
		button.OnClick (func () {
			button.SetText (
				"stop clicking me you idiot!\n" +
				"you've already seen it all!")
			button.OnClick(tomo.Stop)
		})
	})
	window.Adopt(button)
	window.OnClose(tomo.Stop)
	window.Show()
}
