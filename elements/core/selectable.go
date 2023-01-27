package core

import "git.tebibyte.media/sashakoshka/tomo"

// SelectableCore is a struct that can be embedded into objects to make them
// selectable, giving them the default selectability behavior.
type SelectableCore struct {
	selected bool
	enabled  bool
	drawSelectionChange func ()
	onSelectionRequest func () (granted bool)
	onSelectionMotionRequest func(tomo.SelectionDirection) (granted bool)
}

// NewSelectableCore creates a new selectability core and its corresponding
// control. If your element needs to visually update itself when it's selection
// state changes (which it should), a callback to draw and push the update can
// be specified. 
func NewSelectableCore (
	drawSelectionChange func (),
) (	
	core *SelectableCore,
	control SelectableCoreControl,
) {
	core = &SelectableCore {
		drawSelectionChange: drawSelectionChange,
		enabled: true,
	}
	control = SelectableCoreControl { core: core }
	return
}

// Selected returns whether or not this element is currently selected.
func (core *SelectableCore) Selected () (selected bool) {
	return core.selected
}

// Select selects this element, if its parent element grants the request.
func (core *SelectableCore) Select () {
	if !core.enabled { return }
	if core.onSelectionRequest != nil {
		core.onSelectionRequest()
	}
}

// HandleSelection causes this element to mark itself as selected, if it can
// currently be. Otherwise, it will return false and do nothing.
func (core *SelectableCore) HandleSelection (
	direction tomo.SelectionDirection,
) (
	accepted bool,
) {
	direction = direction.Canon()
	if !core.enabled { return false }
	if core.selected && direction != tomo.SelectionDirectionNeutral {
		return false
	}
	
	core.selected = true
	if core.drawSelectionChange != nil { core.drawSelectionChange() }
	return true
}

// HandleDeselection causes this element to mark itself as deselected.
func (core *SelectableCore) HandleDeselection () {
	core.selected = false
	if core.drawSelectionChange != nil { core.drawSelectionChange() }
}

// OnSelectionRequest sets a function to be called when this element
// wants its parent element to select it. Parent elements should return
// true if the request was granted, and false if it was not.
func (core *SelectableCore) OnSelectionRequest (callback func () (granted bool)) {
	core.onSelectionRequest = callback
}

// OnSelectionMotionRequest sets a function to be called when this
// element wants its parent element to select the element behind or in
// front of it, depending on the specified direction. Parent elements
// should return true if the request was granted, and false if it was
// not.
func (core *SelectableCore) OnSelectionMotionRequest (
	callback func (direction tomo.SelectionDirection) (granted bool),
) {
	core.onSelectionMotionRequest = callback
}

// Enabled returns whether or not the element is enabled.
func (core *SelectableCore) Enabled () (enabled bool) {
	return core.enabled
}

// SelectableCoreControl is a struct that can be used to exert control over a
// selectability core. It must not be directly embedded into an element, but
// instead kept as a private member. When a SelectableCore struct is created, a
// corresponding SelectableCoreControl struct is linked to it and returned
// alongside it.
type SelectableCoreControl struct {
	core *SelectableCore
}

// SetEnabled sets whether the selectability core is enabled. If the state
// changes, this will call drawSelectionChange.
func (control SelectableCoreControl) SetEnabled (enabled bool) {
	if control.core.enabled == enabled { return }
	control.core.enabled = enabled
	if !enabled { control.core.selected = false }
	if control.core.drawSelectionChange != nil {
		control.core.drawSelectionChange()
	}
}
