package basic

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// Checkbox is a toggle-able checkbox with a label.
type Checkbox struct {
	*core.Core
	core core.CoreControl

	pressed  bool
	checked  bool
	enabled  bool
	selected bool

	text   string
	drawer artist.TextDrawer
	
	onClick func ()
	onSelectionRequest func () (granted bool)
	onSelectionMotionRequest func (tomo.SelectionDirection) (granted bool)
}

// NewCheckbox creates a new cbeckbox with the specified label text.
func NewCheckbox (text string, checked bool) (element *Checkbox) {
	element = &Checkbox { enabled: true, checked: checked }
	element.Core, element.core = core.NewCore(element)
	element.drawer.SetFace(theme.FontFaceRegular())
	element.SetText(text)
	return
}

// Resize changes this element's size.
func (element *Checkbox) Resize (width, height int) {
	element.core.AllocateCanvas(width, height)
	element.draw()
}

func (element *Checkbox) HandleMouseDown (x, y int, button tomo.Button) {
	element.Select()
	element.pressed = true
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Checkbox) HandleMouseUp (x, y int, button tomo.Button) {
	if button != tomo.ButtonLeft { return }

	element.pressed = false
	within := image.Point { x, y }.
		In(element.Bounds())
	if within {
		element.checked = !element.checked
	}
	
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
	if within && element.onClick != nil {
		element.onClick()
	}
}

func (element *Checkbox) HandleMouseMove (x, y int) { }
func (element *Checkbox) HandleMouseScroll (x, y int, deltaX, deltaY float64) { }

func (element *Checkbox) HandleKeyDown (key tomo.Key, modifiers tomo.Modifiers) {
	if key == tomo.KeyEnter {
		element.pressed = true
		if element.core.HasImage() {
			element.draw()
			element.core.DamageAll()
		}
	}
}

func (element *Checkbox) HandleKeyUp (key tomo.Key, modifiers tomo.Modifiers) {
	if key == tomo.KeyEnter && element.pressed {
		element.pressed = false
		element.checked = !element.checked
		if element.core.HasImage() {
			element.draw()
			element.core.DamageAll()
		}
		if element.onClick != nil {
			element.onClick()
		}
	}
}

// Selected returns whether or not this element is selected.
func (element *Checkbox) Selected () (selected bool) {
	return element.selected
}

// Select requests that this element be selected.
func (element *Checkbox) Select () {
	if !element.enabled { return }
	if element.onSelectionRequest != nil {
		element.onSelectionRequest()
	}
}

func (element *Checkbox) HandleSelection (
	direction tomo.SelectionDirection,
) (
	accepted bool,
) {
	direction = direction.Canon()
	if !element.enabled { return false }
	if element.selected && direction != tomo.SelectionDirectionNeutral {
		return false
	}
	
	element.selected = true
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
	return true
}

func (element *Checkbox) HandleDeselection () {
	element.selected = false
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Checkbox) OnSelectionRequest (callback func () (granted bool)) {
	element.onSelectionRequest = callback
}

func (element *Checkbox) OnSelectionMotionRequest (
	callback func (direction tomo.SelectionDirection) (granted bool),
) {
	element.onSelectionMotionRequest = callback
}

// OnClick sets the function to be called when the checkbox is toggled.
func (element *Checkbox) OnClick (callback func ()) {
	element.onClick = callback
}

// Value reports whether or not the checkbox is currently checked.
func (element *Checkbox) Value () (checked bool) {
	return element.checked
}

// SetEnabled sets whether this checkbox can be toggled or not.
func (element *Checkbox) SetEnabled (enabled bool) {
	if element.enabled == enabled { return }
	element.enabled = enabled
	if element.core.HasImage () {
		element.draw()
		element.core.DamageAll()
	}
}

// SetText sets the checkbox's label text.
func (element *Checkbox) SetText (text string) {
	if element.text == text { return }

	element.text = text
	element.drawer.SetText([]rune(text))
	textBounds := element.drawer.LayoutBounds()
	element.core.SetMinimumSize (
		textBounds.Dy() + theme.Padding() + textBounds.Dx(),
		textBounds.Dy())
	if element.core.HasImage () {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Checkbox) draw () {
	bounds := element.core.Bounds()
	boxBounds := image.Rect(0, 0, bounds.Dy(), bounds.Dy())

	artist.FillRectangle ( element.core, theme.BackgroundPattern(), bounds)
	artist.FillRectangle (
		element.core,
		theme.ButtonPattern (
			element.enabled,
			element.Selected(),
			element.pressed),
		boxBounds)
		
	innerBounds := bounds
	innerBounds.Min.X += theme.Padding()
	innerBounds.Min.Y += theme.Padding()
	innerBounds.Max.X -= theme.Padding()
	innerBounds.Max.Y -= theme.Padding()

	textBounds := element.drawer.LayoutBounds()
	offset := image.Point {
		X: bounds.Dy() + theme.Padding(),
	}

	offset.Y -= textBounds.Min.Y
	offset.X -= textBounds.Min.X

	foreground := theme.ForegroundPattern(element.enabled)
	element.drawer.Draw(element.core, foreground, offset)
	
	if element.checked {
		checkBounds := boxBounds.Inset(4)
		if element.pressed {
			checkBounds = checkBounds.Add(theme.SinkOffsetVector())
		}
		artist.FillRectangle (
			element.core,
			theme.ForegroundPattern(element.enabled),
			checkBounds)
	}
}
