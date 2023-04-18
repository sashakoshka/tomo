package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/flow"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(tomo.Bounds(0, 0, 192, 192))
	window.SetTitle("adventure")
	container := elements.NewVBox(elements.SpaceBoth)
	window.Adopt(container)

	var world flow.Flow
	world.Transition = container.DisownAll
	world.Stages = map [string] func () {
		"start": func () {
			label := elements.NewLabelWrapped (
				"you are standing next to a river.")
			
			button0 := elements.NewButton("go in the river")
			button0.OnClick(world.SwitchFunc("wet"))
			button1 := elements.NewButton("walk along the river")
			button1.OnClick(world.SwitchFunc("house"))
			button2 := elements.NewButton("turn around")
			button2.OnClick(world.SwitchFunc("bear"))

			container.AdoptExpand(label)
			container.Adopt(button0, button1, button2)
			button0.Focus()
		},
		"wet": func () {
			label := elements.NewLabelWrapped (
				"you get completely soaked.\n" +
				"you die of hypothermia.")
			
			button0 := elements.NewButton("try again")
			button0.OnClick(world.SwitchFunc("start"))
			button1 := elements.NewButton("exit")
			button1.OnClick(tomo.Stop)

			container.AdoptExpand(label)
			container.Adopt(button0, button1)
			button0.Focus()				
		},
		"house": func () {
			label := elements.NewLabelWrapped (
				"you are standing in front of a delapidated " +
				"house.")
			
			button1 := elements.NewButton("go inside")
			button1.OnClick(world.SwitchFunc("inside"))
			button0 := elements.NewButton("turn back")
			button0.OnClick(world.SwitchFunc("start"))
			
			container.AdoptExpand(label)
			container.Adopt(button0, button1)
			button1.Focus()
		},
		"inside": func () {
			label := elements.NewLabelWrapped (
				"you are standing inside of the house.\n" +
				"it is dark, but rays of light stream " +
				"through the window.\n" +
				"there is nothing particularly interesting " +
				"here.")
			
			button0 := elements.NewButton("go back outside")
			button0.OnClick(world.SwitchFunc("house"))
			
			container.AdoptExpand(label)
			container.Adopt(button0)
			button0.Focus()
		},
		"bear": func () {
			label := elements.NewLabelWrapped (
				"you come face to face with a bear.\n" +
				"it eats you (it was hungry).")
			
			button0 := elements.NewButton("try again")
			button0.OnClick(world.SwitchFunc("start"))
			button1 := elements.NewButton("exit")
			button1.OnClick(tomo.Stop)
			
			container.AdoptExpand(label)
			container.Adopt(button0, button1)
			button0.Focus()
		},
	}
	world.Switch("start")

	window.OnClose(tomo.Stop)
	window.Show()
}
