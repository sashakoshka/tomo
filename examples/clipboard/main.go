package main

import "io"
import "image"
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
	window, _ := tomo.NewWindow(256, 2)
	window.SetTitle("Clipboard")

	container := containers.NewContainer(basicLayouts.Vertical { true, true })
	textInput := basicElements.NewTextBox("", "")
	controlRow := containers.NewContainer(basicLayouts.Horizontal { true, false })
	copyButton := basicElements.NewButton("Copy")
	copyButton.SetIcon(theme.IconCopy)
	pasteButton := basicElements.NewButton("Paste")
	pasteButton.SetIcon(theme.IconPaste)
	pasteImageButton := basicElements.NewButton("Image")
	pasteImageButton.SetIcon(theme.IconPictures)

	imageClipboardCallback := func (clipboard data.Data, err error) {
		if err != nil {
			popups.NewDialog (
				popups.DialogKindError,
				window,
				"Error",
				"Cannot get clipboard:\n" + err.Error())
			return 
		}

		imageData, ok := clipboard[data.M("image", "png")]
		if !ok {
			popups.NewDialog (
				popups.DialogKindError,
				window,
				"Clipboard Empty",
				"No image data in clipboard")
			return 
		}

		img, _, err := image.Decode(imageData)
		if err != nil {
			popups.NewDialog (
				popups.DialogKindError,
				window,
				"Error",
				"Cannot decode image:\n" + err.Error())
			return 
		}
		imageWindow(img)
	}
	clipboardCallback := func (clipboard data.Data, err error) {
		if err != nil {
			popups.NewDialog (
				popups.DialogKindError,
				window,
				"Error",
				"Cannot get clipboard:\n" + err.Error())
			return 
		}
		
		textData, ok := clipboard[data.MimePlain]
		if !ok {
			popups.NewDialog (
				popups.DialogKindError,
				window,
				"Clipboard Empty",
				"No text data in clipboard")
			return 
		}

		text, _ := io.ReadAll(textData)
		textInput.SetValue(string(text))
	}
	copyButton.OnClick (func () {
		window.Copy(data.Text(textInput.Value()))
	})
	pasteButton.OnClick (func () {
		window.Paste(clipboardCallback, data.MimePlain)
	})
	pasteImageButton.OnClick (func () {
		window.Paste(imageClipboardCallback, data.M("image", "png"))
	})
	
	container.Adopt(textInput, true)
	controlRow.Adopt(copyButton, true)
	controlRow.Adopt(pasteButton, true)
	controlRow.Adopt(pasteImageButton, true)
	container.Adopt(controlRow, false)
	window.Adopt(container)
		
	window.OnClose(tomo.Stop)
	window.Show()
}

func imageWindow (image image.Image) {
	window, _ := tomo.NewWindow(2, 2)
	window.SetTitle("Clipboard Image")
	container := containers.NewContainer(basicLayouts.Vertical { true, true })
	closeButton := basicElements.NewButton("Ok")
	closeButton.SetIcon(theme.IconYes)
	closeButton.OnClick(window.Close)
	
	container.Adopt(basicElements.NewImage(image), true)
	container.Adopt(closeButton, false)
	window.Adopt(container)
	window.Show()
}
