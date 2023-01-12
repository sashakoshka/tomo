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
			
			button0 := basic.NewButton("go in the river")
			button0.OnClick(world.SwitchFunc("wet"))
			button1 := basic.NewButton("walk along the river")
			button1.OnClick(world.SwitchFunc("house"))
			button2 := basic.NewButton("turn around")
			button2.OnClick(world.SwitchFunc("bear"))

			container.Warp ( func () {
				container.Adopt(label, true)
				container.Adopt(button0, false)
				container.Adopt(button1, false)
				container.Adopt(button2, false)
				button0.Select()
			})
		},
		"wet": func () {
			label := basic.NewLabel (
				"you get completely soaked.\n" +
				"you die of hypothermia.", false)
			
			button0 := basic.NewButton("try again")
			button0.OnClick(world.SwitchFunc("start"))
			button1 := basic.NewButton("exit")
			button1.OnClick(tomo.Stop)

			container.Warp (func () {
				container.Adopt(label, true)
				container.Adopt(button0, false)
				container.Adopt(button1, false)
				button0.Select()				
			})
		},
		"house": func () {
			label := basic.NewLabel (
				"you are standing in front of a delapidated " +
				"house.", false)
			
			button1 := basic.NewButton("go inside")
			button1.OnClick(world.SwitchFunc("inside"))
			button0 := basic.NewButton("turn back")
			button0.OnClick(world.SwitchFunc("start"))
			
			container.Warp (func () {	
				container.Adopt(label, true)
				container.Adopt(button1, false)
				container.Adopt(button0, false)
				button1.Select()
			})
		},
		"inside": func () {
			label := basic.NewLabel (
				"you are standing inside of the house.\n" +
				"it is dark, but rays of light stream " +
				"through the window.\n" +
				"there is nothing particularly interesting " +
				"here.", false)
			
			button0 := basic.NewButton("go back outside")
			button0.OnClick(world.SwitchFunc("house"))
			
			container.Warp (func () {	
				container.Adopt(label, true)
				container.Adopt(button0, false)
				button0.Select()
			})
		},
		"bear": func () {
			label := basic.NewLabel (
				"you come face to face with a bear.\n" +
				"it eats you (it was hungry).", false)
			
			button0 := basic.NewButton("try again")
			button0.OnClick(world.SwitchFunc("start"))
			button1 := basic.NewButton("exit")
			button1.OnClick(tomo.Stop)
			
			container.Warp (func () {	
				container.Adopt(label, true)
				container.Adopt(button0, false)
				container.Adopt(button1, false)
				button0.Select()
			})
		},
	}
	world.Switch("start")

	window.OnClose(tomo.Stop)
	window.Show()
}
