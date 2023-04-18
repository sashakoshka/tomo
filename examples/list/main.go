package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/popups"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import "git.tebibyte.media/sashakoshka/tomo/elements/testing"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(tomo.Bounds(0, 0, 300, 0))
	window.SetTitle("List Sidebar")

	container := elements.NewHBox(true, true)
	window.Adopt(container)

	var currentPage tomo.Element
	turnPage := func (newPage tomo.Element) {
		if currentPage != nil {
			container.Disown(currentPage)
		}
		container.Adopt(newPage, true)
		currentPage = newPage
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
	form := elements.NewVBox(false, true)
		form.Adopt(elements.NewLabel("I have:", false), false)
		form.Adopt(elements.NewSpacer(true), false)
		form.Adopt(elements.NewCheckbox("Skin", true), false)
		form.Adopt(elements.NewCheckbox("Blood", false), false)
		form.Adopt(elements.NewCheckbox("Bone", false), false)
	art := testing.NewArtist()

	makePage := func (name string, callback func ()) tomo.Selectable {
		cell := elements.NewCell(elements.NewLabel(name, false))
		cell.OnSelectionChange (func () {
			if cell.Selected() { callback() }
		})
		return cell
	}

	list := elements.NewList (
		1,
		makePage("button", func () { turnPage(button) }),
		makePage("mouse",  func () { turnPage(mouse) }),
		makePage("input",  func () { turnPage(input) }),
		makePage("form",   func () { turnPage(form) }),
		makePage("art",    func () { turnPage(art) }))
	list.Collapse(96, 0)
	
	container.Adopt(list, false)
	turnPage(intro)
	
	window.OnClose(tomo.Stop)
	window.Show()
}
