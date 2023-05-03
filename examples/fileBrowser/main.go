package main

import "os"
import "path/filepath"
import "tomo"
import "tomo/nasin"
import "tomo/elements"

func main () {
	nasin.Run(Application { })
}

type Application struct { }

func (Application) Init () error {
	window, err := nasin.NewWindow(tomo.Bounds(0, 0, 384, 384))
	if err != nil { return err }
	window.SetTitle("File browser")
	container := elements.NewVBox(elements.SpaceBoth)
	window.Adopt(container)
	homeDir, err := os.UserHomeDir()
	if err != nil { return err }

	controlBar := elements.NewHBox(elements.SpaceNone)
	backButton := elements.NewButton("Back")
	backButton.SetIcon(tomo.IconBackward)
	backButton.ShowText(false)
	forwardButton := elements.NewButton("Forward")
	forwardButton.SetIcon(tomo.IconForward)
	forwardButton.ShowText(false)
	refreshButton := elements.NewButton("Refresh")
	refreshButton.SetIcon(tomo.IconRefresh)
	refreshButton.ShowText(false)
	upwardButton := elements.NewButton("Go Up")
	upwardButton.SetIcon(tomo.IconUpward)
	upwardButton.ShowText(false)
	locationInput := elements.NewTextBox("Location", "")
	
	statusBar := elements.NewHBox(elements.SpaceMargin)
	directory, _ := elements.NewFile(homeDir, nil)
	baseName := elements.NewLabel(filepath.Base(homeDir))
	
	directoryView, _ := elements.NewDirectory(homeDir, nil)
	updateStatus := func () {
		filePath, _ := directoryView.Location()
		directory.SetLocation(filePath, nil)
		locationInput.SetValue(filePath)
		baseName.SetText(filepath.Base(filePath))
	}
	choose := func (filePath string) {
		directoryView.SetLocation(filePath, nil)
		updateStatus()
	}
	directoryView.OnChoose(choose)
	locationInput.OnEnter (func () {
		choose(locationInput.Value())
	})
	choose(homeDir)
	backButton.OnClick (func () {
		directoryView.Backward()
		updateStatus()
	})
	forwardButton.OnClick (func () {
		directoryView.Forward()
		updateStatus()
	})
	refreshButton.OnClick (func () {
		directoryView.Update()
		updateStatus()
	})
	upwardButton.OnClick (func () {
		filePath, _ := directoryView.Location()
		choose(filepath.Dir(filePath))
	})

	controlBar.Adopt(backButton, forwardButton, refreshButton, upwardButton)
	controlBar.AdoptExpand(locationInput)
	statusBar.Adopt(directory, baseName)
	
	container.Adopt(controlBar)
	container.AdoptExpand (
		elements.NewScroll(elements.ScrollVertical, directoryView))
	container.Adopt(statusBar)

	window.OnClose(nasin.Stop)
	window.Show()
	return nil
}
