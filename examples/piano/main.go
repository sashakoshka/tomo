package main

import "math"
import "errors"
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
var   waveform   = 0
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

	controlBar := basicElements.NewContainer(basicLayouts.Horizontal { true, false })
	label := basicElements.NewLabel("Play a song!", false)
	controlBar.Adopt(label, true)
	waveformButton := basicElements.NewButton("Sine")
	waveformButton.OnClick (func () {
		waveform = (waveform + 1) % 2
		switch waveform {
		case 0: waveformButton.SetText("Sine")
		case 1: waveformButton.SetText("Square")
		}
	})
	controlBar.Adopt(waveformButton, false)
	container.Adopt(controlBar, false)
	
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
	var streamer beep.Streamer
	switch waveform {
	case 0: streamer, _ = generators.SinTone(sampleRate, int(tuning.Tune(note)))
	case 1: streamer, _ = SquareTone(sampleRate, int(tuning.Tune(note)))
	}

	stopNote(note)
	speaker.Lock()
	playing[note] = &beep.Ctrl { Streamer: streamer }
	speaker.Unlock()
	speaker.Play(playing[note])
}

// https://github.com/faiface/beep/blob/v1.1.0/generators/toner.go
// Adapted to be a square wave instead

type toneStreamer struct {
	stat  float64
	delta float64
}

func SquareTone (sr beep.SampleRate, freq int) (beep.Streamer, error) {
	if int(sr) / freq < 2 {
		return nil, errors.New (
			"square tone generator: samplerate must be at least " +
			"2 times greater then frequency")
	}
	tone := new(toneStreamer)
	tone.stat = 0.0
	srf := float64(sr)
	ff := float64(freq)
	steps := srf / ff
	tone.delta = 1.0 / steps
	return tone, nil
}

func (tone *toneStreamer) nextSample () (sample float64) {
	if tone.stat > 0.5 {
		sample = 1
	} else {
		sample = 0
	}
	_, tone.stat = math.Modf(tone.stat + tone.delta)
	return
}

func (tone *toneStreamer) Stream (buf [][2]float64) (int, bool) {
	for i := 0; i < len(buf); i++ {
		s := tone.nextSample()
		buf[i] = [2]float64{s, s}
	}
	return len(buf), true
}
func (tone *toneStreamer) Err () error {
	return nil
}
