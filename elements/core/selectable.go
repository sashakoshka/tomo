package core

// import "runtime/debug"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/elements"

// FocusableCore is a struct that can be embedded into objects to make them
// focusable, giving them the default keynav behavior.
type FocusableCore struct {
	core CoreControl
	focused bool
	enabled bool
	drawFocusChange func ()
}

// NewFocusableCore creates a new focusability core and its corresponding
// control. If your element needs to visually update itself when it's focus
// state changes (which it should), a callback to draw and push the update can
// be specified. 
func NewFocusableCore (
	core CoreControl,
	drawFocusChange func (),
) (	
	focusable *FocusableCore,
	control FocusableCoreControl,
) {
	focusable = &FocusableCore {
		core: core,
		drawFocusChange: drawFocusChange,
		enabled: true,
	}
	control = FocusableCoreControl { core: focusable }
	return
}

// Focused returns whether or not this element is currently focused.
func (core *FocusableCore) Focused () (focused bool) {
	return core.focused
}

// Focus focuses this element, if its parent element grants the request.
func (core *FocusableCore) Focus () {
	if !core.enabled || core.focused { return }
	parent := core.core.Parent()
	if parent, ok := parent.(elements.FocusableParent); ok && parent != nil {
		core.focused = parent.RequestFocus (
			core.core.Outer().(elements.Focusable),
			input.KeynavDirectionNeutral)
	}
}

// HandleFocus causes this element to mark itself as focused, if it can
// currently be. Otherwise, it will return false and do nothing.
func (core *FocusableCore) HandleFocus (
	direction input.KeynavDirection,
) (
	accepted bool,
) {
	direction = direction.Canon()
	if !core.enabled { return false }
	if core.focused && direction != input.KeynavDirectionNeutral {
		return false
	}

	if core.focused == false {
		core.focused = true
		if core.drawFocusChange != nil { core.drawFocusChange() }
	}
	return true
}

// HandleUnfocus causes this element to mark itself as unfocused.
func (core *FocusableCore) HandleUnfocus () {
	core.focused = false
	// debug.PrintStack()
	if core.drawFocusChange != nil { core.drawFocusChange() }
}

// Enabled returns whether or not the element is enabled.
func (core *FocusableCore) Enabled () (enabled bool) {
	return core.enabled
}

// FocusableCoreControl is a struct that can be used to exert control over a
// focusability core. It must not be directly embedded into an element, but
// instead kept as a private member. When a FocusableCore struct is created, a
// corresponding FocusableCoreControl struct is linked to it and returned
// alongside it.
type FocusableCoreControl struct {
	core *FocusableCore
}

// SetEnabled sets whether the focusability core is enabled. If the state
// changes, this will call drawFocusChange.
func (control FocusableCoreControl) SetEnabled (enabled bool) {
	if control.core.enabled == enabled { return }
	control.core.enabled = enabled
	if !enabled { control.core.focused = false }
	if control.core.drawFocusChange != nil {
		control.core.drawFocusChange()
	}
}
