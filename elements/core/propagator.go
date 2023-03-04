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
func NewPropagator (iterator ChildIterator) (propagator *Propagator) {
	propagator = &Propagator {
		iterator: iterator,
	}
	return
}

// Focused returns whether or not this element or any of its children
// are currently focused.
func (propagator *Propagator) Focused () (focused bool)

// Focus focuses this element, if its parent element grants the
// request.
func (propagator *Propagator) Focus ()

// HandleFocus causes this element to mark itself as focused. If the
// element does not have children, it is disabled, or there are no more
// selectable children in the given direction, it should return false
// and do nothing. Otherwise, it should select itself and any children
// (if applicable) and return true.
func (propagator *Propagator) HandleFocus (direction input.KeynavDirection) (accepted bool)

// HandleDeselection causes this element to mark itself and all of its
// children as unfocused.
func (propagator *Propagator) HandleUnfocus ()

// OnFocusRequest sets a function to be called when this element wants
// its parent element to focus it. Parent elements should return true if
// the request was granted, and false if it was not. If the parent
// element returns true, the element must act as if a HandleFocus call
// was made with KeynavDirectionNeutral.
func (propagator *Propagator) OnFocusRequest (func () (granted bool))

// OnFocusMotionRequest sets a function to be called when this
// element wants its parent element to focus the element behind or in
// front of it, depending on the specified direction. Parent elements
// should return true if the request was granted, and false if it was
// not.
func (propagator *Propagator) OnFocusMotionRequest (func (direction input.KeynavDirection) (granted bool))

// HandleKeyDown is called when a key is pressed down or repeated while
// this element has keyboard focus. It is important to note that not
// every key down event is guaranteed to be paired with exactly one key
// up event. This is the reason a list of modifier keys held down at the
// time of the key press is given.
func (propagator *Propagator) HandleKeyDown (key input.Key, modifiers input.Modifiers)

// HandleKeyUp is called when a key is released while this element has
// keyboard focus.
func (propagator *Propagator) HandleKeyUp (key input.Key, modifiers input.Modifiers)

// HandleMouseDown is called when a mouse button is pressed down on this
// element.
func (propagator *Propagator) HandleMouseDown (x, y int, button input.Button)

// HandleMouseUp is called when a mouse button is released that was
// originally pressed down on this element.
func (propagator *Propagator) HandleMouseUp (x, y int, button input.Button)

// HandleMouseMove is called when the mouse is moved over this element,
// or the mouse is moving while being held down and originally pressed
// down on this element.
func (propagator *Propagator) HandleMouseMove (x, y int)

// HandleScroll is called when the mouse is scrolled. The X and Y
// direction of the scroll event are passed as deltaX and deltaY.
func (propagator *Propagator) HandleMouseScroll (x, y int, deltaX, deltaY float64)

// SetTheme sets the element's theme to something fulfilling the
// theme.Theme interface.
func (propagator *Propagator) SetTheme (theme.Theme)

// SetConfig sets the element's configuration to something fulfilling
// the config.Config interface.
func (propagator *Propagator) SetConfig (config.Config)
