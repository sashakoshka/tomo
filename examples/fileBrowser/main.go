package main

import "os"
import "path/filepath"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/layouts/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/file"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"
import "git.tebibyte.media/sashakoshka/tomo/elements/containers"
import _ "git.tebibyte.media/sashakoshka/tomo/backends/all"

func main () {
	tomo.Run(run)
}

func run () {
	window, _ := tomo.NewWindow(384, 384)
	window.SetTitle("File browser")
	container := containers.NewContainer(basicLayouts.Vertical { true, true })
	window.Adopt(container)
	homeDir, _ := os.UserHomeDir()

	controlBar := containers.NewContainer(basicLayouts.Horizontal { })
	backButton := basicElements.NewButton("Back")
	backButton.SetIcon(theme.IconBackward)
	backButton.ShowText(false)
	forwardButton := basicElements.NewButton("Forward")
	forwardButton.SetIcon(theme.IconForward)
	forwardButton.ShowText(false)
	refreshButton := basicElements.NewButton("Refresh")
	refreshButton.SetIcon(theme.IconRefresh)
	refreshButton.ShowText(false)
	upwardButton := basicElements.NewButton("Go Up")
	upwardButton.SetIcon(theme.IconUpward)
	upwardButton.ShowText(false)
	locationInput := basicElements.NewTextBox("Location", "")
	
	statusBar := containers.NewContainer(basicLayouts.Horizontal { true, false })
	directory, _ := fileElements.NewFile(homeDir, nil)
	baseName := basicElements.NewLabel(filepath.Base(homeDir), false)
	
	scrollContainer  := containers.NewScrollContainer(false, true)
	directoryView, _ := fileElements.NewDirectoryView(homeDir, nil)
	choose := func (filePath string) {
		directoryView.SetLocation(filePath, nil)
		directory.SetLocation(filePath, nil)
		locationInput.SetValue(filePath)
		baseName.SetText(filepath.Base(filePath))
	}
	directoryView.OnChoose(choose)
	locationInput.OnEnter (func () {
		choose(locationInput.Value())
	})
	choose(homeDir)
	
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
