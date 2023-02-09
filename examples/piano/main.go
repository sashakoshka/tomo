package main

import "github.com/faiface/beep"
import "github.com/faiface/beep/speaker"
import "github.com/faiface/beep/generators"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements/fun"
import "git.tebibyte.media/sashakoshka/tomo/layouts/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/fun/music"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/x"

const sampleRate = 44100
const bufferSize = 256
var   tuning     = music.EqualTemparment { A4: 440 }

var playing = map[music.Note] *beep.Ctrl { }

func main () {
	speaker.Init(sampleRate, bufferSize)
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("Piano")
	container := basicElements.NewContainer(basicLayouts.Vertical { true, true })
	window.Adopt(container)

	label := basicElements.NewLabel("Play a song!", false)
	container.Adopt(label, false)
	piano := fun.NewPiano(3, 5)
	container.Adopt(piano, true)
	piano.OnPress(playNote)
	piano.OnRelease(stopNote)
	
	window.OnClose(tomo.Stop)
	window.Show()
}

func stopNote (note music.Note) {
	if _, is := playing[note]; !is { return }
	
	speaker.Lock() 
	playing[note].Streamer = nil
	delete(playing, note)
	speaker.Unlock()
}

func playNote (note music.Note) {
	streamer, err := generators.SinTone(sampleRate, int(tuning.Tune(note)))
	if err != nil { panic(err.Error()) }

	stopNote(note)
	speaker.Lock()
	playing[note] = &beep.Ctrl { Streamer: streamer }
	speaker.Unlock()
	speaker.Play(playing[note])
}
