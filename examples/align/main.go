package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/textdraw"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(tomo.Bounds(0, 0, 256, 256))
	window.SetTitle("Text alignment")

	left    := elements.NewLabelWrapped(text)
	center  := elements.NewLabelWrapped(text)
	right   := elements.NewLabelWrapped(text)
	justify := elements.NewLabelWrapped(text)

	left.SetAlign(textdraw.AlignLeft)
	center.SetAlign(textdraw.AlignCenter)
	right.SetAlign(textdraw.AlignRight)
	justify.SetAlign(textdraw.AlignJustify)

	window.Adopt (elements.NewScroll (elements.ScrollVertical,
		elements.NewDocument(left, center, right, justify)))
	
	window.OnClose(tomo.Stop)
	window.Show()
}

const text = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Fermentum et sollicitudin ac orci phasellus egestas tellus rutrum. Aliquam vestibulum morbi blandit cursus risus at ultrices mi."
