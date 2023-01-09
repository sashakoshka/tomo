package tomo

// Event represents any event. Use a type switch to figure out what sort of
// event it is.
type Event interface { }

// EventResize is sent to an element when its parent decides to resize it.
// Elements should not do anything if the width and height do not change.
type EventResize struct {
	// The width and height the element should not be less than the
	// element's reported minimum width and height. If by some chance they
	// are anyways, the element should use its minimum width and height
	// instead.
	Width, Height int
}

// EventKeyDown is sent to the currently selected element when a key on the
// keyboard is pressed. Containers must propagate this event downwards.
type EventKeyDown struct {
	Key
	Modifiers
	Repeated bool
}

// EventKeyDown is sent to the currently selected element when a key on the
// keyboard is released. Containers must propagate this event downwards.
type EventKeyUp struct {
	Key
	Modifiers
}

// EventMouseDown is sent to the element the mouse is positioned over when it is
// clicked. Containers must propagate this event downwards, with X and Y values
// relative to the top left corner of the child element.
type EventMouseDown struct {
	// The button that was released
	Button
	
	// The X and Y position of the mouse cursor at the time of the event,
	// relative to the top left corner of the element
	X, Y int
}

// EventMouseUp is sent to the element that was positioned under the mouse the
// last time this particular mouse button was pressed down when it is released.
// Containers must propagate this event downwards, with X and Y values relative
// to the top left corner of the child element.
type EventMouseUp struct {
	// The button that was released
	Button
	
	// The X and Y position of the mouse cursor at the time of the event,
	// relative to the top left corner of the element
	X, Y int
}

// EventMouseMove is sent to the element positioned under the mouse cursor when
// the mouse moves, or if a mouse button is currently being pressed, the element
// that the mouse was positioned under when it was pressed down. Containers must
// propogate this event downwards, with X and Y values relative to the top left
// corner of the child element.
type EventMouseMove struct {
	// The X and Y position of the mouse cursor at the time of the event,
	// relative to the top left corner of the element
	X, Y int
}

// EventScroll is sent to the element positioned under the mouse cursor when the
// scroll wheel (or equivalent) is spun. Containers must propogate this event
// downwards.
type EventScroll struct {
	// The X and Y position of the mouse cursor at the time of the event,
	// relative to the top left corner of the element
	X, Y int

	// The X and Y amount the scroll wheel moved
	ScrollX, ScrollY int
}

// EventSelect is sent to selectable elements when they become selected, whether
// by a mouse click or by keyboard navigation. Containers must propagate this
// event downwards.
type EventSelect struct { }

// EventDeselect is sent to selectable elements when they stop being selected,
// whether by a mouse click or by keyboard navigation. Containers must propagate
// this event downwards.
type EventDeselect struct { }
