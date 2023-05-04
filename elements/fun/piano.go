package fun

import "image"
import "tomo"
import "tomo/input"
import "art"
import "art/artutil"
import "tomo/elements/fun/music"

var pianoCase = tomo.C("tomo", "piano")
var flatCase  = tomo.C("tomo", "piano", "flatKey")
var sharpCase = tomo.C("tomo", "piano", "sharpKey")

const pianoKeyWidth = 18

type pianoKey struct {
	image.Rectangle
	music.Note
}

// Piano is an element that can be used to input midi notes.
type Piano struct {
	entity tomo.Entity

	low, high music.Octave
	flatKeys  []pianoKey
	sharpKeys []pianoKey
	contentBounds image.Rectangle

	enabled bool
	pressed *pianoKey
	keynavPressed map[music.Note] bool

	onPress   func (music.Note)
	onRelease func (music.Note)
}

// NewPiano returns a new piano element with a lowest and highest octave,
// inclusive. If low is greater than high, they will be swapped.
func NewPiano (low, high music.Octave) (element *Piano) {
	if low > high { low, high = high, low }
	
	element = &Piano {
		low:  low,
		high: high,
		keynavPressed: make(map[music.Note] bool),
	}
	
	element.entity = tomo.GetBackend().NewEntity(element)
	element.updateMinimumSize()
	return
}

// Entity returns this element's entity.
func (element *Piano) Entity () tomo.Entity {
	return element.entity
}

// Draw causes the element to draw to the specified destination canvas.
func (element *Piano) Draw (destination art.Canvas) {
	element.recalculate()

	state := tomo.State {
		Focused:  element.entity.Focused(),
		Disabled: !element.Enabled(),
	}

	for _, key := range element.flatKeys {
		_, keynavPressed := element.keynavPressed[key.Note]
		element.drawFlat (
			destination,
			key.Rectangle,
			element.pressed != nil &&
			(*element.pressed).Note == key.Note || keynavPressed,
			state)
	}
	for _, key := range element.sharpKeys {
		_, keynavPressed := element.keynavPressed[key.Note]
		element.drawSharp (
			destination,
			key.Rectangle,
			element.pressed != nil &&
			(*element.pressed).Note == key.Note || keynavPressed,
			state)
	}
	
	pattern := element.entity.Theme().Pattern(tomo.PatternPinboard, state, pianoCase)
	artutil.DrawShatter (
		destination, pattern, element.entity.Bounds(),
		element.contentBounds)
}

// Focus gives this element input focus.
func (element *Piano) Focus () {
	element.entity.Focus()
}

// Enabled returns whether this piano can be played or not.
func (element *Piano) Enabled () bool {
	return element.enabled
}

// SetEnabled sets whether this piano can be played or not.
func (element *Piano) SetEnabled (enabled bool) {
	if element.enabled == enabled { return }
	element.enabled = enabled
	element.entity.Invalidate()
}


// OnPress sets a function to be called when a key is pressed.
func (element *Piano) OnPress (callback func (note music.Note)) {
	element.onPress = callback
}

// OnRelease sets a function to be called when a key is released.
func (element *Piano) OnRelease (callback func (note music.Note)) {
	element.onRelease = callback
}

func (element *Piano) HandleMouseDown (x, y int, button input.Button) {
	element.Focus()
	if button != input.ButtonLeft { return }
	element.pressUnderMouseCursor(image.Pt(x, y))
}

func (element *Piano) HandleMouseUp (x, y int, button input.Button) {
	if button != input.ButtonLeft { return }
	if element.onRelease != nil && element.pressed != nil {
		element.onRelease((*element.pressed).Note)
	}
	element.pressed = nil
	element.entity.Invalidate()
}

func (element *Piano) HandleMotion (x, y int) {
	if element.pressed == nil { return }
	element.pressUnderMouseCursor(image.Pt(x, y))
}

func (element *Piano) pressUnderMouseCursor (point image.Point) {
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
		// release previous note
		if element.pressed != nil && element.onRelease != nil {
			element.onRelease((*element.pressed).Note)
		}
	
		// press new note
		element.pressed = newKey
		if element.onPress != nil {
			element.onPress((*element.pressed).Note)
		}
		element.entity.Invalidate()
	}
}

var noteForKey = map[input.Key] music.Note {
	'a': 46,
	'z': 47,

	'x': 48,
	'd': 49,
	'c': 50,
	'f': 51,
	'v': 52,
	'b': 53,
	'h': 54,
	'n': 55,
	'j': 56,
	'm': 57,
	'k': 58,
	',': 59,
	'.': 60,
	';': 61,
	'/': 62,
	'\'': 63,

	'1': 56,
	'q': 57,
	'2': 58,
	'w': 59,
	
	'e': 60,
	'4': 61,
	'r': 62,
	'5': 63,
	't': 64,
	'y': 65,
	'7': 66,
	'u': 67,
	'8': 68,
	'i': 69,
	'9': 70,
	'o': 71,
	
	'p': 72,
	'-': 73,
	'[': 74,
	'=': 75,
	']': 76,
	'\\': 77,
}

func (element *Piano) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
	if !element.Enabled() { return }
	note, exists := noteForKey[key]
	if !exists { return }
	if !element.keynavPressed[note] {
		element.keynavPressed[note] = true
		if element.onPress != nil {
			element.onPress(note)
		}
		element.entity.Invalidate()
	}
}

func (element *Piano) HandleKeyUp (key input.Key, modifiers input.Modifiers) {
	note, exists := noteForKey[key]
	if !exists { return }
	_, pressed := element.keynavPressed[note]
	if !pressed { return }
	delete(element.keynavPressed, note)
	if element.onRelease != nil {
		element.onRelease(note)
	}
	element.entity.Invalidate()
}

func (element *Piano) HandleThemeChange () {
	element.updateMinimumSize()
	element.entity.Invalidate()
}

func (element *Piano) updateMinimumSize () {
	padding := element.entity.Theme().Padding(tomo.PatternPinboard, pianoCase)
	element.entity.SetMinimumSize (
		pianoKeyWidth * 7 * element.countOctaves() +
		padding.Horizontal(),
		64 + padding.Vertical())
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

func (element *Piano) recalculate () {
	element.flatKeys  = make([]pianoKey, element.countFlats())
	element.sharpKeys = make([]pianoKey, element.countSharps())

	padding := element.entity.Theme().Padding(tomo.PatternPinboard, pianoCase)
	bounds  := padding.Apply(element.entity.Bounds())

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
				(bounds.Dy() * 5) / 8).Add(dot)
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

	element.contentBounds = image.Rectangle {
		bounds.Min,
		image.Pt(dot.X, bounds.Max.Y),
	}
}

func (element *Piano) drawFlat (
	destination art.Canvas,
	bounds image.Rectangle,
	pressed bool,
	state tomo.State,
) {
	state.Pressed = pressed
	pattern := element.entity.Theme().Pattern(tomo.PatternButton, state, flatCase)
	pattern.Draw(destination, bounds)
}

func (element *Piano) drawSharp (
	destination art.Canvas,
	bounds image.Rectangle,
	pressed bool,
	state tomo.State,
) {
	state.Pressed = pressed
	pattern := element.entity.Theme().Pattern(tomo.PatternButton, state, sharpCase)
	pattern.Draw(destination, bounds)
}
