package popups

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/layouts"
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
	window tomo.Window,
) {
	window, _ = tomo.NewWindow(2, 2)
	window.SetTitle(title)
	
	container := basic.NewContainer(layouts.Dialog { true, true })
	window.Adopt(container)

	container.Adopt(basic.NewLabel(message, false), true)
	if len(buttons) == 0 {
		button := basic.NewButton("OK")
		button.OnClick(window.Close)
		container.Adopt(button, false)
		button.Select()
	} else {
		var button *basic.Button
		for _, buttonDescriptor := range buttons {
			button = basic.NewButton(buttonDescriptor.Name)
			button.SetEnabled(buttonDescriptor.OnPress != nil)
			button.OnClick (func () {
				buttonDescriptor.OnPress()
				window.Close()
			})
			container.Adopt(button, false)
		}
		button.Select()
	}
	
	window.Show()
	return
}
