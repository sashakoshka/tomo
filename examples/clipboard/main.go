package main

import "io"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/data"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/popups"
import "git.tebibyte.media/sashakoshka/tomo/layouts/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("Clipboard")

	container := containers.NewContainer(basicLayouts.Vertical { true, true })
	textInput := basicElements.NewTextBox("", "")
	controlRow := containers.NewContainer(basicLayouts.Horizontal { true, false })
	copyButton := basicElements.NewButton("Copy")
	copyButton.SetIcon(theme.IconCopy)
	pasteButton := basicElements.NewButton("Paste")
	pasteButton.SetIcon(theme.IconPaste)

	clipboardCallback := func (clipboard io.Reader, err error) {
		if err != nil {
			popups.NewDialog (
				popups.DialogKindError,
				window,
				"Error",
				"Cannot get clipboard:\n" + err.Error())
			return 
		}
		
		if clipboard == nil {
			popups.NewDialog (
				popups.DialogKindError,
				window,
				"Clipboard Empty",
				"No text data in clipboard")
			return 
		}

		text, _ := io.ReadAll(clipboard)
		tomo.Do (func () {
			textInput.SetValue(string(text))
		})
	}
	copyButton.OnClick (func () {
		window.Copy(data.Text(textInput.Value()))
	})
	pasteButton.OnClick (func () {
		window.Paste(data.MimePlain, clipboardCallback)
	})
	
	container.Adopt(textInput, true)
	controlRow.Adopt(copyButton, true)
	controlRow.Adopt(pasteButton, true)
	container.Adopt(controlRow, false)
	window.Adopt(container)
		
	window.OnClose(tomo.Stop)
	window.Show()
}
