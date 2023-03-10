package main

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/layouts/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(480, 360)
	window.SetTitle("Scroll")
	container := basicElements.NewContainer(basicLayouts.Vertical { true, true })
	window.Adopt(container)

	textBox := basicElements.NewTextBox("", copypasta)
	scrollContainer := basicElements.NewScrollContainer(true, false)

	disconnectedContainer := basicElements.NewContainer (basicLayouts.Horizontal {
		Gap: true,
	})
	list := basicElements.NewList (
		basicElements.NewListEntry("This is list item 0", nil),
		basicElements.NewListEntry("This is list item 1", nil),
		basicElements.NewListEntry("This is list item 2", nil),
		basicElements.NewListEntry("This is list item 3", nil),
		basicElements.NewListEntry("This is list item 4", nil),
		basicElements.NewListEntry("This is list item 5", nil),
		basicElements.NewListEntry("This is list item 6", nil),
		basicElements.NewListEntry("This is list item 7", nil),
		basicElements.NewListEntry("This is list item 8", nil),
		basicElements.NewListEntry("This is list item 9", nil),
		basicElements.NewListEntry("This is list item 10", nil),
		basicElements.NewListEntry("This is list item 11", nil),
		basicElements.NewListEntry("This is list item 12", nil),
		basicElements.NewListEntry("This is list item 13", nil),
		basicElements.NewListEntry("This is list item 14", nil),
		basicElements.NewListEntry("This is list item 15", nil),
		basicElements.NewListEntry("This is list item 16", nil),
		basicElements.NewListEntry("This is list item 17", nil),
		basicElements.NewListEntry("This is list item 18", nil),
		basicElements.NewListEntry("This is list item 19", nil),
		basicElements.NewListEntry("This is list item 20", nil))
	list.Collapse(0, 32)
	scrollBar := basicElements.NewScrollBar(true)
	list.OnScrollBoundsChange (func () {
		scrollBar.SetBounds (
			list.ScrollContentBounds(),
			list.ScrollViewportBounds())
	})
	scrollBar.OnScroll (func (viewport image.Point) {
		list.ScrollTo(viewport)
	})
	
	scrollContainer.Adopt(textBox)
	container.Adopt(basicElements.NewLabel("A ScrollContainer:", false), false)
	container.Adopt(scrollContainer, false)
	disconnectedContainer.Adopt(list, false)
	disconnectedContainer.Adopt (basicElements.NewLabel (
		"Notice how the scroll bar to the right can be used to " +
		"control the list, despite not even touching it. It is " +
		"indeed a thing you can do. It is also terrible UI design so " +
		"don't do it.", true), true)
	disconnectedContainer.Adopt(scrollBar, false)
	container.Adopt(disconnectedContainer, true)
	
	window.OnClose(tomo.Stop)
	window.Show()
}

const copypasta = `"I use Linux as my operating system," I state proudly to the unkempt, bearded man. He swivels around in his desk chair with a devilish gleam in his eyes, ready to mansplain with extreme precision. "Actually", he says with a grin, "Linux is just the kernel. You use GNU+Linux!' I don't miss a beat and reply with a smirk, "I use Alpine, a distro that doesn't include the GNU Coreutils, or any other GNU code. It's Linux, but it's not GNU+Linux." The smile quickly drops from the man's face. His body begins convulsing and he foams at the mouth and drops to the floor with a sickly thud. As he writhes around he screams "I-IT WAS COMPILED WITH GCC! THAT MEANS IT'S STILL GNU!" Coolly, I reply "If windows were compiled with GCC, would that make it GNU?" I interrupt his response with "-and work is being made on the kernel to make it more compiler-agnostic. Even if you were correct, you won't be for long." With a sickly wheeze, the last of the man's life is ejected from his body. He lies on the floor, cold and limp. I've womansplained him to death.`
