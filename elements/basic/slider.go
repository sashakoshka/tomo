package basicElements

import "image"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// Slider is a slider control with a floating point value between zero and one.
type Slider struct {
	*core.Core
	*core.FocusableCore
	core core.CoreControl
	focusableControl core.FocusableCoreControl

	value   float64
	vertical bool
	dragging bool
	dragOffset int
	track image.Rectangle
	bar image.Rectangle
	
	config config.Wrapped
	theme  theme.Wrapped
	
	onSlide   func ()
	onRelease func ()
}

// NewSlider creates a new slider with the specified value. If vertical is set
// to true, 
func NewSlider (value float64, vertical bool) (element *Slider) {
	element = &Slider {
		value: value,
		vertical: vertical,
	}
	if vertical {
		element.theme.Case = theme.C("basic", "sliderVertical")
	} else {
		element.theme.Case = theme.C("basic", "sliderHorizontal")
	}
	element.Core, element.core = core.NewCore(element, element.draw)
	element.FocusableCore,
	element.focusableControl = core.NewFocusableCore(element.core, element.redo)
	element.updateMinimumSize()
	return
}

func (element *Slider) HandleMouseDown (x, y int, button input.Button) {
	if !element.Enabled() { return }
	element.Focus()
	if button == input.ButtonLeft {
		element.dragging = true
		element.value = element.valueFor(x, y)
		if element.onSlide != nil {
			element.onSlide()
		}
		element.redo()
	}
}

func (element *Slider) HandleMouseUp (x, y int, button input.Button) {
	if button != input.ButtonLeft || !element.dragging { return }
	element.dragging = false
	if element.onRelease != nil {
		element.onRelease()
	}
	element.redo()
}

func (element *Slider) HandleMotion (x, y int) {
	if element.dragging {
		element.dragging = true
		element.value = element.valueFor(x, y)
		if element.onSlide != nil {
			element.onSlide()
		}
		element.redo()
	}
}

func (element *Slider) HandleScroll (x, y int, deltaX, deltaY float64) { }

func (element *Slider) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
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

func (element *Slider) HandleKeyUp (key input.Key, modifiers input.Modifiers) { }

// Value returns the slider's value.
func (element *Slider) Value () (value float64) {
	return element.value
}

// SetEnabled sets whether or not the slider can be interacted with.
func (element *Slider) SetEnabled (enabled bool) {
	element.focusableControl.SetEnabled(enabled)
}

// SetValue sets the slider's value.
func (element *Slider) SetValue (value float64) {
	if value < 0 { value = 0 }
	if value > 1 { value = 1 }
	
	if element.value == value { return }

	element.value = value
	if element.onRelease != nil {
		element.onRelease()
	}
	element.redo()
}

// OnSlide sets a function to be called every time the slider handle changes
// position while being dragged.
func (element *Slider) OnSlide (callback func ()) {
	element.onSlide = callback
}

// OnRelease sets a function to be called when the handle stops being dragged.
func (element *Slider) OnRelease (callback func ()) {
	element.onRelease = callback
}

// SetTheme sets the element's theme.
func (element *Slider) SetTheme (new theme.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.redo()
}

// SetConfig sets the element's configuration.
func (element *Slider) SetConfig (new config.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
	element.updateMinimumSize()
	element.redo()
}

func (element *Slider) changeValue (delta float64) {
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
	element.redo()
}

func (element *Slider) valueFor (x, y int) (value float64) {
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

func (element *Slider) updateMinimumSize () {
	if element.vertical {
		element.core.SetMinimumSize (
			element.config.HandleWidth(),
			element.config.HandleWidth() * 2)
	} else {
		element.core.SetMinimumSize (
			element.config.HandleWidth() * 2,
			element.config.HandleWidth())
	}
}

func (element *Slider) redo () {
	if element.core.HasImage () {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Slider) draw () {
	bounds := element.Bounds()
	element.track = element.theme.Padding(theme.PatternGutter).Apply(bounds)
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

	state := theme.State {
		Focused:  element.Focused(),
		Disabled: !element.Enabled(),
		Pressed:  element.dragging,
	}
	element.theme.Pattern(theme.PatternGutter, state).Draw (
		element.core,
		bounds)
	element.theme.Pattern(theme.PatternHandle, state).Draw (
		element.core,
		element.bar)
}
