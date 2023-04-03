package tomo

import "image"

// Parent represents a type capable of containing child elements.
type Parent interface {
	// NotifyMinimumSizeChange notifies the container that a child element's
	// minimum size has changed. This method is expected to be called by
	// child elements when their minimum size changes.
	NotifyMinimumSizeChange (child Element)

	// Window returns the window containing the parent.
	Window () Window
}

// FocusableParent represents a parent with keyboard navigation support.
type FocusableParent interface {
	Parent
	
	// RequestFocus notifies the parent that a child element is requesting
	// keyboard focus. If the parent grants the request, the method will
	// return true and the child element should behave as if a HandleFocus
	// call was made.
	RequestFocus (child Focusable) (granted bool)

	// RequestFocusMotion notifies the parent that a child element wants the
	// focus to be moved to the next focusable element.
	RequestFocusNext (child Focusable)
	
	// RequestFocusMotion notifies the parent that a child element wants the
	// focus to be moved to the previous focusable element.
	RequestFocusPrevious (child Focusable)
}

// FlexibleParent represents a parent that accounts for elements with
// flexible height.
type FlexibleParent interface {
	Parent

	// NotifyFlexibleHeightChange notifies the parent that the parameters
	// affecting a child's flexible height have changed. This method is
	// expected to be called by flexible child element when their content
	// changes.
	NotifyFlexibleHeightChange (child Flexible)
}

// ScrollableParent represents a parent that can change the scroll
// position of its child element(s).
type ScrollableParent interface {
	Parent
	
	// NotifyScrollBoundsChange notifies the parent that a child's scroll
	// content bounds or viewport bounds have changed. This is expected to
	// be called by child elements when they change their supported scroll
	// axes, their scroll position (either autonomously or as a result of a
	// call to ScrollTo()), or their content size.
	NotifyScrollBoundsChange (child Scrollable)
}

// BackgroundParent represents a parent that is able to re-draw a portion of its
// background upon request. This is intended to be used by transparent elements
// that want to adopt their parent's background pattern. If a parent implements
// this interface, it should call a child's DrawTo method when its area of the
// background is affected.
type BackgroundParent interface {
	Parent

	// DrawBackground draws a portion of the parent's background pattern
	// within the specified bounds. The parent will not push these changes.
	DrawBackground (bounds image.Rectangle)
}
