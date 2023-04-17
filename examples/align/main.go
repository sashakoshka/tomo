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

	container := elements.NewDocument()

	left    := elements.NewLabel(text, true)
	center  := elements.NewLabel(text, true)
	right   := elements.NewLabel(text, true)
	justify := elements.NewLabel(text, true)

	left.SetAlign(textdraw.AlignLeft)
	center.SetAlign(textdraw.AlignCenter)
	right.SetAlign(textdraw.AlignRight)
	justify.SetAlign(textdraw.AlignJustify)

	container.Adopt(left, true)
	container.Adopt(center, true)
	container.Adopt(right, true)
	container.Adopt(justify, true)
	window.Adopt(elements.NewScroll(container, false, true))
	
	window.OnClose(tomo.Stop)
	window.Show()
}

const text = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Fermentum et sollicitudin ac orci phasellus egestas tellus rutrum. Aliquam vestibulum morbi blandit cursus risus at ultrices mi."
