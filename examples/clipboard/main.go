package main

import "io"
import "image"
import _ "image/png"
import _ "image/gif"
import _ "image/jpeg"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/data"
import "git.tebibyte.media/sashakoshka/tomo/popups"
import "git.tebibyte.media/sashakoshka/tomo/layouts"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"

func main () {
	tomo.Run(run)
}

var validImageTypes = []data.Mime {
	data.M("image", "png"),
	data.M("image", "gif"),
	data.M("image", "jpeg"),
}

func run () {
	window, _ := tomo.NewWindow(tomo.Bounds(0, 0, 256, 0))
	window.SetTitle("Clipboard")

	container := containers.NewContainer(layouts.Vertical { true, true })
	textInput := elements.NewTextBox("", "")
	controlRow := containers.NewContainer(layouts.Horizontal { true, false })
	copyButton := elements.NewButton("Copy")
	copyButton.SetIcon(tomo.IconCopy)
	pasteButton := elements.NewButton("Paste")
	pasteButton.SetIcon(tomo.IconPaste)
	pasteImageButton := elements.NewButton("Image")
	pasteImageButton.SetIcon(tomo.IconPictures)

	imageClipboardCallback := func (clipboard data.Data, err error) {
		if err != nil {
			popups.NewDialog (
				popups.DialogKindError,
				window,
				"Error",
				"Cannot get clipboard:\n" + err.Error())
			return 
		}

		var imageData io.Reader
		var ok bool
		for mime, reader := range clipboard {
		for _, mimeCheck := range validImageTypes {
		if mime == mimeCheck {
			imageData = reader
			ok = true
		}}}
		
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
		imageWindow(window, img)
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
		window.Paste(imageClipboardCallback, validImageTypes...)
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

func imageWindow (parent tomo.Window, image image.Image) {
	window, _ := parent.NewModal(tomo.Bounds(0, 0, 0, 0))
	window.SetTitle("Clipboard Image")
	container := containers.NewContainer(layouts.Vertical { true, true })
	closeButton := elements.NewButton("Ok")
	closeButton.SetIcon(tomo.IconYes)
	closeButton.OnClick(window.Close)
	
	container.Adopt(elements.NewImage(image), true)
	container.Adopt(closeButton, false)
	window.Adopt(container)
	window.Show()
}
