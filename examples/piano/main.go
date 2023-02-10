package main

import "math"
import "time"
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
var playing = map[music.Note] *toneStreamer { }

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
	piano.Focus()
	
	window.OnClose(tomo.Stop)
	window.Show()
}

func stopNote (note music.Note) {
	if _, is := playing[note]; !is { return }
	
	speaker.Lock() 
	playing[note].Release()
	delete(playing, note)
	speaker.Unlock()
}

func playNote (note music.Note) {
	streamer, _ := Tone (
		sampleRate,
		int(tuning.Tune(note)),
		waveform,
		0.3,
		ADSR {
			Attack:  100 * time.Millisecond,
			Decay:   400 * time.Millisecond,
			Sustain: 0.7,
			Release: 500 * time.Millisecond,
		})

	stopNote(note)
	speaker.Lock()
	playing[note] = streamer
	speaker.Unlock()
	speaker.Play(playing[note])
}

// https://github.com/faiface/beep/blob/v1.1.0/generators/toner.go
// Adapted to be a bit more versatile.

type toneStreamer struct {
	position float64
	delta    float64
	
	waveform int
	gain     float64

	adsr     ADSR
	released bool
	complete bool

	adsrPhase    int
	adsrPosition float64
	adsrDeltas   [4]float64
}

type ADSR struct {
	Attack  time.Duration
	Decay   time.Duration
	Sustain float64
	Release time.Duration
}

func Tone (
	sampleRate beep.SampleRate,
	frequency int,
	waveform int,
	gain float64,
	adsr ADSR,
) (
	*toneStreamer,
	error,
) {
	if int(sampleRate) / frequency < 2 {
		return nil, errors.New (
			"tone generator: samplerate must be at least " +
			"2 times greater then frequency")
	}
	
	tone := new(toneStreamer)
	tone.waveform = waveform
	tone.position = 0.0
	steps := float64(sampleRate) / float64(frequency)
	tone.delta = 1.0 / steps
	tone.gain = gain

	if adsr.Attack  < time.Millisecond { adsr.Attack  = time.Millisecond }
	if adsr.Decay   < time.Millisecond { adsr.Decay   = time.Millisecond }
	if adsr.Release < time.Millisecond { adsr.Release = time.Millisecond }
	tone.adsr = adsr

	attackSteps  := adsr.Attack.Seconds()  * float64(sampleRate)
	decaySteps   := adsr.Decay.Seconds()   * float64(sampleRate)
	releaseSteps := adsr.Release.Seconds() * float64(sampleRate)
	tone.adsrDeltas[0] = 1 / attackSteps
	tone.adsrDeltas[1] = 1 / decaySteps
	tone.adsrDeltas[2] = 0
	tone.adsrDeltas[3] = 1 / releaseSteps
	
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

	adsrGain := 0.0
	switch tone.adsrPhase {
	case 0: adsrGain = tone.adsrPosition
		if tone.adsrPosition > 1 {
			tone.adsrPosition = 0
			tone.adsrPhase = 1
		}
		
	case 1: adsrGain = 1 + tone.adsrPosition * (tone.adsr.Sustain - 1)
		if tone.adsrPosition > 1 {
			tone.adsrPosition = 0
			tone.adsrPhase = 2
		}
		
	case 2: adsrGain = tone.adsr.Sustain
		if tone.released {
			tone.adsrPhase = 3
		}
		
	case 3: adsrGain = (1 - tone.adsrPosition) * tone.adsr.Sustain
		if tone.adsrPosition > 1 {
			tone.adsrPosition = 0
			tone.complete = true
		}
	}

	sample *= adsrGain * adsrGain
	
	tone.adsrPosition += tone.adsrDeltas[tone.adsrPhase]
	_, tone.position = math.Modf(tone.position + tone.delta)
	return
}

func (tone *toneStreamer) Stream (buf [][2]float64) (int, bool) {
	if tone.complete {
		return 0, false
	}

	for i := 0; i < len(buf); i++ {
		sample := 0.0
		if !tone.complete {
			sample = tone.nextSample() * tone.gain
		}
		buf[i] = [2]float64{sample, sample}
	}
	return len(buf), true
}

func (tone *toneStreamer) Err () error {
	return nil
}

func (tone *toneStreamer) Release () {
	tone.released = true
}
