package fun

import "image"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/shatter"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"
import "git.tebibyte.media/sashakoshka/tomo/elements/fun/music"

const pianoKeyWidth = 18

type pianoKey struct {
	image.Rectangle
	music.Note
}

// Piano is an element that can be used to input midi notes.
type Piano struct {
	*core.Core
	*core.FocusableCore
	core core.CoreControl
	focusableControl core.FocusableCoreControl
	low, high music.Octave
	
	config config.Wrapped
	theme  theme.Wrapped

	flatKeys  []pianoKey
	sharpKeys []pianoKey
	contentBounds image.Rectangle

	pressed *pianoKey
	keynavPressed map[music.Note] bool

	onPress   func (music.Note)
	onRelease func (music.Note)
}

// NewPiano returns a new piano element with a lowest and highest octave,
// inclusive. If low is greater than high, they will be swapped.
func NewPiano (low, high music.Octave) (element *Piano) {
	if low > high {
		temp := low
		low = high
		high = temp
	}
	
	element = &Piano {
		low:  low,
		high: high,
		keynavPressed: make(map[music.Note] bool),
	}
	
	element.theme.Case = theme.C("fun", "piano")
	element.Core, element.core = core.NewCore (func () {
		element.recalculate()
		element.draw()
	})
	element.FocusableCore,
	element.focusableControl = core.NewFocusableCore(element.redo)
	element.updateMinimumSize()
	return
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
	element.redo()
}

func (element *Piano) HandleMouseMove (x, y int) {
	if element.pressed == nil { return }
	element.pressUnderMouseCursor(image.Pt(x, y))
}

func (element *Piano) HandleMouseScroll (x, y int, deltaX, deltaY float64) { }

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
		element.redo()
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
		element.redo()
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
	element.redo()
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
	inset := element.theme.Inset(theme.PatternSunken)
	element.core.SetMinimumSize (
		pianoKeyWidth * 7 * element.countOctaves() + inset[1] + inset[3],
		64 + inset[0] + inset[2])
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

	inset  := element.theme.Inset(theme.PatternPinboard)
	bounds := inset.Apply(element.Bounds())

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

func (element *Piano) draw () {
	state := theme.PatternState {
		Focused: element.Focused(),
		Disabled: !element.Enabled(),
	}

	for _, key := range element.flatKeys {
		_, keynavPressed := element.keynavPressed[key.Note]
		element.drawFlat (
			key.Rectangle,
			element.pressed != nil &&
			(*element.pressed).Note == key.Note || keynavPressed,
			state)
	}
	for _, key := range element.sharpKeys {
		_, keynavPressed := element.keynavPressed[key.Note]
		element.drawSharp (
			key.Rectangle,
			element.pressed != nil &&
			(*element.pressed).Note == key.Note || keynavPressed,
			state)
	}
	
	pattern := element.theme.Pattern(theme.PatternPinboard, state)
	tiles := shatter.Shatter(element.Bounds(), element.contentBounds)
	for _, tile := range tiles {
		artist.FillRectangleClip (
			element.core, pattern, element.Bounds(), tile)
	}
}

func (element *Piano) drawFlat (
	bounds image.Rectangle,
	pressed bool,
	state theme.PatternState,
) {
	state.Pressed = pressed
	pattern := element.theme.Theme.Pattern (
		theme.PatternButton, state, theme.C("fun", "flatKey"))
	artist.FillRectangle(element.core, pattern, bounds)
}

func (element *Piano) drawSharp (
	bounds image.Rectangle,
	pressed bool,
	state theme.PatternState,
) {
	state.Pressed = pressed
	pattern := element.theme.Theme.Pattern (
		theme.PatternButton, state, theme.C("fun", "sharpKey"))
	artist.FillRectangle(element.core, pattern, bounds)
}
