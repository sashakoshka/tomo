package core

import "image"
import "git.tebibyte.media/sashakoshka/tomo/layouts"
import "git.tebibyte.media/sashakoshka/tomo/elements"

// ContainerCore is a struct that can be embedded into an object to allow it to
// have one or more children. It also implements Flexible and Focusable and
// provides the standard behavior for selecting multiple children, and
// propagating user input events to them.
type ContainerCore struct {
	bounds    image.Rectangle
	layout    layouts.Layout
	children  []layouts.LayoutEntry
	drags     [10]elements.MouseTarget
	warping   bool
	focused   bool
	focusable bool
	flexible  bool
}

func NewContainerCore (
	layout         layouts.Layout,
	onFocusChange  func (),
	onLayoutChange func (),
) (
	core    *ContainerCore,
	control ContainerCoreControl,
) {
	core = &ContainerCore {
		layout: layout,
	}
	control = ContainerCoreControl {
		core: core,
	}
	return
}

// TODO fulfill interfaces here. accessors and mutators need to be in the
// container core control, because elements will have different ways of adopting
// and disowning child elements.

type ContainerCoreControl struct {
	core *ContainerCore
}

// Resize sets the size of the control, and
func (control ContainerCoreControl) Resize (bounds image.Rectangle) {
	// TODO do a layout
	// TODO call onLayoutChange
}

func (control ContainerCoreControl) Adopt (element elements.Element, expand bool) {
	
}
