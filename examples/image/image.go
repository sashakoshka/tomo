package main

import "os"
import "image"
import "bytes"
import _ "image/png"
import "github.com/jezek/xgbutil/gopher"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/nasin"
import "git.tebibyte.media/sashakoshka/tomo/popups"
import "git.tebibyte.media/sashakoshka/tomo/elements"

func main () {
	nasin.Run(Application { })
}

type Application struct { }

func (Application) Init () error {
	window, _ := nasin.NewWindow(tomo.Bounds(0, 0, 0, 0))
	window.SetTitle("Tomo Logo")

	file, err := os.Open("assets/banner.png")
	if err != nil { return err }
	logo, _, err := image.Decode(file)
	file.Close()
	if err != nil { return err }

	container := elements.NewVBox(elements.SpaceBoth)
	logoImage := elements.NewImage(logo)
	button    := elements.NewButton("Show me a gopher instead")
	button.OnClick (func () {
		window.SetTitle("Not the Tomo Logo")
		container.DisownAll()
		gopher, _, err :=
			image.Decode(bytes.NewReader(gopher.GopherPng()))
		if err != nil { fatalError(window, err); return }
		container.AdoptExpand(elements.NewImage(gopher))
	})

	container.AdoptExpand(logoImage)
	container.Adopt(button)
	window.Adopt(container)

	button.Focus()
	
	window.OnClose(nasin.Stop)
	window.Show()
	return nil
}

func fatalError (window tomo.Window, err error) {
	popups.NewDialog (
		popups.DialogKindError,
		window,
		"Error",
		err.Error(),
		popups.Button {
			Name: "OK",
			OnPress: nasin.Stop,
		})
} 

