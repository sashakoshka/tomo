package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/popups"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import "git.tebibyte.media/sashakoshka/tomo/layouts/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/testing"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(300, 2)
	window.SetTitle("List Sidebar")

	container := basicElements.NewContainer(basicLayouts.Horizontal { true, true })
	window.Adopt(container)

	var currentPage elements.Element
	turnPage := func (newPage elements.Element) {
		container.Warp (func () {
			if currentPage != nil {
				container.Disown(currentPage)
			}
			container.Adopt(newPage, true)
			currentPage = newPage
		})
	}

	intro := basicElements.NewLabel (
		"The List element can be easily used as a sidebar. " +
		"Click on entries to flip pages!", true)
	button := basicElements.NewButton("I do nothing!")
	button.OnClick (func () {
		popups.NewDialog(popups.DialogKindInfo, "", "Sike!")
	})
	mouse  := testing.NewMouse()
	input  := basicElements.NewTextBox("Write some text", "")
	form := basicElements.NewContainer(basicLayouts.Vertical { true, false})
		form.Adopt(basicElements.NewLabel("I have:", false), false)
		form.Adopt(basicElements.NewSpacer(true), false)
		form.Adopt(basicElements.NewCheckbox("Skin", true), false)
		form.Adopt(basicElements.NewCheckbox("Blood", false), false)
		form.Adopt(basicElements.NewCheckbox("Bone", false), false)
	art := testing.NewArtist()

	list := basicElements.NewList (
		basicElements.NewListEntry("button", func () { turnPage(button) }),
		basicElements.NewListEntry("mouse",  func () { turnPage(mouse) }),
		basicElements.NewListEntry("input",  func () { turnPage(input) }),
		basicElements.NewListEntry("form",   func () { turnPage(form) }),
		basicElements.NewListEntry("art",    func () { turnPage(art) }))
	list.OnNoEntrySelected(func () { turnPage (intro) })
	list.Collapse(96, 0)
	
	container.Adopt(list, false)
	turnPage(intro)
	
	window.OnClose(tomo.Stop)
	window.Show()
}
