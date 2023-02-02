package popups

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import "git.tebibyte.media/sashakoshka/tomo/layouts/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"

// DialogKind defines the semantic role of a dialog window.
type DialogKind int

const (
	DialogKindInfo DialogKind = iota
	DialogKindQuestion
	DialogKindWarning
	DialogKindError
)

// Button represents a dialog response button.
type Button struct {
	// Name contains the text to display on the button.
	Name string

	// OnPress specifies a callback to run when the button is pressed. If
	// this callback is nil, the button will appear disabled.
	OnPress func ()
}

// NewDialog creates a new dialog window and returns it.
func NewDialog (
	kind DialogKind,
	title, message string,
	buttons ...Button,
) (
	window elements.Window,
) {
	window, _ = tomo.NewWindow(2, 2)
	window.SetTitle(title)
	
	container := basicElements.NewContainer(basicLayouts.Dialog { true, true })
	window.Adopt(container)

	container.Adopt(basicElements.NewLabel(message, false), true)
	if len(buttons) == 0 {
		button := basicElements.NewButton("OK")
		button.OnClick(window.Close)
		container.Adopt(button, false)
		button.Focus()
	} else {
		var button *basicElements.Button
		for _, buttonDescriptor := range buttons {
			button = basicElements.NewButton(buttonDescriptor.Name)
			button.SetEnabled(buttonDescriptor.OnPress != nil)
			button.OnClick (func () {
				buttonDescriptor.OnPress()
				window.Close()
			})
			container.Adopt(button, false)
		}
		button.Focus()
	}
	
	window.Show()
	return
}
