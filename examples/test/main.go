package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/nasin"
import "git.tebibyte.media/sashakoshka/tomo/elements/testing"

func main () {
	nasin.Run(Application { })
}

type Application struct { }

func (Application) Init () error {
	window, err := nasin.NewWindow(tomo.Bounds(0, 0, 0, 0))
	if err != nil { return err }
	window.SetTitle("Mouse Test")	
	window.Adopt(testing.NewMouse())
	window.OnClose(nasin.Stop)
	window.Show()
	return nil
}
