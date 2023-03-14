package core

import "image"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/elements"

// Parent represents an object that can provide access to a list of child
// elements.
type Parent interface {
	Child         (index int) elements.Element
	CountChildren () int
}

// Propagator is a struct that can be embedded into elements that contain one or
// more children in order to propagate events to them without having to write
// all of the event handlers. It also implements standard behavior for focus
// propagation and keyboard navigation.
type Propagator struct {
	parent   Parent
	drags    [10]elements.MouseTarget
	focused  bool

	onFocusRequest func () (granted bool)
}

// NewPropagator creates a new event propagator that uses the specified parent
// to access a list of child elements that will have events propagated to them.
// If parent is nil, the function will return nil.
func NewPropagator (parent Parent) (propagator *Propagator) {
	if parent == nil { return nil }
	propagator = &Propagator {
		parent: parent,
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
	if propagator.onFocusRequest != nil {
		propagator.onFocusRequest()
	}
}

// HandleFocus causes this element to mark itself as focused. If the
// element does not have children or there are no more focusable children in
// the given direction, it should return false and do nothing. Otherwise, it
// marks itself as focused along with any applicable children and returns
// true.
func (propagator *Propagator) HandleFocus (direction input.KeynavDirection) (accepted bool) {
	direction = direction.Canon()

	firstFocused := propagator.firstFocused()
	if firstFocused < 0 {
		// no element is currently focused, so we need to focus either
		// the first or last focusable element depending on the
		// direction.
		switch direction {
		case input.KeynavDirectionForward:
			// if we recieve a forward direction, focus the first
			// focusable element.
			return propagator.focusFirstFocusableElement(direction)
		
		case input.KeynavDirectionBackward:
			// if we recieve a backward direction, focus the last
			// focusable element.
			return propagator.focusLastFocusableElement(direction)

		case input.KeynavDirectionNeutral:
			// if we recieve a neutral direction, just focus this
			// element and nothing else.
			propagator.focused = true
			return true
		}
	} else {
		// an element is currently focused, so we need to move the
		// focus in the specified direction
		firstFocusedChild :=
			propagator.parent.Child(firstFocused).
			(elements.Focusable)

		// before we move the focus, the currently focused child
		// may also be able to move its focus. if the child is able
		// to do that, we will let it and not move ours.
		if firstFocusedChild.HandleFocus(direction) {
			return true
		}

		// find the previous/next focusable element relative to the
		// currently focused element, if it exists.
		for index := firstFocused + int(direction);
			index < propagator.parent.CountChildren() && index >= 0;
			index += int(direction) {

			child, focusable :=
				propagator.parent.Child(index).
				(elements.Focusable)
			if focusable && child.HandleFocus(direction) {
				// we have found one, so we now actually move
				// the focus.
				firstFocusedChild.HandleUnfocus()
				propagator.focused = true
				return true
			}
		}
	}
	
	return false
}

// HandleDeselection causes this element to mark itself and all of its children
// as unfocused.
func (propagator *Propagator) HandleUnfocus () {
	propagator.forFocusable (func (child elements.Focusable) bool {
		child.HandleUnfocus()
		return true
	})
	propagator.focused = false
}

// OnFocusRequest sets a function to be called when this element wants its
// parent element to focus it. Parent elements should return true if the request
// was granted, and false if it was not. If the parent element returns true, the
// element acts as if a HandleFocus call was made with KeynavDirectionNeutral.
func (propagator *Propagator) OnFocusRequest (callback func () (granted bool)) {
	propagator.onFocusRequest = callback
}

// OnFocusMotionRequest sets a function to be called when this element wants its
// parent element to focus the element behind or in front of it, depending on
// the specified direction. Parent elements should return true if the request
// was granted, and false if it was not.
func (propagator *Propagator) OnFocusMotionRequest (
	callback func (direction input.KeynavDirection) (granted bool),
) { }

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
	propagator.forChildren (func (child elements.Element) bool {
		typedChild, themeable := child.(elements.Themeable)
		if themeable {
			typedChild.SetTheme(theme)
		}
		return true
	})
}

// SetConfig sets the theme of all children to the specified config.
func (propagator *Propagator) SetConfig (config config.Config) {
	propagator.forChildren (func (child elements.Element) bool {
		typedChild, configurable := child.(elements.Configurable)
		if configurable {
			typedChild.SetConfig(config)
		}
		return true
	})
}

// ----------- Focusing utilities ----------- //

func (propagator *Propagator) focusFirstFocusableElement (
	direction input.KeynavDirection,
) (
	ok bool,
) {
	propagator.forFocusable (func (child elements.Focusable) bool {
		if child.HandleFocus(direction) {
			propagator.focused = true
			ok = true
			return false
		}
		return true
	})
	return
}

func (propagator *Propagator) focusLastFocusableElement (
	direction input.KeynavDirection,
) (
	ok bool,
) {
	propagator.forChildrenReverse (func (child elements.Element) bool {
		typedChild, focusable := child.(elements.Focusable)
		if focusable && typedChild.HandleFocus(direction) {
			propagator.focused = true
			ok = true
			return false
		}
		return true
	})
	return
}

// ----------- Iterator utilities ----------- //

func (propagator *Propagator) forChildren (callback func (child elements.Element) bool) {
	for index := 0; index < propagator.parent.CountChildren(); index ++ {
		child := propagator.parent.Child(index)
		if child == nil     { continue }
		if !callback(child) { break    }
	}
}

func (propagator *Propagator) forChildrenReverse (callback func (child elements.Element) bool) {
	for index := propagator.parent.CountChildren() - 1; index > 0; index -- {
		child := propagator.parent.Child(index)
		if child == nil     { continue }
		if !callback(child) { break    }
	}
}

func (propagator *Propagator) childAt (position image.Point) (child elements.Element) {
	propagator.forChildren (func (current elements.Element) bool {
		if position.In(current.Bounds()) {
			child = current
		}
		return true
	})
	return
}

func (propagator *Propagator) forFocused (callback func (child elements.Focusable) bool) {
	propagator.forChildren (func (child elements.Element) bool {
		typedChild, focusable := child.(elements.Focusable)
		if focusable && typedChild.Focused() {
			if !callback(typedChild) { return false }
		}
		return true
	})
}

func (propagator *Propagator) forFocusable (callback func (child elements.Focusable) bool) {
	propagator.forChildren (func (child elements.Element) bool {
		typedChild, focusable := child.(elements.Focusable)
		if focusable {
			if !callback(typedChild) { return false }
		}
		return true
	})
}

func (propagator *Propagator) firstFocused () int {
	for index := 0; index < propagator.parent.CountChildren(); index ++ {
		child, focusable := propagator.parent.Child(index).(elements.Focusable)
		if focusable && child.Focused() {
			return index
		}
	}
	return -1
}
