package main

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/nasin"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import "git.tebibyte.media/sashakoshka/tomo/elements/testing"

func main () {
	nasin.Run(Application { })
}

type Application struct { }

func (Application) Init () error {
	window, err := nasin.NewWindow(tomo.Bounds(0, 0, 128, 128))
	if err != nil { return err }
	window.SetTitle("vertical stack")

	container := elements.NewVBox(elements.SpaceBoth)

	label    := elements.NewLabelWrapped("it is a label hehe")
	button   := elements.NewButton("drawing pad")
	okButton := elements.NewButton("OK")
	button.OnClick (func () {
		container.DisownAll()
		container.Adopt(elements.NewLabel("Draw here (not really):"))
		container.AdoptExpand(testing.NewMouse())
		container.Adopt(okButton)
		okButton.Focus()
	})
	okButton.OnClick(nasin.Stop)

	container.AdoptExpand(label)
	container.Adopt(button, okButton)
	window.Adopt(container)
	
	okButton.Focus()
	window.OnClose(nasin.Stop)
	window.Show()
	
	return nil
}
