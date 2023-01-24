package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/popups"
import "git.tebibyte.media/sashakoshka/tomo/layouts"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/testing"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(300, 2)
	window.SetTitle("List Sidebar")

	container := basic.NewContainer(layouts.Horizontal { true, true })
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

	button := basic.NewButton("I do nothing!")
	button.OnClick (func () {
		popups.NewDialog(popups.DialogKindInfo, "", "Sike!")
	})
	mouse  := testing.NewMouse()
	input  := basic.NewTextBox("Write some text", "")
	form := basic.NewContainer(layouts.Vertical { true, false})
		form.Adopt(basic.NewLabel("I have:", false), false)
		form.Adopt(basic.NewSpacer(true), false)
		form.Adopt(basic.NewCheckbox("Skin", true), false)
		form.Adopt(basic.NewCheckbox("Blood", false), false)
		form.Adopt(basic.NewCheckbox("Bone", false), false)

	list := basic.NewList (
		basic.NewListEntry("button", func () { turnPage(button) }),
		basic.NewListEntry("mouse",  func () { turnPage(mouse) }),
		basic.NewListEntry("input",  func () { turnPage(input) }),
		basic.NewListEntry("form",   func () { turnPage(form) }))
	list.Collapse(96, 0)
	
	container.Adopt(list, false)
	turnPage (basic.NewLabel (
		"The List element can be easily used as a sidebar. " +
		"Click on entries to flip pages!", true))
	
	window.OnClose(tomo.Stop)
	window.Show()
}
