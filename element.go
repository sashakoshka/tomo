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

	// SelectabilityChange is called when the chid element becomes
	// selectable or non-selectable.
	SelectabilityChange func (selectable bool)

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

// RunSelectabilityChange runs the SelectionRequest hook if it is not nil. If it
// is nil, it does nothing.
func (hooks ParentHooks) RunSelectabilityChange (selectable bool) {
	if hooks.SelectabilityChange != nil {
		hooks.SelectabilityChange(selectable)
	}
}

// Element represents a basic on-screen object.
type Element interface {
	// Element must implement the Canvas interface. Elements should start
	// out with a completely blank buffer, and only allocate memory and draw
	// on it for the first time when sent an EventResize event.
	Canvas

	// Handle handles an event, propagating it to children if necessary.
	Handle (event Event)

	// Selectable returns whether this element can be selected. If this
	// element contains other selectable elements, it must return true.
	Selectable () (selectable bool)

	// Selected returns whether or not this element is currently selected.
	// This will always return false if it is not selectable.
	Selected () (selected bool)

	// If this element contains other elements, and one is selected, this
	// method will advance the selection in the specified direction. If
	// the element contains selectable elements but none of them are
	// selected, it will select the first selectable element. If there are
	// no more children to be selected in the specified direction, the
	// element will return false. If the selection could be advanced, it
	// will return true. If the element contains no selectable child
	// elements, it will  always return false.
	AdvanceSelection (direction int) (ok bool)

	// SetParentHooks gives the element callbacks that let it send
	// information to its parent element without it knowing anything about
	// the parent element or containing any reference to it. When a parent
	// element adopts a child element, it must set these callbacks.
	SetParentHooks (callbacks ParentHooks)

	// MinimumSize specifies the minimum amount of pixels this element's
	// width and height may be set to. If the element is given a resize
	// event with dimensions smaller than this, it will use its minimum
	// instead of the offending dimension(s).
	MinimumSize () (width, height int)
}
