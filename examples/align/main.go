package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/textdraw"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(tomo.Bounds(0, 0, 0, 0))
	window.SetTitle("Text alignment")

	container := containers.NewDocumentContainer()
	scrollContainer := containers.NewScrollContainer(false, true)
	scrollContainer.Adopt(container)
	window.Adopt(scrollContainer)

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
	
	window.OnClose(tomo.Stop)
	window.Show()
}

const text = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Fermentum et sollicitudin ac orci phasellus egestas tellus rutrum. Aliquam vestibulum morbi blandit cursus risus at ultrices mi."
