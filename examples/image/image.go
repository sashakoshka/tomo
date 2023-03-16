package main

import "os"
import "image"
import "bytes"
import _ "image/png"
import "github.com/jezek/xgbutil/gopher"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/popups"
import "git.tebibyte.media/sashakoshka/tomo/layouts/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("Tomo Logo")

	file, err := os.Open("assets/banner.png")
	if err != nil { fatalError(err); return  }
	logo, _, err := image.Decode(file)
	file.Close()
	if err != nil { fatalError(err); return  }

	container := basicElements.NewContainer(basicLayouts.Vertical { true, true })
	logoImage := basicElements.NewImage(logo)
	button    := basicElements.NewButton("Show me a gopher instead")
	button.OnClick (func () { container.Warp (func () {
			container.DisownAll()
			gopher, _, err :=
				image.Decode(bytes.NewReader(gopher.GopherPng()))
			if err != nil { fatalError(err); return }
			container.Adopt(basicElements.NewImage(gopher),true)
	}) })

	container.Adopt(logoImage, true)
	container.Adopt(button, false)
	window.Adopt(container)

	button.Focus()
	
	window.OnClose(tomo.Stop)
	window.Show()
}

func fatalError (err error) {
	popups.NewDialog (
		popups.DialogKindError,
		"Error",
		err.Error(),
		popups.Button {
			Name: "OK",
			OnPress: tomo.Stop,
		})
} 
