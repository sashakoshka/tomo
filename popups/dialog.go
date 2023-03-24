package popups

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import "git.tebibyte.media/sashakoshka/tomo/layouts/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"

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

// NewDialog creates a new modal dialog window and returns it. If parent is nil,
// the dialog will just be a normal window
func NewDialog (
	kind DialogKind,
	parent elements.Window,
	title, message string,
	buttons ...Button,
) (
	window elements.Window,
) {
	if parent == nil {
		window, _ = tomo.NewWindow(2, 2)
	} else {
		window, _ = parent.NewModal(2, 2)
	}
	window.SetTitle(title)

	container := containers.NewContainer(basicLayouts.Dialog { true, true })
	window.Adopt(container)

	messageContainer := containers.NewContainer(basicLayouts.Horizontal { true, false })
	iconId := theme.IconInformation
	switch kind {
	case DialogKindInfo:     iconId = theme.IconInformation
	case DialogKindQuestion: iconId = theme.IconQuestion
	case DialogKindWarning:  iconId = theme.IconWarning
	case DialogKindError:    iconId = theme.IconError
	}
	
	messageContainer.Adopt(basicElements.NewIcon(iconId, theme.IconSizeLarge), false)
	messageContainer.Adopt(basicElements.NewLabel(message, false), true)
	container.Adopt(messageContainer, true)
	
	if len(buttons) == 0 {
		button := basicElements.NewButton("OK")
		button.SetIcon(theme.IconYes)
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
