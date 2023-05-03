package elements

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/artist"

// Slider is a slider control with a floating point value between zero and one.
type Slider struct {
	slider
}

// NewVSlider creates a new horizontal slider with the specified value.
func NewVSlider (value float64) (element *Slider) {
	element = NewHSlider(value)
	element.vertical = true
	return
}

// NewHSlider creates a new horizontal slider with the specified value.
func NewHSlider (value float64) (element *Slider) {
	element = &Slider { }
	element.value = value
	element.entity = tomo.GetBackend().NewEntity(element)
	element.construct()
	return
}

type slider struct {
	entity tomo.Entity

	c tomo.Case

	value      float64
	vertical   bool
	dragging   bool
	enabled    bool
	dragOffset int
	track      image.Rectangle
	bar        image.Rectangle
	
	onSlide   func ()
	onRelease func ()
}

func (element *slider) construct () {
	element.enabled = true
	if element.vertical {
		element.c = tomo.C("tomo", "sliderVertical")
	} else {
		element.c = tomo.C("tomo", "sliderHorizontal")
	}
	element.updateMinimumSize()
}

// Entity returns this element's entity.
func (element *slider) Entity () tomo.Entity {
	return element.entity
}

// Draw causes the element to draw to the specified destination canvas.
func (element *slider) Draw (destination artist.Canvas) {
	bounds := element.entity.Bounds()
	element.track = element.entity.Theme().Padding(tomo.PatternGutter, element.c).Apply(bounds)
	if element.vertical {
		barSize := element.track.Dx()
		element.bar = image.Rect(0, 0, barSize, barSize).Add(element.track.Min)
		barOffset :=
			float64(element.track.Dy() - barSize) *
			(1 - element.value)
		element.bar = element.bar.Add(image.Pt(0, int(barOffset)))
	} else {
		barSize := element.track.Dy()
		element.bar = image.Rect(0, 0, barSize, barSize).Add(element.track.Min)
		barOffset :=
			float64(element.track.Dx() - barSize) *
			element.value
		element.bar = element.bar.Add(image.Pt(int(barOffset), 0))
	}

	state := tomo.State {
		Disabled: !element.Enabled(),
		Focused:  element.entity.Focused(),
		Pressed:  element.dragging,
	}
	element.entity.Theme().Pattern(tomo.PatternGutter, state, element.c).Draw(destination, bounds)
	element.entity.Theme().Pattern(tomo.PatternHandle, state, element.c).Draw(destination, element.bar)
}

// Focus gives this element input focus.
func (element *slider) Focus () {
	if !element.entity.Focused() { element.entity.Focus() }
}

// Enabled returns whether this slider can be dragged or not.
func (element *slider) Enabled () bool {
	return element.enabled
}

// SetEnabled sets whether this slider can be dragged or not.
func (element *slider) SetEnabled (enabled bool) {
	if element.enabled == enabled { return }
	element.enabled = enabled
	element.entity.Invalidate()
}

func (element *slider) HandleFocusChange () {
	element.entity.Invalidate()
}

func (element *slider) HandleMouseDown  (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) {
	if !element.Enabled() { return }
	element.Focus()
	if button == input.ButtonLeft {
		element.dragging = true
		element.value = element.valueFor(position.X, position.Y)
		if element.onSlide != nil {
			element.onSlide()
		}
		element.entity.Invalidate()
	}
}

func (element *slider) HandleMouseUp  (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) {
	if button != input.ButtonLeft || !element.dragging { return }
	element.dragging = false
	if element.onRelease != nil {
		element.onRelease()
	}
	element.entity.Invalidate()
}

func (element *slider) HandleMotion (position image.Point) {
	if element.dragging {
		element.dragging = true
		element.value = element.valueFor(position.X, position.Y)
		if element.onSlide != nil {
			element.onSlide()
		}
		element.entity.Invalidate()
	}
}

func (element *slider) HandleScroll (
	position image.Point,
	deltaX, deltaY float64,
	modifiers input.Modifiers,
) { }

func (element *slider) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
	switch key {
	case input.KeyUp:
		element.changeValue(0.1)
	case input.KeyDown:
		element.changeValue(-0.1)
	case input.KeyRight:
		if element.vertical {
			element.changeValue(-0.1)
		} else {
			element.changeValue(0.1)
		}
	case input.KeyLeft:
		if element.vertical {
			element.changeValue(0.1)
		} else {
			element.changeValue(-0.1)
		}
	}
}

func (element *slider) HandleKeyUp (key input.Key, modifiers input.Modifiers) { }

// Value returns the slider's value.
func (element *slider) Value () (value float64) {
	return element.value
}

// SetValue sets the slider's value.
func (element *slider) SetValue (value float64) {
	if value < 0 { value = 0 }
	if value > 1 { value = 1 }
	
	if element.value == value { return }

	element.value = value
	if element.onRelease != nil {
		element.onRelease()
	}
	element.entity.Invalidate()
}

// OnSlide sets a function to be called every time the slider handle changes
// position while being dragged.
func (element *slider) OnSlide (callback func ()) {
	element.onSlide = callback
}

// OnRelease sets a function to be called when the handle stops being dragged.
func (element *slider) OnRelease (callback func ()) {
	element.onRelease = callback
}

func (element *slider) HandleThemeChange () {
	element.updateMinimumSize()
	element.entity.Invalidate()
}


func (element *slider) changeValue (delta float64) {
	element.value += delta
	if element.value < 0 {
		element.value = 0
	}
	if element.value > 1 {
		element.value = 1
	}
	if element.onRelease != nil {
		element.onRelease()
	}
	element.entity.Invalidate()
}

func (element *slider) valueFor (x, y int) (value float64) {
	if element.vertical {
		value =
			float64(y - element.track.Min.Y - element.bar.Dy() / 2) /
			float64(element.track.Dy() - element.bar.Dy())
		value = 1 - value
	} else {
		value =
			float64(x - element.track.Min.X - element.bar.Dx() / 2) /
			float64(element.track.Dx() - element.bar.Dx())
	}
	
	if value < 0 { value = 0 }
	if value > 1 { value = 1 }
	return
}

func (element *slider) updateMinimumSize () {
	gutterPadding := element.entity.Theme().Padding(tomo.PatternGutter, element.c)
	handlePadding := element.entity.Theme().Padding(tomo.PatternHandle, element.c)
	if element.vertical {
		element.entity.SetMinimumSize (
			gutterPadding.Horizontal() + handlePadding.Horizontal(),
			gutterPadding.Vertical()   + handlePadding.Vertical() * 2)
	} else {
		element.entity.SetMinimumSize (
			gutterPadding.Horizontal() + handlePadding.Horizontal() * 2,
			gutterPadding.Vertical()   + handlePadding.Vertical())
	}
}
