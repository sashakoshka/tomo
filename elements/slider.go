package elements

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"
import "git.tebibyte.media/sashakoshka/tomo/default/config"

// Slider is a slider control with a floating point value between zero and one.
type Slider struct {
	slider
}

// NewSlider creates a new slider with the specified value.
func NewSlider (value float64, orientation Orientation) (element *Slider) {
	element = &Slider { }
	element.value = value
	element.vertical = bool(orientation)
	element.entity = tomo.NewEntity(element).(tomo.FocusableEntity)
	element.construct()
	return
}

type slider struct {
	entity tomo.FocusableEntity
	
	value      float64
	vertical   bool
	dragging   bool
	enabled    bool
	dragOffset int
	track      image.Rectangle
	bar        image.Rectangle
	
	config config.Wrapped
	theme  theme.Wrapped
	
	onSlide   func ()
	onRelease func ()
}

func (element *slider) construct () {
	element.enabled = true
	if element.vertical {
		element.theme.Case = tomo.C("tomo", "sliderVertical")
	} else {
		element.theme.Case = tomo.C("tomo", "sliderHorizontal")
	}
	element.updateMinimumSize()
}

// Entity returns this element's entity.
func (element *slider) Entity () tomo.Entity {
	return element.entity
}

// Draw causes the element to draw to the specified destination canvas.
func (element *slider) Draw (destination canvas.Canvas) {
	bounds := element.entity.Bounds()
	element.track = element.theme.Padding(tomo.PatternGutter).Apply(bounds)
	if element.vertical {
		barSize := element.track.Dx()
		element.bar = image.Rect(0, 0, barSize, barSize).Add(bounds.Min)
		barOffset :=
			float64(element.track.Dy() - barSize) *
			(1 - element.value)
		element.bar = element.bar.Add(image.Pt(0, int(barOffset)))
	} else {
		barSize := element.track.Dy()
		element.bar = image.Rect(0, 0, barSize, barSize).Add(bounds.Min)
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
	element.theme.Pattern(tomo.PatternGutter, state).Draw(destination, bounds)
	element.theme.Pattern(tomo.PatternHandle, state).Draw(destination, bounds)
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

func (element *slider) HandleMouseDown (x, y int, button input.Button) {
	if !element.Enabled() { return }
	element.Focus()
	if button == input.ButtonLeft {
		element.dragging = true
		element.value = element.valueFor(x, y)
		if element.onSlide != nil {
			element.onSlide()
		}
		element.entity.Invalidate()
	}
}

func (element *slider) HandleMouseUp (x, y int, button input.Button) {
	if button != input.ButtonLeft || !element.dragging { return }
	element.dragging = false
	if element.onRelease != nil {
		element.onRelease()
	}
	element.entity.Invalidate()
}

func (element *slider) HandleMotion (x, y int) {
	if element.dragging {
		element.dragging = true
		element.value = element.valueFor(x, y)
		if element.onSlide != nil {
			element.onSlide()
		}
		element.entity.Invalidate()
	}
}

func (element *slider) HandleScroll (x, y int, deltaX, deltaY float64) { }

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

// SetTheme sets the element's theme.
func (element *slider) SetTheme (new tomo.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.entity.Invalidate()
}

// SetConfig sets the element's configuration.
func (element *slider) SetConfig (new tomo.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
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
	gutterPadding := element.theme.Padding(tomo.PatternGutter)
	handlePadding := element.theme.Padding(tomo.PatternHandle)
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
