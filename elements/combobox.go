package elements

import "image"
import "tomo"
import "tomo/input"
import "tomo/artist"
import "tomo/ability"
import "tomo/textdraw"

var comboBoxCase = tomo.C("tomo", "comboBox")

// Option specifies a ComboBox option. A blank option will display as "(None)".
type Option string

func (option Option) Title () string {
	if option == "" {
		return "(None)"
	} else {
		return string(option)
	}
}

// ComboBox is an input that can be one of several predetermined values.
type ComboBox struct {
	entity tomo.Entity
	drawer textdraw.Drawer

	options  []Option
	selected Option

	enabled bool
	pressed bool
	
	onChange func ()
}

// NewComboBox creates a new ComboBox with the specifed options.
func NewComboBox (options ...Option) (element *ComboBox) {
	if len(options) == 0 { options = []Option { "" } }
	element = &ComboBox { enabled: true, options: options }
	element.entity = tomo.GetBackend().NewEntity(element)
	element.drawer.SetFace (element.entity.Theme().FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal,
		comboBoxCase))
	element.Select(options[0])
	return
}

// Entity returns this element's entity.
func (element *ComboBox) Entity () tomo.Entity {
	return element.entity
}

// Draw causes the element to draw to the specified destination canvas.
func (element *ComboBox) Draw (destination artist.Canvas) {
	state   := element.state()
	bounds  := element.entity.Bounds()
	pattern := element.entity.Theme().Pattern(tomo.PatternButton, state, comboBoxCase)

	pattern.Draw(destination, bounds)
	
	foreground := element.entity.Theme().Color(tomo.ColorForeground, state, comboBoxCase)
	sink       := element.entity.Theme().Sink(tomo.PatternButton, comboBoxCase)
	margin     := element.entity.Theme().Margin(tomo.PatternButton, comboBoxCase)
	padding    := element.entity.Theme().Padding(tomo.PatternButton, comboBoxCase)
	
	offset := image.Pt(0, bounds.Dy() / 2).Add(bounds.Min)

	textBounds := element.drawer.LayoutBounds()
	offset.Y -= textBounds.Dy() / 2
	offset.Y -= textBounds.Min.Y
	offset.X -= textBounds.Min.X

	icon := element.entity.Theme().Icon(tomo.IconExpand, tomo.IconSizeSmall, comboBoxCase) 
	if icon != nil {
		iconBounds := icon.Bounds()
		addedWidth := iconBounds.Dx() + margin.X
		iconOffset := bounds.Min

		iconOffset.X += padding[3]
		iconOffset.Y =
			bounds.Min.Y +
			(bounds.Dy() -
			iconBounds.Dy()) / 2
		if element.pressed {
			iconOffset = iconOffset.Add(sink)
		}
		offset.X += addedWidth + padding[3]

		icon.Draw(destination, foreground, iconOffset)
	}

	if element.pressed {
		offset = offset.Add(sink)
	}
	element.drawer.Draw(destination, foreground, offset)
}

// OnChange sets the function to be called when this element's value is changed.
func (element *ComboBox) OnChange (callback func ()) {
	element.onChange = callback
}

// Value returns this element's value.
func (element *ComboBox) Value () Option {
	return element.selected
}

// Select sets this element's value.
func (element *ComboBox) Select (option Option) {
	element.selected = option
	element.drawer.SetText([]rune(option.Title()))
	element.updateMinimumSize()
	element.entity.Invalidate()
	if element.onChange != nil {
		element.onChange()
	}
}

// Filled returns whether this element has a value other than (None).
func (element *ComboBox) Filled () bool {
	return element.selected != ""
}

// Focus gives this element input focus.
func (element *ComboBox) Focus () {
	if !element.entity.Focused() { element.entity.Focus() }
}

// Enabled returns whether this element is enabled or not.
func (element *ComboBox) Enabled () bool {
	return element.enabled
}

// SetEnabled sets whether this element is enabled or not.
func (element *ComboBox) SetEnabled (enabled bool) {
	if element.enabled == enabled { return }
	element.enabled = enabled
	element.entity.Invalidate()
}

func (element *ComboBox) HandleThemeChange () {
	element.drawer.SetFace (element.entity.Theme().FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal,
		comboBoxCase))
	element.updateMinimumSize()
	element.entity.Invalidate()
}

func (element *ComboBox) HandleFocusChange () {
	element.entity.Invalidate()
}

func (element *ComboBox) HandleMouseDown (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) {
	if !element.Enabled() { return }
	element.Focus()
	if button != input.ButtonLeft { return }
	element.dropDown()
}

func (element *ComboBox) HandleMouseUp (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) { }

func (element *ComboBox) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
	if !element.Enabled() { return }

	selectionDelta := 0
	switch key {
	case input.KeyEnter:
		element.pressed = true
		element.entity.Invalidate()
	case input.KeyUp, input.KeyLeft:
		selectionDelta = -1
	case input.KeyDown, input.KeyRight:
		selectionDelta = 1
	}

	if selectionDelta != 0 {
		selected := 0
		for index, option := range element.options {
			if option == element.selected {
				selected = index
			}
		}
		selected += selectionDelta
		if selected < 0 {
			selected = len(element.options) - 1
		} else if selected >= len(element.options) {
			selected = 0
		}

		element.Select(element.options[selected])
	}
}

func (element *ComboBox) HandleKeyUp(key input.Key, modifiers input.Modifiers) {
	if key == input.KeyEnter && element.pressed {
		element.pressed = false
		element.entity.Invalidate()
		if !element.Enabled() { return }
		element.dropDown()
	}
}

func (element *ComboBox) dropDown () {
	window := element.entity.Window()
	menu, err := window.NewMenu(element.entity.Bounds())
	if err != nil { return }

	cellToOption := make(map[ability.Selectable] Option)

	list := NewList()
	for _, option := range element.options {
		option := option
		cell := NewCell(NewLabel(option.Title()))
		cellToOption[cell] = option
		list.Adopt(cell)

		if option == element.selected {
			list.Select(cell)
		}
	}
	list.OnClick(func () {
		selected := list.Selected()
		if selected == nil { return }
		element.Select(cellToOption[selected])
		menu.Close()
	})

	menu.Adopt(list)
	list.Focus()
	menu.Show()
}

func (element *ComboBox) updateMinimumSize () {
	padding := element.entity.Theme().Padding(tomo.PatternButton, comboBoxCase)
	margin  := element.entity.Theme().Margin(tomo.PatternButton, comboBoxCase)

	textBounds  := element.drawer.LayoutBounds()
	minimumSize := textBounds.Sub(textBounds.Min)
	
	icon := element.entity.Theme().Icon(tomo.IconExpand, tomo.IconSizeSmall, comboBoxCase) 
	if icon != nil {
		bounds := icon.Bounds()
		minimumSize.Max.X += bounds.Dx()
		minimumSize.Max.X += margin.X
	}
	
	minimumSize = padding.Inverse().Apply(minimumSize)
	element.entity.SetMinimumSize(minimumSize.Dx(), minimumSize.Dy())
}

func (element *ComboBox) state () tomo.State {
	return tomo.State {
		Disabled: !element.Enabled(),
		Focused:  element.entity.Focused(),
		Pressed:  element.pressed,
	}
}
