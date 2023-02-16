package elements

import "image"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/config"

// Element represents a basic on-screen object.
type Element interface {
	// Bounds reports the element's bounding box. This must reflect the
	// bounding box of the last canvas given to the element by DrawTo.
	Bounds () (bounds image.Rectangle)

	// DrawTo sets this element's canvas. This should only be called by the
	// parent element. This is typically a region of the parent element's
	// canvas.
	DrawTo (canvas canvas.Canvas)

	// OnDamage sets a function to be called when an area of the element is
	// drawn on and should be pushed to the screen.
	OnDamage (callback func (region canvas.Canvas))

	// MinimumSize specifies the minimum amount of pixels this element's
	// width and height may be set to. If the element is given a resize
	// event with dimensions smaller than this, it will use its minimum
	// instead of the offending dimension(s).
	MinimumSize () (width, height int)

	// OnMinimumSizeChange sets a function to be called when the element's
	// minimum size is changed.
	OnMinimumSizeChange (callback func ())
}

// Focusable represents an element that has keyboard navigation support. This
// includes inputs, buttons, sliders, etc. as well as any elements that have
// children (so keyboard navigation events can be propagated downward).
type Focusable interface {
	Element

	// Focused returns whether or not this element or any of its children
	// are currently focused.
	Focused () (selected bool)

	// Focus focuses this element, if its parent element grants the
	// request.
	Focus ()

	// HandleFocus causes this element to mark itself as focused. If the
	// element does not have children, it is disabled, or there are no more
	// selectable children in the given direction, it should return false
	// and do nothing. Otherwise, it should select itself and any children
	// (if applicable) and return true.
	HandleFocus (direction input.KeynavDirection) (accepted bool)

	// HandleDeselection causes this element to mark itself and all of its
	// children as unfocused.
	HandleUnfocus ()

	// OnFocusRequest sets a function to be called when this element wants
	// its parent element to focus it. Parent elements should return true if
	// the request was granted, and false if it was not. If the parent
	// element returns true, the element must act as if a HandleFocus call
	// was made with KeynavDirectionNeutral.
	OnFocusRequest (func () (granted bool))

	// OnFocusMotionRequest sets a function to be called when this
	// element wants its parent element to focus the element behind or in
	// front of it, depending on the specified direction. Parent elements
	// should return true if the request was granted, and false if it was
	// not.
	OnFocusMotionRequest (func (direction input.KeynavDirection) (granted bool))
}

// KeyboardTarget represents an element that can receive keyboard input.
type KeyboardTarget interface {
	Element

	// HandleKeyDown is called when a key is pressed down or repeated while
	// this element has keyboard focus. It is important to note that not
	// every key down event is guaranteed to be paired with exactly one key
	// up event. This is the reason a list of modifier keys held down at the
	// time of the key press is given.
	HandleKeyDown (key input.Key, modifiers input.Modifiers)

	// HandleKeyUp is called when a key is released while this element has
	// keyboard focus.
	HandleKeyUp (key input.Key, modifiers input.Modifiers)
}

// MouseTarget represents an element that can receive mouse events.
type MouseTarget interface {
	Element

	// Each of these handler methods is passed the position of the mouse
	// cursor at the time of the event as x, y.

	// HandleMouseDown is called when a mouse button is pressed down on this
	// element.
	HandleMouseDown (x, y int, button input.Button)

	// HandleMouseUp is called when a mouse button is released that was
	// originally pressed down on this element.
	HandleMouseUp (x, y int, button input.Button)

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
	// ScrollContentBounds, ScrollViewportBounds, or ScrollAxes are changed.
	OnScrollBoundsChange (callback func ())
}

// Collapsible represents an element who's minimum width and height can be
// manually resized. Scrollable elements should implement this if possible.
type Collapsible interface {
	Element

	// Collapse collapses the element's minimum width and height. A value of
	// zero for either means that the element's normal value is used.
	Collapse (width, height int)
}

// Themeable represents an element that can modify its appearance to fit within
// a theme.
type Themeable interface {
	Element
	
	// SetTheme sets the element's theme to something fulfilling the
	// theme.Theme interface.
	SetTheme (theme.Theme)
}

// Configurable represents an element that can modify its behavior to fit within
// a set of configuration parameters.
type Configurable interface {
	Element
	
	// SetConfig sets the element's configuration to something fulfilling
	// the config.Config interface.
	SetConfig (config.Config)
}
