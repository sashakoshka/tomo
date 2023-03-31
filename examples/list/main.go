package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/popups"
import "git.tebibyte.media/sashakoshka/tomo/layouts"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import "git.tebibyte.media/sashakoshka/tomo/elements/testing"
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(300, 2)
	window.SetTitle("List Sidebar")

	container := containers.NewContainer(layouts.Horizontal { true, true })
	window.Adopt(container)

	var currentPage tomo.Element
	turnPage := func (newPage tomo.Element) {
		container.Warp (func () {
			if currentPage != nil {
				container.Disown(currentPage)
			}
			container.Adopt(newPage, true)
			currentPage = newPage
		})
	}

	intro := elements.NewLabel (
		"The List element can be easily used as a sidebar. " +
		"Click on entries to flip pages!", true)
	button := elements.NewButton("I do nothing!")
	button.OnClick (func () {
		popups.NewDialog(popups.DialogKindInfo, window, "", "Sike!")
	})
	mouse  := testing.NewMouse()
	input  := elements.NewTextBox("Write some text", "")
	form := containers.NewContainer(layouts.Vertical { true, false})
		form.Adopt(elements.NewLabel("I have:", false), false)
		form.Adopt(elements.NewSpacer(true), false)
		form.Adopt(elements.NewCheckbox("Skin", true), false)
		form.Adopt(elements.NewCheckbox("Blood", false), false)
		form.Adopt(elements.NewCheckbox("Bone", false), false)
	art := testing.NewArtist()

	list := elements.NewList (
		elements.NewListEntry("button", func () { turnPage(button) }),
		elements.NewListEntry("mouse",  func () { turnPage(mouse) }),
		elements.NewListEntry("input",  func () { turnPage(input) }),
		elements.NewListEntry("form",   func () { turnPage(form) }),
		elements.NewListEntry("art",    func () { turnPage(art) }))
	list.OnNoEntrySelected(func () { turnPage (intro) })
	list.Collapse(96, 0)
	
	container.Adopt(list, false)
	turnPage(intro)
	
	window.OnClose(tomo.Stop)
	window.Show()
}
