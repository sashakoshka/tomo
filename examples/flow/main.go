package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/flow"
import "git.tebibyte.media/sashakoshka/tomo/layouts"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(tomo.Bounds(0, 0, 192, 192))
	window.SetTitle("adventure")
	container := containers.NewContainer(layouts.Vertical { true, true })
	window.Adopt(container)

	var world flow.Flow
	world.Transition = container.DisownAll
	world.Stages = map [string] func () {
		"start": func () {
			label := elements.NewLabel (
				"you are standing next to a river.", true)
			
			button0 := elements.NewButton("go in the river")
			button0.OnClick(world.SwitchFunc("wet"))
			button1 := elements.NewButton("walk along the river")
			button1.OnClick(world.SwitchFunc("house"))
			button2 := elements.NewButton("turn around")
			button2.OnClick(world.SwitchFunc("bear"))

			container.Warp ( func () {
				container.Adopt(label, true)
				container.Adopt(button0, false)
				container.Adopt(button1, false)
				container.Adopt(button2, false)
				button0.Focus()
			})
		},
		"wet": func () {
			label := elements.NewLabel (
				"you get completely soaked.\n" +
				"you die of hypothermia.", true)
			
			button0 := elements.NewButton("try again")
			button0.OnClick(world.SwitchFunc("start"))
			button1 := elements.NewButton("exit")
			button1.OnClick(tomo.Stop)

			container.Warp (func () {
				container.Adopt(label, true)
				container.Adopt(button0, false)
				container.Adopt(button1, false)
				button0.Focus()				
			})
		},
		"house": func () {
			label := elements.NewLabel (
				"you are standing in front of a delapidated " +
				"house.", true)
			
			button1 := elements.NewButton("go inside")
			button1.OnClick(world.SwitchFunc("inside"))
			button0 := elements.NewButton("turn back")
			button0.OnClick(world.SwitchFunc("start"))
			
			container.Warp (func () {	
				container.Adopt(label, true)
				container.Adopt(button1, false)
				container.Adopt(button0, false)
				button1.Focus()
			})
		},
		"inside": func () {
			label := elements.NewLabel (
				"you are standing inside of the house.\n" +
				"it is dark, but rays of light stream " +
				"through the window.\n" +
				"there is nothing particularly interesting " +
				"here.", true)
			
			button0 := elements.NewButton("go back outside")
			button0.OnClick(world.SwitchFunc("house"))
			
			container.Warp (func () {	
				container.Adopt(label, true)
				container.Adopt(button0, false)
				button0.Focus()
			})
		},
		"bear": func () {
			label := elements.NewLabel (
				"you come face to face with a bear.\n" +
				"it eats you (it was hungry).", true)
			
			button0 := elements.NewButton("try again")
			button0.OnClick(world.SwitchFunc("start"))
			button1 := elements.NewButton("exit")
			button1.OnClick(tomo.Stop)
			
			container.Warp (func () {	
				container.Adopt(label, true)
				container.Adopt(button0, false)
				container.Adopt(button1, false)
				button0.Focus()
			})
		},
	}
	world.Switch("start")

	window.OnClose(tomo.Stop)
	window.Show()
}
