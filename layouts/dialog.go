package layouts

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"

// Dialog arranges elements in the form of a dialog box. The first element is
// positioned above as the main focus of the dialog, and is set to expand
// regardless of whether it is expanding or not. The remaining elements are
// arranged at the bottom in a row called the control row, which is aligned to
// the right, the last element being the rightmost one.
type Dialog struct {
	// If Gap is true, a gap will be placed between each element.
	Gap bool

	// If Pad is true, there will be padding running along the inside of the
	// layout's border.
	Pad bool
}

// FIXME

// Arrange arranges a list of entries into a dialog.
func (layout Dialog) Arrange (entries []tomo.LayoutEntry, bounds image.Rectangle) {
	if layout.Pad { bounds = bounds.Inset(theme.Margin()) }
	
	controlRowWidth, controlRowHeight := 0, 0
	if len(entries) > 1 {
		controlRowWidth,
		controlRowHeight = layout.minimumSizeOfControlRow(entries[1:])
	}

	if len(entries) > 0 {
		entries[0].Bounds.Min = image.Point { }
		if layout.Pad {
			entries[0].Bounds.Min.X += theme.Margin()
			entries[0].Bounds.Min.Y += theme.Margin()
		}
		mainHeight := bounds.Dy() - controlRowHeight
		if layout.Gap {
			mainHeight -= theme.Margin()
		}
		mainBounds := entries[0].Bounds
		if mainBounds.Dy() != mainHeight || mainBounds.Dx() != bounds.Dx() {
			entries[0].Bounds.Max =
				mainBounds.Min.Add(image.Pt(bounds.Dx(), mainHeight))
		}
	}

	if len(entries) > 1 {
		freeSpace := bounds.Dx()
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
				freeSpace -= theme.Margin()
			}
		}
		expandingElementWidth := 0
		if expandingElements > 0 {
			expandingElementWidth = freeSpace / expandingElements
		}

		// determine starting position and dimensions for control row
		x, y := 0, bounds.Dy() - controlRowHeight
		if expandingElements == 0 {
			x = bounds.Dx() - controlRowWidth
		}
		if layout.Pad {
			x += theme.Margin()
			y += theme.Margin()
		}
		bounds.Max.Y -= controlRowHeight

		// set the size and position of each element in the control row
		for index, entry := range entries[1:] {
			if index > 0 && layout.Gap { x += theme.Margin() }
			
			entries[index + 1].Bounds.Min = image.Pt(x, y)
			entryWidth := 0
			if entry.Expand {
				entryWidth = expandingElementWidth
			} else {
				entryWidth, _ = entry.MinimumSize()
			}
			x += entryWidth
			entryBounds := entry.Bounds
			if entryBounds.Dy() != controlRowHeight ||
				entryBounds.Dx() != entryWidth {
				entries[index].Bounds.Max = entryBounds.Min.Add (
					image.Pt(entryWidth, controlRowHeight))
			}
		}
	}

	
}

// MinimumSize returns the minimum width and height that will be needed to
// arrange the given list of entries.
func (layout Dialog) MinimumSize (
	entries []tomo.LayoutEntry,
) (
	width, height int,
) {
	if len(entries) > 0 {
		mainChildHeight := 0
		width, mainChildHeight = entries[0].MinimumSize()
		height += mainChildHeight
	}

	if len(entries) > 1 {
		if layout.Gap { height += theme.Margin() }
		additionalWidth,
		additionalHeight := layout.minimumSizeOfControlRow(entries[1:])
		height += additionalHeight
		if additionalWidth > width {
			width = additionalWidth
		}
	}

	if layout.Pad {
		width  += theme.Margin() * 2
		height += theme.Margin() * 2
	}
	return
}

// FlexibleHeightFor Returns the minimum height the layout needs to lay out the
// specified elements at the given width, taking into account flexible elements.
func (layout Dialog) FlexibleHeightFor (
	entries []tomo.LayoutEntry,
	width int,
) (
	height int,
) {
	if layout.Pad {
		width -= theme.Margin() * 2
	}
	
	if len(entries) > 0 {
		mainChildHeight := 0
		if child, flexible := entries[0].Element.(tomo.Flexible); flexible {
			mainChildHeight = child.FlexibleHeightFor(width)
		} else {
			_, mainChildHeight = entries[0].MinimumSize()
		}
		height += mainChildHeight
	}

	if len(entries) > 1 {
		if layout.Gap { height += theme.Margin() }
		_, additionalHeight := layout.minimumSizeOfControlRow(entries[1:])
		height += additionalHeight
	}

	if layout.Pad {
		height += theme.Margin() * 2
	}
	return
}

// TODO: possibly flatten this method to account for flexible elements within
// the control row.
func (layout Dialog) minimumSizeOfControlRow (
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
			width += theme.Margin()
		}
	}
	return
}
