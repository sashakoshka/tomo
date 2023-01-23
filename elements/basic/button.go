package basic

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// Button is a clickable button.
type Button struct {
	*core.Core
	core core.CoreControl
	
	pressed  bool
	enabled  bool
	selected bool

	text   string
	drawer artist.TextDrawer
	
	onClick func ()
	onSelectionRequest func () (granted bool)
	onSelectionMotionRequest func (tomo.SelectionDirection) (granted bool)
}

// NewButton creates a new button with the specified label text.
func NewButton (text string) (element *Button) {
	element = &Button { enabled: true }
	element.Core, element.core = core.NewCore(element)
	element.drawer.SetFace(theme.FontFaceRegular())
	element.SetText(text)
	return
}

func (element *Button) Resize (width, height int) {
	element.core.AllocateCanvas(width, height)
	element.draw()
}

func (element *Button) HandleMouseDown (x, y int, button tomo.Button) {
	if !element.enabled { return }
	if !element.selected { element.Select() }
	if button != tomo.ButtonLeft { return }
	element.pressed = true
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Button) HandleMouseUp (x, y int, button tomo.Button) {
	if button != tomo.ButtonLeft { return }
	element.pressed = false
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}

	within := image.Point { x, y }.
		In(element.Bounds())
		
	if !element.enabled { return }
	if within && element.onClick != nil {
		element.onClick()
	}
}

func (element *Button) HandleMouseMove (x, y int) { }
func (element *Button) HandleMouseScroll (x, y int, deltaX, deltaY float64) { }

func (element *Button) HandleKeyDown (key tomo.Key, modifiers tomo.Modifiers) {
	if !element.enabled { return }
	if key == tomo.KeyEnter {
		element.pressed = true
		if element.core.HasImage() {
			element.draw()
			element.core.DamageAll()
		}
	}
}

func (element *Button) HandleKeyUp(key tomo.Key, modifiers tomo.Modifiers) {
	if key == tomo.KeyEnter && element.pressed {
		element.pressed = false
		if element.core.HasImage() {
			element.draw()
			element.core.DamageAll()
		}
		if !element.enabled { return }
		if element.onClick != nil {
			element.onClick()
		}
	}
}

func (element *Button) Selected () (selected bool) {
	return element.selected
}

func (element *Button) Select () {
	if !element.enabled { return }
	if element.onSelectionRequest != nil {
		element.onSelectionRequest()
	}
}

func (element *Button) HandleSelection (
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

func (element *Button) HandleDeselection () {
	element.selected = false
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Button) OnSelectionRequest (callback func () (granted bool)) {
	element.onSelectionRequest = callback
}

func (element *Button) OnSelectionMotionRequest (
	callback func (direction tomo.SelectionDirection) (granted bool),
) {
	element.onSelectionMotionRequest = callback
}

// OnClick sets the function to be called when the button is clicked.
func (element *Button) OnClick (callback func ()) {
	element.onClick = callback
}

// SetEnabled sets whether this button can be clicked or not.
func (element *Button) SetEnabled (enabled bool) {
	if element.enabled == enabled { return }
	element.enabled = enabled
	if element.core.HasImage () {
		element.draw()
		element.core.DamageAll()
	}
}

// SetText sets the button's label text.
func (element *Button) SetText (text string) {
	if element.text == text { return }

	element.text = text
	element.drawer.SetText([]rune(text))
	textBounds := element.drawer.LayoutBounds()
	element.core.SetMinimumSize (
		theme.Padding() * 2 + textBounds.Dx(),
		theme.Padding() * 2 + textBounds.Dy())
	if element.core.HasImage () {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Button) draw () {
	bounds := element.core.Bounds()

	artist.FillRectangle (
		element.core,
		theme.ButtonPattern (
			element.enabled,
			element.Selected(),
			element.pressed),
		bounds)
		
	innerBounds := bounds
	innerBounds.Min.X += theme.Padding()
	innerBounds.Min.Y += theme.Padding()
	innerBounds.Max.X -= theme.Padding()
	innerBounds.Max.Y -= theme.Padding()

	textBounds := element.drawer.LayoutBounds()
	offset := image.Point {
		X: theme.Padding() + (innerBounds.Dx() - textBounds.Dx()) / 2,
		Y: theme.Padding() + (innerBounds.Dy() - textBounds.Dy()) / 2,
	}

	// account for the fact that the bounding rectangle will be shifted over
	// due to the bounds origin being at the baseline of the first line
	offset.Y -= textBounds.Min.Y
	offset.X -= textBounds.Min.X

	if element.pressed {
		offset = offset.Add(theme.SinkOffsetVector())
	}

	foreground := theme.ForegroundPattern(element.enabled)
	element.drawer.Draw(element.core, foreground, offset)
}
