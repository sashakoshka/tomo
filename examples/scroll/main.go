package main

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(tomo.Bounds(0, 0, 360, 240))
	window.SetTitle("Scroll")
	container := elements.NewVBox(elements.SpaceBoth)
	window.Adopt(container)

	textBox := elements.NewTextBox("", copypasta)

	disconnectedContainer := elements.NewHBox(elements.SpaceMargin)
	list := elements.NewList (
		2,
		elements.NewCell(elements.NewCheckbox("Item 0", true)),
		elements.NewCell(elements.NewCheckbox("Item 1", false)),
		elements.NewCell(elements.NewCheckbox("Item 2", false)),
		elements.NewCell(elements.NewCheckbox("Item 3", true)),
		elements.NewCell(elements.NewCheckbox("Item 4", false)),
		elements.NewCell(elements.NewCheckbox("Item 5", false)),
		elements.NewCell(elements.NewCheckbox("Item 6", false)),
		elements.NewCell(elements.NewCheckbox("Item 7", true)),
		elements.NewCell(elements.NewCheckbox("Item 8", true)),
		elements.NewCell(elements.NewCheckbox("Item 9", false)),
		elements.NewCell(elements.NewCheckbox("Item 10", false)),
		elements.NewCell(elements.NewCheckbox("Item 11", true)),
		elements.NewCell(elements.NewCheckbox("Item 12", false)),
		elements.NewCell(elements.NewCheckbox("Item 13", true)),
		elements.NewCell(elements.NewCheckbox("Item 14", false)),
		elements.NewCell(elements.NewCheckbox("Item 15", false)),
		elements.NewCell(elements.NewCheckbox("Item 16", true)),
		elements.NewCell(elements.NewCheckbox("Item 17", true)),
		elements.NewCell(elements.NewCheckbox("Item 18", false)),
		elements.NewCell(elements.NewCheckbox("Item 19", false)),
		elements.NewCell(elements.NewCheckbox("Item 20", true)),
		elements.NewCell(elements.NewCheckbox("Item 21", false)),
		elements.NewCell(elements.NewScroll (
			elements.ScrollHorizontal,
			elements.NewTextBox("", "I bet you weren't expecting this!"))))
	list.Collapse(0, 32)
	scrollBar := elements.NewVScrollBar()
	list.OnScrollBoundsChange (func () {
		scrollBar.SetBounds (
			list.ScrollContentBounds(),
			list.ScrollViewportBounds())
	})
	scrollBar.OnScroll (func (viewport image.Point) {
		list.ScrollTo(viewport)
	})
	
	container.Adopt(elements.NewLabel("A ScrollContainer:"))
	container.Adopt(elements.NewScroll(elements.ScrollHorizontal, textBox))
	disconnectedContainer.Adopt(list)
	disconnectedContainer.AdoptExpand(elements.NewLabelWrapped (
		"Notice how the scroll bar to the right can be used to " +
		"control the list, despite not even touching it. It is " +
		"indeed a thing you can do. It is also terrible UI design so " +
		"don't do it."))
	disconnectedContainer.Adopt(scrollBar)
	container.AdoptExpand(disconnectedContainer)
	
	window.OnClose(tomo.Stop)
	window.Show()
}

const copypasta = `"I use Linux as my operating system," I state proudly to the unkempt, bearded man. He swivels around in his desk chair with a devilish gleam in his eyes, ready to mansplain with extreme precision. "Actually", he says with a grin, "Linux is just the kernel. You use GNU+Linux!' I don't miss a beat and reply with a smirk, "I use Alpine, a distro that doesn't include the GNU Coreutils, or any other GNU code. It's Linux, but it's not GNU+Linux." The smile quickly drops from the man's face. His body begins convulsing and he foams at the mouth and drops to the floor with a sickly thud. As he writhes around he screams "I-IT WAS COMPILED WITH GCC! THAT MEANS IT'S STILL GNU!" Coolly, I reply "If windows were compiled with GCC, would that make it GNU?" I interrupt his response with "-and work is being made on the kernel to make it more compiler-agnostic. Even if you were correct, you won't be for long." With a sickly wheeze, the last of the man's life is ejected from his body. He lies on the floor, cold and limp. I've womansplained him to death.`
