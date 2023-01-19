package tomo

import "image"

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

	// OnDamage sets a function to be called when an area of the element is
	// drawn on and should be pushed to the screen.
	OnDamage (callback func (region Canvas))

	// OnMinimumSizeChange sets a function to be called when the element's
	// minimum size is changed.
	OnMinimumSizeChange (callback func ())
}

// SelectionDirection represents a keyboard navigation direction.
type SelectionDirection int

const (
	SelectionDirectionNeutral  SelectionDirection =  0
	SelectionDirectionBackward SelectionDirection = -1
	SelectionDirectionForward  SelectionDirection =  1
)

// Canon returns a well-formed direction.
func (direction SelectionDirection) Canon () (canon SelectionDirection) {
	if direction > 0 {
		return SelectionDirectionForward
	} else if direction == 0 {
		return SelectionDirectionNeutral
	} else {
		return SelectionDirectionBackward
	}
}

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

	// OnSelectionRequest sets a function to be called when this element
	// wants its parent element to select it. Parent elements should return
	// true if the request was granted, and false if it was not.
	OnSelectionRequest (func () (granted bool))

	// OnSelectionMotionRequest sets a function to be called when this
	// element wants its parent element to select the element behind or in
	// front of it, depending on the specified direction. Parent elements
	// should return true if the request was granted, and false if it was
	// not.
	OnSelectionMotionRequest (func (SelectionDirection) (granted bool))
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
	HandleMouseScroll (x, y int, deltaX, deltaY float64)
}

// Flexible represents an element who's preferred minimum height can change in
// response to its width.
type Flexible interface {
	Element

	// FlexibleHeightFor returns what the element's minimum height would be
	// if resized to a specified width. This does not actually alter the
	// state of the element in any way, but it may perform significant work,
	// so it should be called sparingly.
	//
	// It is reccomended that parent containers check for this interface and
	// take this method's value into account in order to support things like
	// flow layouts and text wrapping, but it is not absolutely necessary.
	// The element's MinimumSize method will still return the absolute
	// minimum size that the element may be resized to.
	//
	// It is important to note that if a parent container checks for
	// flexible chilren, it itself will likely need to be flexible.
	FlexibleHeightFor (width int) (height int)

	// OnFlexibleHeightChange sets a function to be called when the
	// parameters affecting this element's flexible height are changed.
	OnFlexibleHeightChange (callback func ())
}

// Scrollable represents an element that can be scrolled. It acts as a viewport
// through which its contents can be observed.
type Scrollable interface {
	Element

	// ScrollContentBounds returns the full content size of the element.
	ScrollContentBounds () (bounds image.Rectangle)

	// ScrollViewportBounds returns the size and position of the element's
	// viewport relative to ScrollBounds.
	ScrollViewportBounds () (bounds image.Rectangle)

	// ScrollTo scrolls the viewport to the specified point relative to
	// ScrollBounds.
	ScrollTo (position image.Point)

	// ScrollAxes returns the supported axes for scrolling.
	ScrollAxes () (horizontal, vertical bool)

	// OnScrollBoundsChange sets a function to be called when the element's
	// ScrollContentBounds or ScrollViewportBounds are changed.
	OnScrollBoundsChange (callback func ())
}
