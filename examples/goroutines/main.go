package main

import "os"
import "time"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import "git.tebibyte.media/sashakoshka/tomo/elements/fun"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"

func main () {
	tomo.Run(run)
	os.Exit(0)
}

func run () {
	window, _ := tomo.NewWindow(tomo.Bounds(0, 0, 200, 216))
	window.SetTitle("Clock")
	container := containers.NewVBox(true, true)
	window.Adopt(container)

	clock := fun.NewAnalogClock(time.Now())
	container.Adopt(clock, true)
	label := elements.NewLabel(formatTime(), false)
	container.Adopt(label, false)
	
	window.OnClose(tomo.Stop)
	window.Show()
	go tick(label, clock)
}

func formatTime () (timeString string) {
	return time.Now().Format("2006-01-02 15:04:05")
}

func tick (label *elements.Label, clock *fun.AnalogClock) {
	for {
		tomo.Do (func () {
			label.SetText(formatTime())
			clock.SetTime(time.Now())
		})
		time.Sleep(time.Second)
	}
}
