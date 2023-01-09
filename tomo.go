package tomo

import "image"
import "errors"
import "image/draw"
import "image/color"

// Image represents a simple image buffer that fulfills the image.Image
// interface while also having methods that do away with the use of the
// color.Color interface to facilitate more efficient drawing. This interface
// can be easily satisfied using an image.RGBA struct.
type Image interface {
	image.Image
	RGBAAt (x, y int) (c color.RGBA)
}

// Canvas is like Image but also requires Set and SetRGBA methods. This
// interface can be easily satisfied using an image.RGBA struct.
type Canvas interface {
	draw.Image
	RGBAAt (x, y int) (c color.RGBA)
	SetRGBA (x, y int, c color.RGBA)
}

// ParentHooks is a struct that contains callbacks that let child elements send
// information to their parent element without the child element knowing
// anything about the parent element or containing any reference to it. When a
// parent element adopts a child element, it must set these callbacks.
type ParentHooks struct {
	// Draw is called when a part of the child element's surface is updated.
	// The updated region will be passed to the callback as a sub-image.
	Draw func (region Image)

	// MinimumSizeChange is called when the child element's minimum width
	// and/or height changes. When this function is called, the element will
	// have already been resized and there is no need to send it a resize
	// event.
	MinimumSizeChange func (width, height int)

	// SelectionRequest is called when the child element element wants
	// itself to be selected. If the parent element chooses to grant the
	// request, it must send the child element a selection event.
	SelectionRequest func ()
}

// RunDraw runs the Draw hook if it is not nil. If it is nil, it does nothing.
func (hooks ParentHooks) RunDraw (region Image) {
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
func (hooks ParentHooks) RunSelectionRequest () {
	if hooks.SelectionRequest != nil {
		hooks.SelectionRequest()
	}
}

// Element represents a basic on-screen object.
type Element interface {
	// Element must implement the Image interface. Elements should start out
	// with a completely blank image buffer, and only set its size and draw
	// on it for the first time when sent an EventResize event.
	Image

	// Handle handles an event, propagating it to children if necessary.
	Handle (event Event)

	// Selectable returns whether this element can be selected. If this
	// element contains other selectable elements, it must return true.
	Selectable () (bool)

	// If this element contains other elements, and one is selected, this
	// method will advance the selection in the specified direction. If no
	// children are selected, or there are no more children to be selected
	// in the specified direction, the element will unselect all of its
	// children and return false. If the selection could be advanced, it
	// will return true. If the element contains no child elements, it will
	// always return false.
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

// Window represents a top-level container generated by the currently running
// backend. It can contain a single element. It is hidden by default, and must
// be explicitly shown with the Show() method. If it contains no element, it
// displays a black (or transprent) background.
type Window interface {
	// Adopt sets the root element of the window. There can only be one of
	// these at one time.
	Adopt (child Element)

	// Child returns the root element of the window.
	Child () (child Element)

	// SetTitle sets the title that appears on the window's title bar. This
	// method might have no effect with some backends.
	SetTitle (title string)

	// SetIcon taks in a list different sizes of the same icon and selects
	// the best one to display on the window title bar, dock, or whatever is
	// applicable for the given backend. This method might have no effect
	// with some backends.
	SetIcon (sizes []image.Image)

	// Show shows the window. The window starts off hidden, so this must be
	// called after initial setup to make sure it is visible.
	Show ()

	// Hide hides the window.
	Hide ()

	// Close closes the window.
	Close ()

	// OnClose specifies a function to be called when the window is closed.
	OnClose (func ())
}

var backend Backend

// Run initializes a backend, calls the callback function, and begins the event
// loop in that order. This function does not return until Stop() is called, or
// the backend experiences a fatal error.
func Run (callback func ()) (err error) {
	backend, err = instantiateBackend()
	if callback != nil { callback() }
	err = backend.Run()
	backend = nil
	return
}

// Stop gracefully stops the event loop and shuts the backend down. Call this
// before closing your application.
func Stop () {
	if backend != nil { backend.Stop() }
}

// Do executes the specified callback within the main thread as soon as
// possible. This function can be safely called from other threads.
func Do (callback func ()) {
	
}

// NewWindow creates a new window using the current backend, and returns it as a
// Window. If the window could not be created, an error is returned explaining
// why. If this function is called without a running backend, an error is
// returned as well.
func NewWindow (width, height int) (window Window, err error) {
	if backend == nil {
		err = errors.New("no backend is running.")
		return
	}
	window, err = backend.NewWindow(width, height)
	return
}
