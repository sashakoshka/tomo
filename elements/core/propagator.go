package core

import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/elements"

// ChildIterator represents an object that can iterate over a list of children,
// calling a specified iterator function for each one. When keepGoing is false,
// the iterator stops the current loop and OverChildren returns.
type ChildIterator interface {
	OverChildren (func (child elements.Element) (keepGoing bool))
}

// Propagator is a struct that can be embedded into elements that contain one or
// more children in order to propagate events to them without having to write
// all of the event handlers. It also implements standard behavior for focus
// propagation and keyboard navigation.
type Propagator struct {
	iterator ChildIterator
	drags    [10]elements.MouseTarget
}

// NewPropagator creates a new event propagator that uses the specified iterator
// to access a list of child elements that will have events propagated to them.
// If iterator is nil, the function will return nil.
func NewPropagator (iterator ChildIterator) (propagator *Propagator) {
	if iterator == nil { return nil }
	propagator = &Propagator {
		iterator: iterator,
	}
	return
}

// Focused returns whether or not this element or any of its children
// are currently focused.
func (propagator *Propagator) Focused () (focused bool) {
	
}

// Focus focuses this element, if its parent element grants the
// request.
func (propagator *Propagator) Focus () {
	
}

// HandleFocus causes this element to mark itself as focused. If the
// element does not have children or there are no more focusable children in
// the given direction, it should return false and do nothing. Otherwise, it
// marks itself as focused along with any applicable children and returns
// true.
func (propagator *Propagator) HandleFocus (direction input.KeynavDirection) (accepted bool) {
	
}

// HandleDeselection causes this element to mark itself and all of its children
// as unfocused.
func (propagator *Propagator) HandleUnfocus () {
	
}

// OnFocusRequest sets a function to be called when this element wants its
// parent element to focus it. Parent elements should return true if the request
// was granted, and false if it was not. If the parent element returns true, the
// element acts as if a HandleFocus call was made with KeynavDirectionNeutral.
func (propagator *Propagator) OnFocusRequest (func () (granted bool)) {
	
}

// OnFocusMotionRequest sets a function to be called when this element wants its
// parent element to focus the element behind or in front of it, depending on
// the specified direction. Parent elements should return true if the request
// was granted, and false if it was not.
func (propagator *Propagator) OnFocusMotionRequest (func (direction input.KeynavDirection) (granted bool)) {
	
}

// HandleKeyDown propogates the keyboard event to the currently selected child.
func (propagator *Propagator) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
	
}

// HandleKeyUp propogates the keyboard event to the currently selected child.
func (propagator *Propagator) HandleKeyUp (key input.Key, modifiers input.Modifiers) {
	
}

// HandleMouseDown propagates the mouse event to the element under the mouse
// pointer.
func (propagator *Propagator) HandleMouseDown (x, y int, button input.Button) {
	
}

// HandleMouseUp propagates the mouse event to the element that the released
// mouse button was originally pressed on.
func (propagator *Propagator) HandleMouseUp (x, y int, button input.Button) {
	
}

// HandleMouseMove propagates the mouse event to the element that was last
// pressed down by the mouse if the mouse is currently being held down, else it
// propagates the event to whichever element is underneath the mouse pointer.
func (propagator *Propagator) HandleMouseMove (x, y int) {
	
}

// HandleScroll propagates the mouse event to the element under the mouse
// pointer.
func (propagator *Propagator) HandleMouseScroll (x, y int, deltaX, deltaY float64) {
	
}

// SetTheme sets the theme of all children to the specified theme.
func (propagator *Propagator) SetTheme (theme.Theme) {
	
}

// SetConfig sets the theme of all children to the specified config.
func (propagator *Propagator) SetConfig (config.Config) {
	
}
