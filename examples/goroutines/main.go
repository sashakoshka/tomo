package main

import "os"
import "time"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import "git.tebibyte.media/sashakoshka/tomo/elements/fun"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
	os.Exit(0)
}

func run () {
	window, _ := tomo.NewWindow(tomo.Bounds(0, 0, 200, 216))
	window.SetTitle("Clock")
	container := elements.NewVBox(elements.SpaceBoth)
	window.Adopt(container)

	clock := fun.NewAnalogClock(time.Now())
	label := elements.NewLabel(formatTime())
	container.AdoptExpand(clock)
	container.Adopt(label)
	
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
