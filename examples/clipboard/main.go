package main

import "io"
import "image"
import _ "image/png"
import _ "image/gif"
import _ "image/jpeg"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/data"
import "git.tebibyte.media/sashakoshka/tomo/nasin"
import "git.tebibyte.media/sashakoshka/tomo/popups"
import "git.tebibyte.media/sashakoshka/tomo/elements"

var validImageTypes = []data.Mime {
	data.M("image", "png"),
	data.M("image", "gif"),
	data.M("image", "jpeg"),
}


func main () {
	nasin.Run(Application { })
}

type Application struct { }

func (Application) Init () error {
	window, err:= nasin.NewWindow(tomo.Bounds(0, 0, 256, 0))
	if err != nil { return err }
	window.SetTitle("Clipboard")

	container := elements.NewVBox(elements.SpaceBoth)
	textInput := elements.NewTextBox("", "")
	controlRow := elements.NewHBox(elements.SpaceMargin)
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
	
	container.AdoptExpand(textInput)
	controlRow.AdoptExpand(copyButton)
	controlRow.AdoptExpand(pasteButton)
	controlRow.AdoptExpand(pasteImageButton)
	container.Adopt(controlRow)
	window.Adopt(container)
		
	window.OnClose(nasin.Stop)
	window.Show()
	return nil
}

func imageWindow (parent tomo.Window, image image.Image) {
	window, _ := parent.NewModal(tomo.Bounds(0, 0, 0, 0))
	window.SetTitle("Clipboard Image")
	container := elements.NewVBox(elements.SpaceBoth)
	closeButton := elements.NewButton("Ok")
	closeButton.SetIcon(tomo.IconYes)
	closeButton.OnClick(window.Close)
	
	container.AdoptExpand(elements.NewImage(image))
	container.Adopt(closeButton)
	window.Adopt(container)

	closeButton.Focus()
	window.Show()
}
