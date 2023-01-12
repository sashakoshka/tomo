package popups

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/layouts"

type DialogKind int

const (
	DialogKindInfo DialogKind = iota
	DialogKindQuestion
	DialogKindWarning
	DialogKindError
)

type Button struct {
	Name string
	OnPress func ()
}

func NewDialog (kind DialogKind, title, message string, buttons ...Button) {
	window, _ := tomo.NewWindow(2, 2)
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
}
