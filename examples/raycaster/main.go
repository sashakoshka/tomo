package main

import "bytes"
import _ "embed"
import _ "image/png"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/layouts/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"

//go:embed wall.png
var wallTextureBytes []uint8

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(640, 480)
	window.SetTitle("Raycaster")

	container := basicElements.NewContainer(basicLayouts.Vertical { true, true })
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

	container.Adopt(basicElements.NewLabel("Explore a 3D world!", false), false)
	container.Adopt(game, true)
	game.Focus()
	
	window.OnClose(tomo.Stop)
	window.Show()
}
