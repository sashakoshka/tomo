package main

import "tomo"
import "tomo/nasin"
import "tomo/popups"
import "tomo/ability"
import "tomo/elements"
import "tomo/elements/testing"

func main () {
	nasin.Run(Application { })
}

type Application struct { }

func (Application) Init () error {
	window, err := nasin.NewWindow(tomo.Bounds(0, 0, 300, 0))
	if err != nil { return err }
	window.SetTitle("List Sidebar")

	container := elements.NewHBox(elements.SpaceBoth)
	window.Adopt(container)

	var currentPage tomo.Element
	turnPage := func (newPage tomo.Element) {
		if currentPage != nil {
			container.Disown(currentPage)
		}
		container.AdoptExpand(newPage)
		currentPage = newPage
	}

	intro := elements.NewLabelWrapped (
		"The List element can be easily used as a sidebar. " +
		"Click on entries to flip pages!")
	button := elements.NewButton("I do nothing!")
	button.OnClick (func () {
		popups.NewDialog(popups.DialogKindInfo, window, "", "Sike!")
	})
	mouse  := testing.NewMouse()
	input  := elements.NewTextBox("Write some text", "")
	form := elements.NewVBox (
		elements.SpaceMargin,
		elements.NewLabel("I have:"),
		elements.NewLine(),
		elements.NewCheckbox("Skin", true),
		elements.NewCheckbox("Blood", false),
		elements.NewCheckbox("Bone", false))
	art := testing.NewArtist()

	makePage := func (name string, callback func ()) ability.Selectable {
		cell := elements.NewCell(elements.NewLabel(name))
		cell.OnSelectionChange (func () {
			if cell.Selected() { callback() }
		})
		return cell
	}

	list := elements.NewList (
		makePage("button", func () { turnPage(button) }),
		makePage("mouse",  func () { turnPage(mouse) }),
		makePage("input",  func () { turnPage(input) }),
		makePage("form",   func () { turnPage(form) }),
		makePage("art",    func () { turnPage(art) }))
	list.Collapse(96, 0)
	
	container.Adopt(list)
	turnPage(intro)
	
	window.OnClose(nasin.Stop)
	window.Show()
	return nil
}
