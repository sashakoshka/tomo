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
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

const sampleRate = 44100
const bufferSize = 256
var   tuning     = music.EqualTemparment { A4: 440 }
var   waveform   = 0
var   playing    = map[music.Note] *toneStreamer { }
var   adsr = ADSR {
	Attack:  5 * time.Millisecond,
	Decay:   400 * time.Millisecond,
	Sustain: 0.7,
	Release: 500 * time.Millisecond,
}
var gain = 0.3

func main () {
	speaker.Init(sampleRate, bufferSize)
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("Piano")
	container := containers.NewContainer(basicLayouts.Vertical { true, true })
	controlBar := containers.NewContainer(basicLayouts.Horizontal { true, false })

	waveformColumn := containers.NewContainer(basicLayouts.Vertical { true, false })
	waveformList := basicElements.NewList (
		basicElements.NewListEntry("Sine",     func(){ waveform = 0 }),
		basicElements.NewListEntry("Triangle", func(){ waveform = 3 }),
		basicElements.NewListEntry("Square",   func(){ waveform = 1 }),
		basicElements.NewListEntry("Saw",      func(){ waveform = 2 }),
		basicElements.NewListEntry("Supersaw", func(){ waveform = 4 }),
	)
	waveformList.OnNoEntrySelected (func(){waveformList.Select(0)})
	waveformList.Select(0)

	adsrColumn := containers.NewContainer(basicLayouts.Vertical { true, false })
	adsrGroup := containers.NewContainer(basicLayouts.Horizontal { true, false })
	attackSlider  := basicElements.NewLerpSlider(0, 3 * time.Second, adsr.Attack, true)
	decaySlider   := basicElements.NewLerpSlider(0, 3 * time.Second, adsr.Decay, true)
	sustainSlider := basicElements.NewSlider(adsr.Sustain, true)
	releaseSlider := basicElements.NewLerpSlider(0, 3 * time.Second, adsr.Release, true)
	gainSlider    := basicElements.NewSlider(math.Sqrt(gain), false)

	attackSlider.OnRelease (func () {
		adsr.Attack = attackSlider.Value()
	})
	decaySlider.OnRelease (func () {
		adsr.Decay = decaySlider.Value()
	})
	sustainSlider.OnRelease (func () {
		adsr.Sustain = sustainSlider.Value()
	})
	releaseSlider.OnRelease (func () {
		adsr.Release = releaseSlider.Value()
	})
	gainSlider.OnRelease (func () {
		gain = math.Pow(gainSlider.Value(), 2)
	})

	patchColumn := containers.NewContainer(basicLayouts.Vertical { true, false })
	patch := func (w int, a, d time.Duration, s float64, r time.Duration) func () {
		return func () {
			waveform = w
			adsr     = ADSR {
				a * time.Millisecond,
				d * time.Millisecond,
				s,
				r * time.Millisecond,
			}
			waveformList.Select(w)
			attackSlider .SetValue(adsr.Attack)
			decaySlider  .SetValue(adsr.Decay)
			sustainSlider.SetValue(adsr.Sustain)
			releaseSlider.SetValue(adsr.Release)
		}
	}
	patchList := basicElements.NewList (
		basicElements.NewListEntry ("Bones", patch (
			0, 0, 100, 0.0, 0)),
		basicElements.NewListEntry ("Staccato", patch (
			4, 70, 500, 0, 0)),
		basicElements.NewListEntry ("Sustain", patch (
			4, 70, 200, 0.8, 500)),
		basicElements.NewListEntry ("Upright", patch (
			1, 0, 500, 0.4, 70)),
		basicElements.NewListEntry ("Space Pad", patch (
			4, 1500, 0, 1.0, 3000)),
		basicElements.NewListEntry ("Popcorn", patch (
			2, 0, 40, 0.0, 0)),
		basicElements.NewListEntry ("Racer", patch (
			3, 70, 0, 0.7, 400)),
		basicElements.NewListEntry ("Reverse", patch (
			2, 3000, 60, 0, 0)),
	)
	patchList.Collapse(0, 32)
	patchScrollBox := containers.NewScrollContainer(false, true)
	
	piano := fun.NewPiano(2, 5)
	piano.OnPress(playNote)
	piano.OnRelease(stopNote)

	// honestly, if you were doing something like this for real, i'd
	// encourage you to build a custom layout because this is a bit cursed.
	// i need to add more layouts...

	window.Adopt(container)
	
	controlBar.Adopt(patchColumn, true)
	patchColumn.Adopt(basicElements.NewLabel("Presets", false), false)
	patchColumn.Adopt(patchScrollBox, true)
	patchScrollBox.Adopt(patchList)

	controlBar.Adopt(basicElements.NewSpacer(true), false)
	
	controlBar.Adopt(waveformColumn, false)
	waveformColumn.Adopt(basicElements.NewLabel("Waveform", false), false)
	waveformColumn.Adopt(waveformList, true)
	
	controlBar.Adopt(basicElements.NewSpacer(true), false)
	
	adsrColumn.Adopt(basicElements.NewLabel("ADSR", false), false)
	adsrGroup.Adopt(attackSlider, false)
	adsrGroup.Adopt(decaySlider, false)
	adsrGroup.Adopt(sustainSlider, false)
	adsrGroup.Adopt(releaseSlider, false)
	adsrColumn.Adopt(adsrGroup, true)
	adsrColumn.Adopt(gainSlider, false)
	
	controlBar.Adopt(adsrColumn, false)
	container.Adopt(controlBar, true)
	container.Adopt(piano, false)
	
	piano.Focus()
	window.OnClose(tomo.Stop)
	window.Show()
}

type Patch struct {
	ADSR
	Waveform int
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
		gain,
		adsr)

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
	cycles   uint64
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
		unison := 5
		detuneDelta := 0.00005
		
		detune := 0.0 - (float64(unison) / 2) * detuneDelta
		for i := 0; i < unison; i ++ {
			_, offset := math.Modf(detune * float64(tone.cycles) + tone.position)
			sample += (offset - 0.5) * 2
			detune += detuneDelta
		}

		sample /= float64(unison)
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
	tone.cycles ++
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
