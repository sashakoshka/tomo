package main

import "time"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/nasin"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import "git.tebibyte.media/sashakoshka/tomo/elements/fun"

func main () {
	nasin.Run(Application { })
}

type Application struct { }

func (Application) Init () error {
	window, err := nasin.NewWindow(tomo.Bounds(0, 0, 200, 216))
	if err != nil { return err }
	window.SetTitle("Clock")
	window.SetApplicationName("TomoClock")
	container := elements.NewVBox(elements.SpaceBoth)
	window.Adopt(container)

	clock := fun.NewAnalogClock(time.Now())
	label := elements.NewLabel(formatTime())
	container.AdoptExpand(clock)
	container.Adopt(label)
	
	window.OnClose(nasin.Stop)
	window.Show()
	go tick(label, clock)
	return nil
}

func formatTime () (timeString string) {
	return time.Now().Format("2006-01-02 15:04:05")
}

func tick (label *elements.Label, clock *fun.AnalogClock) {
	for {
		nasin.Do (func () {
			label.SetText(formatTime())
			clock.SetTime(time.Now())
		})
		time.Sleep(time.Second)
	}
}
