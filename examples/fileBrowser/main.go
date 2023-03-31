package main

import "os"
import "path/filepath"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/layouts"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import "git.tebibyte.media/sashakoshka/tomo/elements/file"
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(384, 384)
	window.SetTitle("File browser")
	container := containers.NewContainer(layouts.Vertical { true, true })
	window.Adopt(container)
	homeDir, _ := os.UserHomeDir()

	controlBar := containers.NewContainer(layouts.Horizontal { })
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
	
	statusBar := containers.NewContainer(layouts.Horizontal { true, false })
	directory, _ := fileElements.NewFile(homeDir, nil)
	baseName := elements.NewLabel(filepath.Base(homeDir), false)
	
	scrollContainer  := containers.NewScrollContainer(false, true)
	directoryView, _ := fileElements.NewDirectory(homeDir, nil)
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
	
	controlBar.Adopt(backButton,    false)
	controlBar.Adopt(forwardButton, false)
	controlBar.Adopt(refreshButton, false)
	controlBar.Adopt(upwardButton,  false)
	controlBar.Adopt(locationInput, true)
	scrollContainer.Adopt(directoryView)
	statusBar.Adopt(directory, false)
	statusBar.Adopt(baseName, false)
	
	container.Adopt(controlBar,      false)
	container.Adopt(scrollContainer, true)
	container.Adopt(statusBar,       false)

	window.OnClose(tomo.Stop)
	window.Show()
}
