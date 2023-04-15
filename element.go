package tomo

import "image"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/canvas"

// Element represents a basic on-screen object.
type Element interface {
	// Draw causes the element to draw to the specified canvas. The bounds
	// of this canvas specify the area that is actually drawn to, while the
	// Entity bounds specify the actual area of the element.
	Draw (canvas.Canvas)

	// Entity returns this element's entity.
	Entity () Entity
}

// Container is an element capable of containing child elements.
type Container interface {
	Element

	// Layout causes this element to arrange its children.
	Layout ()

	// DrawBackground causes the element to draw its background pattern to
	// the specified canvas. The bounds of this canvas specify the area that
	// is actually drawn to, while the Entity bounds specify the actual area
	// of the element.
	DrawBackground (canvas.Canvas)

	// HandleChildMinimumSizeChange is called when a child's minimum size is
	// changed.
	HandleChildMinimumSizeChange (child Element)
}

// Focusable represents an element that has keyboard navigation support.
type Focusable interface {
	Element

	// HandleFocusChange is called when the element is focused or unfocused.
	HandleFocusChange ()
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

	// HandleMouseDown is called when a mouse button is pressed down on this
	// element.
	HandleMouseDown (x, y int, button input.Button)

	// HandleMouseUp is called when a mouse button is released that was
	// originally pressed down on this element.
	HandleMouseUp (x, y int, button input.Button)
}

// MotionTarget represents an element that can receive mouse motion events.
type MotionTarget interface {
	Element

	// HandleMotion is called when the mouse is moved over this element,
	// or the mouse is moving while being held down and originally pressed
	// down on this element.
	HandleMotion (x, y int)
}

// ScrollTarget represents an element that can receive mouse scroll events.
type ScrollTarget interface {
	Element

	// HandleScroll is called when the mouse is scrolled. The X and Y
	// direction of the scroll event are passed as deltaX and deltaY.
	HandleScroll (x, y int, deltaX, deltaY float64)
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
	FlexibleHeightFor (width int) int
}

// FlexibleContainer represents an element that is capable of containing
// flexible children.
type FlexibleContainer interface {
	Container

	// HandleChildFlexibleHeightChange is called when the parameters
	// affecting a child's flexible height are changed.
	HandleChildFlexibleHeightChange (child Flexible)
}

// Scrollable represents an element that can be scrolled. It acts as a viewport
// through which its contents can be observed.
type Scrollable interface {
	Element

	// ScrollContentBounds returns the full content size of the element.
	ScrollContentBounds () image.Rectangle

	// ScrollViewportBounds returns the size and position of the element's
	// viewport relative to ScrollBounds.
	ScrollViewportBounds () image.Rectangle

	// ScrollTo scrolls the viewport to the specified point relative to
	// ScrollBounds.
	ScrollTo (position image.Point)

	// ScrollAxes returns the supported axes for scrolling.
	ScrollAxes () (horizontal, vertical bool)
}

// ScrollableContainer represents an element that is capable of containing
// scrollable children.
type ScrollableContainer interface {
	Container

	// HandleChildScrollBoundsChange is called when the content bounds,
	// viewport bounds, or scroll axes of a child are changed.
	HandleChildScrollBoundsChange (child Scrollable)
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
	SetTheme (Theme)
}

// Configurable represents an element that can modify its behavior to fit within
// a set of configuration parameters.
type Configurable interface {
	Element
	
	// SetConfig sets the element's configuration to something fulfilling
	// the config.Config interface.
	SetConfig (Config)
}
