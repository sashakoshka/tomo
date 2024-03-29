package tomo

import "image"
import "art"

// Entity is a handle given to elements by the backend. Extended entity
// interfaces are defined in the ability module.
type Entity interface {
	// Invalidate marks the element's current visual as invalid. At the end
	// of every event, the backend will ask all invalid entities to redraw
	// themselves.
	Invalidate ()

	// InvalidateLayout marks the element's layout as invalid. At the end of
	// every event, the backend will ask all invalid elements to recalculate
	// their layouts.
	InvalidateLayout ()

	// Bounds returns the bounds of the element to be used for drawing and
	// layout.
	Bounds () image.Rectangle

	// Window returns the window that the element is in.
	Window () Window

	// SetMinimumSize reports to the system what the element's minimum size
	// can be. The minimum size of child elements should be taken into
	// account when calculating this.
	SetMinimumSize (width, height int)

	// DrawBackground asks the parent element to draw its background pattern
	// to a canvas. This should be used for transparent elements like text
	// labels. If there is no parent element (that is, the element is
	// directly inside of the window), the backend will draw a default
	// background pattern.
	DrawBackground (art.Canvas)

	// --- Behaviors relating to parenting ---

	// Adopt adds an element as a child.
	Adopt (child Element)

	// Insert inserts an element in the child list at the specified
	// location.
	Insert (index int, child Element)

	// Disown removes the child at the specified index.
	Disown (index int)

	// IndexOf returns the index of the specified child.
	IndexOf (child Element) int

	// Child returns the child at the specified index.
	Child (index int) Element

	// CountChildren returns the amount of children the element has.
	CountChildren () int

	// PlaceChild sets the size and position of the child at the specified
	// index to a bounding rectangle.
	PlaceChild (index int, bounds image.Rectangle)

	// SelectChild marks a child as selected or unselected, if it is
	// selectable.
	SelectChild (index int, selected bool)

	// ChildMinimumSize returns the minimum size of the child at the
	// specified index.
	ChildMinimumSize (index int) (width, height int)

	// --- Behaviors relating to input focus ---

	// Focused returns whether the element currently has input focus.
	Focused () bool

	// Focus sets this element as focused. If this succeeds, the element will
	// recieve a HandleFocus call.
	Focus ()

	// FocusNext causes the focus to move to the next element. If this
	// succeeds, the element will recieve a HandleUnfocus call.
	FocusNext ()

	// FocusPrevious causes the focus to move to the next element. If this
	// succeeds, the element will recieve a HandleUnfocus call.	
	FocusPrevious ()

	// Selected returns whether this element is currently selected.
	Selected () bool

	// --- Behaviors relating to scrolling --- //

	// NotifyFlexibleHeightChange notifies the system that the parameters
	// affecting the element's flexible height have changed. This method is
	// expected to be called by flexible elements when their content changes.
	NotifyFlexibleHeightChange ()

	// NotifyScrollBoundsChange notifies the system that the element's
	// scroll content bounds or viewport bounds have changed. This is
	// expected to be called by scrollable elements when they change their
	// supported scroll axes, their scroll position (either autonomously or
	// as a result of a call to ScrollTo()), or their content size.
	NotifyScrollBoundsChange ()

	// --- Behaviors relating to themeing ---

	// Theme returns the currently active theme. When this value changes,
	// the HandleThemeChange method of the element is called.
	Theme () Theme

	// Config returns the currently active config. When this value changes,
	// the HandleThemeChange method of the element is called.
	Config () Config
}
