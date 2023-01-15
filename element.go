package tomo

// ParentHooks is a struct that contains callbacks that let child elements send
// information to their parent element without the child element knowing
// anything about the parent element or containing any reference to it. When a
// parent element adopts a child element, it must set these callbacks.
type ParentHooks struct {
	// Draw is called when a part of the child element's surface is updated.
	// The updated region will be passed to the callback as a sub-image.
	Draw func (region Canvas)

	// MinimumSizeChange is called when the child element's minimum width
	// and/or height changes. When this function is called, the element will
	// have already been resized and there is no need to send it a resize
	// event.
	MinimumSizeChange func (width, height int)
	
	// SelectionRequest is called when the child element element wants
	// itself to be selected. If the parent element chooses to grant the
	// request, it must send the child element a selection event and return
	// true.
	SelectionRequest func () (granted bool)
}

// RunDraw runs the Draw hook if it is not nil. If it is nil, it does nothing.
func (hooks ParentHooks) RunDraw (region Canvas) {
	if hooks.Draw != nil {
		hooks.Draw(region)
	}
}

// RunMinimumSizeChange runs the MinimumSizeChange hook if it is not nil. If it
// is nil, it does nothing.
func (hooks ParentHooks) RunMinimumSizeChange (width, height int) {
	if hooks.MinimumSizeChange != nil {
		hooks.MinimumSizeChange(width, height)
	}
}

// RunSelectionRequest runs the SelectionRequest hook if it is not nil. If it is
// nil, it does nothing.
func (hooks ParentHooks) RunSelectionRequest () (granted bool) {
	if hooks.SelectionRequest != nil {
		granted = hooks.SelectionRequest()
	}
	return
}

// Element represents a basic on-screen object.
type Element interface {
	// Element must implement the Canvas interface. Elements should start
	// out with a completely blank buffer, and only allocate memory and draw
	// on it for the first time when sent an EventResize event.
	Canvas

	// MinimumSize specifies the minimum amount of pixels this element's
	// width and height may be set to. If the element is given a resize
	// event with dimensions smaller than this, it will use its minimum
	// instead of the offending dimension(s).
	MinimumSize () (width, height int)

	// Resize resizes the element. This should only be called by the
	// element's parent.
	Resize (width, height int)

	// SetParentHooks gives the element callbacks that let it send
	// information to its parent element without it knowing anything about
	// the parent element or containing any reference to it. When a parent
	// element adopts a child element, it must set these callbacks.
	SetParentHooks (callbacks ParentHooks)
}

// SelectionDirection represents a keyboard navigation direction.
type SelectionDirection int

const (
	SelectionDirectionNeutral  SelectionDirection =  0
	SelectionDirectionBackward SelectionDirection = -1
	SelectionDirectionForward  SelectionDirection =  1
)

// Selectable represents an element that has keyboard navigation support. This
// includes inputs, buttons, sliders, etc. as well as any elements that have
// children (so keyboard navigation events can be propagated downward).
type Selectable interface {
	Element

	// Selected returns whether or not this element is currently selected.
	Selected () (selected bool)

	// Select selects this element, if its parent element grants the
	// request.
	Select ()

	// HandleSelection causes this element to mark itself as selected, if it
	// can currently be. Otherwise, it will return false and do nothing.
	HandleSelection (direction SelectionDirection) (accepted bool)

	// HandleDeselection causes this element to mark itself and all of its
	// children as deselected.
	HandleDeselection ()
}

// KeyboardTarget represents an element that can receive keyboard input.
type KeyboardTarget interface {
	Element

	// HandleKeyDown is called when a key is pressed down while this element
	// has keyboard focus. It is important to note that not every key down
	// event is guaranteed to be paired with exactly one key up event. This
	// is the reason a list of modifier keys held down at the time of the
	// key press is given.
	HandleKeyDown (key Key, modifiers Modifiers, repeated bool)

	// HandleKeyUp is called when a key is released while this element has
	// keyboard focus.
	HandleKeyUp (key Key, modifiers Modifiers)
}

// MouseTarget represents an element that can receive mouse events.
type MouseTarget interface {
	Element

	// Each of these handler methods is passed the position of the mouse
	// cursor at the time of the event as x, y.

	// HandleMouseDown is called when a mouse button is pressed down on this
	// element.
	HandleMouseDown (x, y int, button Button)

	// HandleMouseUp is called when a mouse button is released that was
	// originally pressed down on this element.
	HandleMouseUp (x, y int, button Button)

	// HandleMouseMove is called when the mouse is moved over this element,
	// or the mouse is moving while being held down and originally pressed
	// down on this element.
	HandleMouseMove (x, y int)

	// HandleScroll is called when the mouse is scrolled. The X and Y
	// direction of the scroll event are passed as deltaX and deltaY.
	HandleScroll (x, y int, deltaX, deltaY float64)
}

// Expanding represents an element who's minimum height can change in response
// to a change in its width.
type Expanding interface {
	Element

	// HeightForWidth returns what the element's minimum height would be if
	// resized to the specified width. This does not actually alter the
	// state of the element in any way, but it may perform significant work,
	// so it should be used sparingly.
	MinimumHeightFor (width int) (height int)
}
