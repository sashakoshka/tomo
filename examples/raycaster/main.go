package main

import "bytes"
import _ "embed"
import _ "image/png"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/popups"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

//go:embed wall.png
var wallTextureBytes []uint8

func main () {
	tomo.Run(run)
}

// FIXME this entire example seems to be broken

func run () {
	window, _ := tomo.NewWindow(tomo.Bounds(0, 0, 640, 480))
	window.SetTitle("Raycaster")

	container := elements.NewVBox(false, false)
	window.Adopt(container)

	wallTexture, _ := TextureFrom(bytes.NewReader(wallTextureBytes))

	game := NewGame (World {
		Data: []int {
		        1,1,1,1,1,1,1,1,1,1,1,1,1,
		        1,0,0,0,0,0,0,0,0,0,0,0,1,
		        1,0,1,1,1,1,1,1,1,0,0,0,1,
		        1,0,0,0,0,0,0,0,1,1,1,0,1,
		        1,0,0,0,0,0,0,0,1,0,0,0,1,
		        1,0,0,0,0,0,0,0,1,0,1,1,1,
		        1,1,1,1,1,1,1,1,1,0,0,0,1,
		        1,0,0,0,0,0,0,0,1,1,0,1,1,
		        1,0,0,1,0,0,0,0,0,0,0,0,1,
		        1,0,1,1,1,0,0,0,0,0,0,0,1,
		        1,0,0,1,0,0,0,0,0,0,0,0,1,
		        1,0,0,0,0,0,0,0,0,0,0,0,1,
		        1,0,0,0,0,1,0,0,0,0,0,0,1,
		        1,1,1,1,1,1,1,1,1,1,1,1,1,
		},
		Stride: 13,
	}, Textures {
		wallTexture,
	})

	topBar := containers.NewHBox(true, true)
	staminaBar := elements.NewProgressBar(game.Stamina())
	healthBar  := elements.NewProgressBar(game.Health())
	
	topBar.Adopt(elements.NewLabel("Stamina:", false), false)
	topBar.Adopt(staminaBar, true)
	topBar.Adopt(elements.NewLabel("Health:", false), false)
	topBar.Adopt(healthBar, true)
	container.Adopt(topBar, false)
	container.Adopt(game, true)
	game.Focus()

	game.OnStatUpdate (func () {
		staminaBar.SetProgress(game.Stamina())
	})
	
	window.OnClose(tomo.Stop)
	window.Show()
	
	popups.NewDialog (
		popups.DialogKindInfo,
		window,
		"Welcome to the backrooms",
		"You've no-clipped into the backrooms!\n" +
		"Move with WASD, and look with the arrow keys.\n" +
		"Keep an eye on your health and stamina.")
}
