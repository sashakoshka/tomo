package core

import "image"
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
	focused  bool
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

// ----------- Interface fulfillment methods ----------- //

// Focused returns whether or not this element or any of its children
// are currently focused.
func (propagator *Propagator) Focused () (focused bool) {
	return propagator.focused
}

// Focus focuses this element, if its parent element grants the
// request.
func (propagator *Propagator) Focus () {
	// TODO
}

// HandleFocus causes this element to mark itself as focused. If the
// element does not have children or there are no more focusable children in
// the given direction, it should return false and do nothing. Otherwise, it
// marks itself as focused along with any applicable children and returns
// true.
func (propagator *Propagator) HandleFocus (direction input.KeynavDirection) (accepted bool) {
	// TODO
}

// HandleDeselection causes this element to mark itself and all of its children
// as unfocused.
func (propagator *Propagator) HandleUnfocus () {
	// TODO
}

// OnFocusRequest sets a function to be called when this element wants its
// parent element to focus it. Parent elements should return true if the request
// was granted, and false if it was not. If the parent element returns true, the
// element acts as if a HandleFocus call was made with KeynavDirectionNeutral.
func (propagator *Propagator) OnFocusRequest (func () (granted bool)) {
	// TODO
}

// OnFocusMotionRequest sets a function to be called when this element wants its
// parent element to focus the element behind or in front of it, depending on
// the specified direction. Parent elements should return true if the request
// was granted, and false if it was not.
func (propagator *Propagator) OnFocusMotionRequest (func (direction input.KeynavDirection) (granted bool)) {
	// TODO
}

// HandleKeyDown propogates the keyboard event to the currently selected child.
func (propagator *Propagator) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
	propagator.forFocused (func (child elements.Focusable) bool {
		typedChild, handlesKeyboard := child.(elements.KeyboardTarget)
		if handlesKeyboard {
			typedChild.HandleKeyDown(key, modifiers)
		}
		return true
	})
}

// HandleKeyUp propogates the keyboard event to the currently selected child.
func (propagator *Propagator) HandleKeyUp (key input.Key, modifiers input.Modifiers) {
	propagator.forFocused (func (child elements.Focusable) bool {
		typedChild, handlesKeyboard := child.(elements.KeyboardTarget)
		if handlesKeyboard {
			typedChild.HandleKeyUp(key, modifiers)
		}
		return true
	})
}

// HandleMouseDown propagates the mouse event to the element under the mouse
// pointer.
func (propagator *Propagator) HandleMouseDown (x, y int, button input.Button) {
	child, handlesMouse :=
		propagator.childAt(image.Pt(x, y)).
		(elements.MouseTarget)
	if handlesMouse {
		propagator.drags[button] = child
		child.HandleMouseDown(x, y, button)
	}
}

// HandleMouseUp propagates the mouse event to the element that the released
// mouse button was originally pressed on.
func (propagator *Propagator) HandleMouseUp (x, y int, button input.Button) {
	child := propagator.drags[button]
	if child != nil {
		propagator.drags[button] = nil
		child.HandleMouseUp(x, y, button)
	}
}

// HandleMouseMove propagates the mouse event to the element that was last
// pressed down by the mouse if the mouse is currently being held down, else it
// propagates the event to whichever element is underneath the mouse pointer.
func (propagator *Propagator) HandleMouseMove (x, y int) {
	handled := false
	for _, child := range propagator.drags {
		if child != nil {
			child.HandleMouseMove(x, y)
			handled = true
		}
	}

	if handled {
		child, handlesMouse :=
			propagator.childAt(image.Pt(x, y)).
			(elements.MouseTarget)
		if handlesMouse {
			child.HandleMouseMove(x, y)
		}
	}
}

// HandleScroll propagates the mouse event to the element under the mouse
// pointer.
func (propagator *Propagator) HandleMouseScroll (x, y int, deltaX, deltaY float64) {
	child, handlesMouse :=
		propagator.childAt(image.Pt(x, y)).
		(elements.MouseTarget)
	if handlesMouse {
		child.HandleMouseScroll(x, y, deltaX, deltaY)
	}
}

// SetTheme sets the theme of all children to the specified theme.
func (propagator *Propagator) SetTheme (theme theme.Theme) {
	propagator.iterator.OverChildren (func (child elements.Element) bool {
		typedChild, themeable := child.(elements.Themeable)
		if themeable {
			typedChild.SetTheme(theme)
		}
		return true
	})
}

// SetConfig sets the theme of all children to the specified config.
func (propagator *Propagator) SetConfig (config config.Config) {
	propagator.iterator.OverChildren (func (child elements.Element) bool {
		typedChild, configurable := child.(elements.Configurable)
		if configurable {
			typedChild.SetConfig(config)
		}
		return true
	})
}

// ----------- Iterator utilities ----------- //

func (propagator *Propagator) childAt (position image.Point) (child elements.Element) {
	propagator.iterator.OverChildren (func (current elements.Element) bool {
		if position.In(current.Bounds()) {
			child = current
		}
		return true
	})
	return
}

func (propagator *Propagator) forFocused (callback func (child elements.Focusable) bool) {
	propagator.iterator.OverChildren (func (child elements.Element) bool {
		typedChild, focusable := child.(elements.Focusable)
		if focusable && typedChild.Focused() {
			if !callback(typedChild) { return false }
		}
		return true
	})
}

func (propagator *Propagator) forFocusable (callback func (child elements.Focusable) bool) {
	propagator.iterator.OverChildren (func (child elements.Element) bool {
		typedChild, focusable := child.(elements.Focusable)
		if focusable {
			if !callback(typedChild) { return false }
		}
		return true
	})
}

func (propagator *Propagator) forFlexible (callback func (child elements.Flexible) bool) {
	propagator.iterator.OverChildren (func (child elements.Element) bool {
		typedChild, flexible := child.(elements.Flexible)
		if flexible {
			if !callback(typedChild) { return false }
		}
		return true
	})
}

// func (propagator *Propagator) forFocusableBackward (callback func (child elements.Focusable) bool) {
	// for index := len(element.children) - 1; index >= 0; index -- {
		// child, focusable := element.children[index].Element.(elements.Focusable)
		// if focusable {
			// if !callback(child) { break }
		// }
	// }
// }

func (propagator *Propagator) firstFocused () (index int) {
	index = -1
	currentIndex := 0
	propagator.forFocusable (func (child elements.Focusable) bool {
		if child.Focused() {
			index = currentIndex
			return false
		}
		currentIndex ++
		return true
	})
	return
}
