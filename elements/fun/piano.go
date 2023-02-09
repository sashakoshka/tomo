package fun

import "image"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// Octave represents a MIDI octave.
type Octave int

// Note returns the note at the specified scale degree in the chromatic scale.
func (octave Octave) Note (degree int) Note {
	return Note(int(octave + 1) * 12 + degree)
}

// Note represents a MIDI note.
type Note int

// Octave returns the octave of the note
func (note Note) Octave () int {
	return int(note / 12 - 1)
}

// Degree returns the scale degree of the note in the chromatic scale.
func (note Note) Degree () int {
	mod := note % 12
	if mod < 0 { mod += 12 }
	return int(mod)
}

// IsSharp returns whether or not the note is a sharp.
func (note Note) IsSharp () bool {
	degree := note.Degree()
	return degree == 1 ||
		degree == 3 ||
		degree == 6 ||
		degree == 8 ||
		degree == 10
}

const pianoKeyWidth = 18

type pianoKey struct {
	image.Rectangle
	Note
}

type Piano struct {
	*core.Core
	core core.CoreControl
	low, high Octave
	
	config config.Wrapped
	theme  theme.Wrapped

	flatKeys  []pianoKey
	sharpKeys []pianoKey

	pressed *pianoKey

	onPress   func (Note)
	onRelease func (Note)
}

func NewPiano (low, high Octave) (element *Piano) {
	element = &Piano {
		low:  low,
		high: high,
	}
	element.theme.Case = theme.C("fun", "piano")
	element.Core, element.core = core.NewCore (func () {
		element.recalculate()
		element.draw()
	})
	element.updateMinimumSize()
	return
}

// OnPress sets a function to be called when a key is pressed.
func (element *Piano) OnPress (callback func (note Note)) {
	element.onPress = callback
}

// OnRelease sets a function to be called when a key is released.
func (element *Piano) OnRelease (callback func (note Note)) {
	element.onRelease = callback
}

func (element *Piano) HandleMouseDown (x, y int, button input.Button) {
	if button != input.ButtonLeft { return }
	element.pressUnderMouseCursor(image.Pt(x, y))
}

func (element *Piano) HandleMouseUp (x, y int, button input.Button) {
	if button != input.ButtonLeft { return }
	if element.onRelease != nil {
		element.onRelease((*element.pressed).Note)
	}
	element.pressed = nil
	element.redo()
}

func (element *Piano) HandleMouseMove (x, y int) {
	if element.pressed == nil { return }
	element.pressUnderMouseCursor(image.Pt(x, y))
}

func (element *Piano) HandleMouseScroll (x, y int, deltaX, deltaY float64) { }

func (element *Piano) pressUnderMouseCursor (point image.Point) {
	// release previous note
	if element.pressed != nil && element.onRelease != nil {
		element.onRelease((*element.pressed).Note)
	}

	// find out which note is being pressed
	newKey := (*pianoKey)(nil)
	for index, key := range element.flatKeys {
		if point.In(key.Rectangle) {
			newKey = &element.flatKeys[index]
			break
		}
	}
	for index, key := range element.sharpKeys {
		if point.In(key.Rectangle) {
			newKey = &element.sharpKeys[index]
			break
		}
	}
	if newKey == nil { return }
	
	if newKey != element.pressed {
		// press new note
		element.pressed = newKey
		if element.onPress != nil {
			element.onPress((*element.pressed).Note)
		}
		element.redo()
	}
}

// SetTheme sets the element's theme.
func (element *Piano) SetTheme (new theme.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.updateMinimumSize()
	element.recalculate()
	element.redo()
}

// SetConfig sets the element's configuration.
func (element *Piano) SetConfig (new config.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
	element.updateMinimumSize()
	element.recalculate()
	element.redo()
}

func (element *Piano) updateMinimumSize () {
	element.core.SetMinimumSize (
		pianoKeyWidth * 7 * element.countOctaves(), 64)
}

func (element *Piano) countOctaves () int {
	return int(element.high - element.low + 1)
}

func (element *Piano) countFlats () int {
	return element.countOctaves() * 8
}

func (element *Piano) countSharps () int {
	return element.countOctaves() * 5
}

func (element *Piano) redo () {
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Piano) recalculate () {
	element.flatKeys  = make([]pianoKey, element.countFlats())
	element.sharpKeys = make([]pianoKey, element.countSharps())

	bounds := element.Bounds()
	dot := bounds.Min
	note := element.low.Note(0)
	limit := element.high.Note(12)
	flatIndex := 0
	sharpIndex := 0
	for note < limit {
		if note.IsSharp() {
			element.sharpKeys[sharpIndex].Rectangle = image.Rect (
				-(pianoKeyWidth * 3) / 7, 0,
				(pianoKeyWidth * 3) / 7,
				bounds.Dy() / 2).Add(dot)
			element.sharpKeys[sharpIndex].Note = note
			sharpIndex ++
		} else {
			element.flatKeys[flatIndex].Rectangle = image.Rect (
				0, 0, pianoKeyWidth, bounds.Dy()).Add(dot)
			dot.X += pianoKeyWidth
			element.flatKeys[flatIndex].Note = note
			flatIndex ++
		}
		note ++
	}
}

func (element *Piano) draw () {
	for _, key := range element.flatKeys {
		element.drawFlat (
			key.Rectangle,
			element.pressed != nil &&
			(*element.pressed).Note == key.Note)
	}
	for _, key := range element.sharpKeys {
		element.drawSharp (
			key.Rectangle,
			element.pressed != nil &&
			(*element.pressed).Note == key.Note)
	}
}

func (element *Piano) drawFlat (bounds image.Rectangle, pressed bool) {
	state := theme.PatternState {
		Pressed: pressed,
	}
	pattern := element.theme.Pattern(theme.PatternButton, state)
	artist.FillRectangle(element, pattern, bounds)
}

func (element *Piano) drawSharp (bounds image.Rectangle, pressed bool) {
	state := theme.PatternState {
		Pressed: pressed,
	}
	pattern := element.theme.Pattern(theme.PatternButton, state)
	artist.FillRectangle(element, pattern, bounds)
}
