package main

import "bytes"
import _ "embed"
import _ "image/png"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/popups"
import "git.tebibyte.media/sashakoshka/tomo/layouts/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"

//go:embed wall.png
var wallTextureBytes []uint8

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(640, 480)
	window.SetTitle("Raycaster")

	container := containers.NewContainer(basicLayouts.Vertical { false, false })
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

	topBar := containers.NewContainer(basicLayouts.Horizontal { true, true })
	staminaBar := basicElements.NewProgressBar(game.Stamina())
	healthBar  := basicElements.NewProgressBar(game.Health())
	
	topBar.Adopt(basicElements.NewLabel("Stamina:", false), false)
	topBar.Adopt(staminaBar, true)
	topBar.Adopt(basicElements.NewLabel("Health:", false), false)
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
		"Welcome to the backrooms",
		"You've no-clipped into the backrooms!\n" +
		"Move with WASD, and look with the arrow keys.\n" +
		"Keep an eye on your health and stamina.")
}
