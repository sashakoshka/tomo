package popups

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"

// DialogKind defines the semantic role of a dialog window.
type DialogKind int

const (
	DialogKindInfo DialogKind = iota
	DialogKindQuestion
	DialogKindWarning
	DialogKindError
)

// TODO: add ability to have an icon for buttons

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
	parent tomo.Window,
	title, message string,
	buttons ...Button,
) (
	window tomo.Window,
) {
	if parent == nil {
		window, _ = tomo.NewWindow(image.Rectangle { })
	} else {
		window, _ = parent.NewModal(image.Rectangle { })
	}
	window.SetTitle(title)
	
	box        := containers.NewVBox(true,  true)
	messageRow := containers.NewHBox(false, true)
	controlRow := containers.NewHBox(false, true)

	iconId := tomo.IconInformation
	switch kind {
	case DialogKindInfo:     iconId = tomo.IconInformation
	case DialogKindQuestion: iconId = tomo.IconQuestion
	case DialogKindWarning:  iconId = tomo.IconWarning
	case DialogKindError:    iconId = tomo.IconError
	}
	
	messageRow.Adopt(elements.NewIcon(iconId, tomo.IconSizeLarge), false)
	messageRow.Adopt(elements.NewLabel(message, false), true)
	
	controlRow.Adopt(elements.NewSpacer(false), true)
	box.Adopt(messageRow, true)
	box.Adopt(controlRow, false)
	window.Adopt(box)
	
	if len(buttons) == 0 {
		button := elements.NewButton("OK")
		button.SetIcon(tomo.IconYes)
		button.OnClick(window.Close)
		controlRow.Adopt(button, false)
		button.Focus()
	} else {
		var button *elements.Button
		for _, buttonDescriptor := range buttons {
			button = elements.NewButton(buttonDescriptor.Name)
			button.SetEnabled(buttonDescriptor.OnPress != nil)
			button.OnClick (func () {
				buttonDescriptor.OnPress()
				window.Close()
			})
			controlRow.Adopt(button, false)
		}
		button.Focus()
	}
	
	window.Show()
	return
}
