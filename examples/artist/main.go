package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements/testing"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"
import _ "net/http/pprof"
import "net/http"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(128, 128)
	window.SetTitle("Draw Test")
	window.Adopt(testing.NewArtist())
	window.OnClose(tomo.Stop)
	window.Show()
	go func () {
		http.ListenAndServe("localhost:6060", nil)
	} ()
}
