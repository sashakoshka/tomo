package layouts

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"

type Dialog struct {
	// If Gap is true, a gap will be placed between each element.
	Gap bool

	// If Pad is true, there will be padding running along the inside of the
	// layout's border.
	Pad bool
}

// Arrange arranges a list of entries into a dialog.
func (layout Dialog) Arrange (entries []tomo.LayoutEntry, width, height int) {
	if layout.Pad {
		width  -= theme.Padding() * 2
		height -= theme.Padding() * 2
	}

	controlRowWidth, controlRowHeight := 0, 0
	if len(entries) > 1 {
		controlRowWidth,
		controlRowHeight = layout.minimumSizeOf(entries[1:])
	}

	if len(entries) > 0 {
		entries[0].Position = image.Point { }
		if layout.Pad {
			entries[0].Position.X += theme.Padding()
			entries[0].Position.Y += theme.Padding()
		}
		mainHeight := height - controlRowHeight
		if layout.Gap {
			mainHeight -= theme.Padding()
		}
		mainBounds := entries[0].Bounds()
		if mainBounds.Dy() != mainHeight ||
			mainBounds.Dx() != width {
			entries[0].Handle (tomo.EventResize {
				Width:  width,
				Height: mainHeight,
			})
		}
	}

	if len(entries) > 1 {
		freeSpace := width
		expandingElements := 0

		// count the number of expanding elements and the amount of free
		// space for them to collectively occupy
		for index, entry := range entries[1:] {
			if entry.Expand {
				expandingElements ++
			} else {
				entryMinWidth, _ := entry.MinimumSize()
				freeSpace -= entryMinWidth
			}
			if index > 0 && layout.Gap {
				freeSpace -= theme.Padding()
			}
		}
		expandingElementWidth := 0
		if expandingElements > 0 {
			expandingElementWidth = freeSpace / expandingElements
		}

		// determine starting position and dimensions for control row
		x, y := 0, height - controlRowHeight
		if expandingElements == 0 {
			x = width - controlRowWidth
		}
		if layout.Pad {
			x += theme.Padding()
			y += theme.Padding()
		}
		height -= controlRowHeight

		// set the size and position of each element in the control row
		for index, entry := range entries[1:] {
			if index > 0 && layout.Gap { x += theme.Padding() }
			
			entries[index + 1].Position = image.Pt(x, y)
			entryWidth := 0
			if entry.Expand {
				entryWidth = expandingElementWidth
			} else {
				entryWidth, _ = entry.MinimumSize()
			}
			x += entryWidth
			entryBounds := entry.Bounds()
			if entryBounds.Dy() != controlRowHeight ||
				entryBounds.Dx() != entryWidth {
				entry.Handle (tomo.EventResize {
					Width:  entryWidth,
					Height: controlRowHeight,
				})
			}
		}
	}

	
}

// MinimumSize returns the minimum width and height that will be needed to
// arrange the given list of entries.
func (layout Dialog) MinimumSize (entries []tomo.LayoutEntry) (width, height int) {
	if len(entries) > 0 {
		mainChildHeight := 0
		width, mainChildHeight = entries[0].MinimumSize()
		height += mainChildHeight
	}

	if len(entries) > 1 {
		if layout.Gap { height += theme.Padding() }
		additionalWidth,
		additionalHeight := layout.minimumSizeOf(entries[1:])
		height += additionalHeight
		if additionalWidth > width {
			width = additionalWidth
		}
	}

	if layout.Pad {
		width  += theme.Padding() * 2
		height += theme.Padding() * 2
	}
	return
}

func (layout Dialog) minimumSizeOf (
	entries []tomo.LayoutEntry,
) (
	width, height int,
) {
	for index, entry := range entries {
		entryWidth, entryHeight := entry.MinimumSize()
		if entryHeight > height {
			height = entryHeight
		}
		width += entryWidth
		if layout.Gap && index > 0 {
			width += theme.Padding()
		}
	}
	return
}
