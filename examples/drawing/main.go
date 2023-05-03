package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/nasin"
import "git.tebibyte.media/sashakoshka/tomo/elements/testing"
import "git.tebibyte.media/sashakoshka/ezprof/ez"

func main () {
	nasin.Run(Application { })
}

type Application struct { }

func (Application) Init () error {
	window, err := nasin.NewWindow(tomo.Bounds(0, 0, 480, 360))
	if err != nil { return err }
	window.Adopt(testing.NewArtist())
	window.OnClose(nasin.Stop)
	window.Show()
	ez.Prof()
	return nil
}
