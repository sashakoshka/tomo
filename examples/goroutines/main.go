package main

import "os"
import "time"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/layouts"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"

func main () {
	tomo.Run(run)
	os.Exit(0)
}

func run () {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("clock")
	container := basic.NewContainer(layouts.Vertical { true, true })
	window.Adopt(container)
	label := basic.NewLabel(formatTime(), false)
	container.Adopt(label, false)
	window.OnClose(tomo.Stop)
	window.Show()
	go tick(label)
}

func formatTime () (timeString string) {
	return time.Now().Format("2006-01-02 15:04:05")
}

func tick (label *basic.Label) {
	for {
		label.SetText(formatTime())
		time.Sleep(time.Second)
	}
}
