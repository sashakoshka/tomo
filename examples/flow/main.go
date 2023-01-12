package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/flow"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/layouts"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("adventure")
	container := basic.NewContainer(layouts.Vertical { true, true })
	window.Adopt(container)

	var world flow.Flow
	world.Transition = container.DisownAll
	world.Stages = map [string] func () {
		"start": func () {
			label := basic.NewLabel (
				"you are standing next to a river.", false)
			container.Adopt(label, true)
			
			button0 := basic.NewButton("go in the river")
			button0.OnClick(world.SwitchFunc("wet"))
			container.Adopt(button0, false)
			button0.Select()
			
			button1 := basic.NewButton("walk along the river")
			button1.OnClick(world.SwitchFunc("house"))
			container.Adopt(button1, false)
			
			button2 := basic.NewButton("turn around")
			button2.OnClick(world.SwitchFunc("bear"))
			container.Adopt(button2, false)
		},
		"wet": func () {
			label := basic.NewLabel (
				"you get completely soaked.\n" +
				"you die of hypothermia.", false)
			container.Adopt(label, true)
			
			button0 := basic.NewButton("try again")
			button0.OnClick(world.SwitchFunc("start"))
			container.Adopt(button0, false)
			button0.Select()
			
			button1 := basic.NewButton("exit")
			button1.OnClick(tomo.Stop)
			container.Adopt(button1, false)
		},
		"house": func () {
			label := basic.NewLabel (
				"you are standing in front of a delapidated " +
				"house.", false)
			container.Adopt(label, true)
			
			button1 := basic.NewButton("go inside")
			button1.OnClick(world.SwitchFunc("inside"))
			container.Adopt(button1, false)
			button1.Select()
			
			button0 := basic.NewButton("turn back")
			button0.OnClick(world.SwitchFunc("start"))
			container.Adopt(button0, false)
		},
		"inside": func () {
			label := basic.NewLabel (
				"you are standing inside of the house.\n" +
				"it is dark, but rays of light stream " +
				"through the window.\n" +
				"there is nothing particularly interesting " +
				"here.", false)
			container.Adopt(label, true)
			
			button0 := basic.NewButton("go back outside")
			button0.OnClick(world.SwitchFunc("house"))
			container.Adopt(button0, false)
			button0.Select()
		},
		"bear": func () {
			label := basic.NewLabel (
				"you come face to face with a bear.\n" +
				"it eats you (it was hungry).", false)
			container.Adopt(label, true)
			
			button0 := basic.NewButton("try again")
			button0.OnClick(world.SwitchFunc("start"))
			container.Adopt(button0, false)
			button0.Select()
			
			button1 := basic.NewButton("exit")
			button1.OnClick(tomo.Stop)
			container.Adopt(button1, false)
		},
	}
	world.Switch("start")

	window.OnClose(tomo.Stop)
	window.Show()
}
