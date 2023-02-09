package main

import "math"
import "errors"
import "github.com/faiface/beep"
import "github.com/faiface/beep/speaker"
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
		waveform = (waveform + 1) % 5
		switch waveform {
		case 0: waveformButton.SetText("Sine")
		case 1: waveformButton.SetText("Square")
		case 2: waveformButton.SetText("Saw")
		case 3: waveformButton.SetText("Triangle")
		case 4: waveformButton.SetText("Supersaw")
		}
	})
	controlBar.Adopt(waveformButton, false)
	container.Adopt(controlBar, false)
	
	piano := fun.NewPiano(2, 5)
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
	streamer, _ := Tone(sampleRate, int(tuning.Tune(note)), waveform)

	stopNote(note)
	speaker.Lock()
	playing[note] = &beep.Ctrl { Streamer: streamer }
	speaker.Unlock()
	speaker.Play(playing[note])
}

// https://github.com/faiface/beep/blob/v1.1.0/generators/toner.go
// Adapted to be a bit more versatile.

type toneStreamer struct {
	position float64
	delta    float64
	waveform int
}

func Tone (sr beep.SampleRate, freq int, waveform int) (beep.Streamer, error) {
	if int(sr) / freq < 2 {
		return nil, errors.New (
			"square tone generator: samplerate must be at least " +
			"2 times greater then frequency")
	}
	tone := new(toneStreamer)
	tone.position = 0.0
	tone.waveform = waveform
	steps := float64(sr) / float64(freq)
	tone.delta = 1.0 / steps
	return tone, nil
}

func (tone *toneStreamer) nextSample () (sample float64) {
	switch tone.waveform {
	case 0:
		sample = math.Sin(tone.position * 2.0 * math.Pi)
	case 1:
		if tone.position > 0.5 {
			sample = 1
		} else {
			sample = -1
		}
	case 2:
		sample = (tone.position - 0.5) * 2
	case 3:
		sample = 1 - math.Abs(tone.position - 0.5) * 4
	case 4:
		sample =
			-1 + 13.7 * tone.position +
			28.32 * tone.position * tone.position +
			15.62 * tone.position * tone.position * tone.position
	}
	_, tone.position = math.Modf(tone.position + tone.delta)
	return
}

func (tone *toneStreamer) Stream (buf [][2]float64) (int, bool) {
	for i := 0; i < len(buf); i++ {
		sample := tone.nextSample()
		buf[i] = [2]float64{sample, sample}
	}
	return len(buf), true
}
func (tone *toneStreamer) Err () error {
	return nil
}
