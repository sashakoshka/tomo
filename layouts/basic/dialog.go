package basicLayouts

import "image"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/layouts"

// Dialog arranges elements in the form of a dialog box. The first element is
// positioned above as the main focus of the dialog, and is set to expand
// regardless of whether it is expanding or not. The remaining elements are
// arranged at the bottom in a row called the control row, which is aligned to
// the right, the last element being the rightmost one.
type Dialog struct {
	// If Mergin is true, a margin will be placed between each element.
	Gap bool

	// If Pad is true, there will be padding running along the inside of the
	// layout's border.
	Pad bool
}

// Arrange arranges a list of entries into a dialog.
func (layout Dialog) Arrange (
	entries []layouts.LayoutEntry,
	margin  image.Point,
	padding artist.Inset,
	bounds image.Rectangle,
) {
	if layout.Pad { bounds = padding.Apply(bounds) }
	
	controlRowWidth, controlRowHeight := 0, 0
	if len(entries) > 1 {
		controlRowWidth,
		controlRowHeight = layout.minimumSizeOfControlRow (
			entries[1:], margin, padding)
	}

	if len(entries) > 0 {
		main := entries[0]
		main.Bounds.Min = bounds.Min
		mainHeight := bounds.Dy() - controlRowHeight
		if layout.Gap {
			mainHeight -= margin.Y
		}
		main.Bounds.Max = main.Bounds.Min.Add(image.Pt(bounds.Dx(), mainHeight))
		entries[0] = main
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
				freeSpace -= margin.X
			}
		}
		expandingElementWidth := 0
		if expandingElements > 0 {
			expandingElementWidth = freeSpace / expandingElements
		}

		// determine starting position and dimensions for control row
		dot := image.Pt(bounds.Min.X, bounds.Max.Y - controlRowHeight)
		if expandingElements == 0 {
			dot.X = bounds.Max.X - controlRowWidth
		}

		// set the size and position of each element in the control row
		for index, entry := range entries[1:] {
			if index > 0 && layout.Gap { dot.X += margin.X }
			
			entry.Bounds.Min = dot
			entryWidth := 0
			if entry.Expand {
				entryWidth = expandingElementWidth
			} else {
				entryWidth, _ = entry.MinimumSize()
			}
			dot.X += entryWidth
			entryBounds := entry.Bounds
			if entryBounds.Dy() != controlRowHeight ||
				entryBounds.Dx() != entryWidth {
				entry.Bounds.Max = entryBounds.Min.Add (
					image.Pt(entryWidth, controlRowHeight))
			}
			entries[index + 1] = entry
		}
	}

	
}

// MinimumSize returns the minimum width and height that will be needed to
// arrange the given list of entries.
func (layout Dialog) MinimumSize (
	entries []layouts.LayoutEntry,
	margin  image.Point,
	padding artist.Inset,
) (
	width, height int,
) {
	if len(entries) > 0 {
		mainChildHeight := 0
		width, mainChildHeight = entries[0].MinimumSize()
		height += mainChildHeight
	}

	if len(entries) > 1 {
		if layout.Gap { height += margin.X }
		additionalWidth,
		additionalHeight := layout.minimumSizeOfControlRow (
			entries[1:], margin, padding)
		height += additionalHeight
		if additionalWidth > width {
			width = additionalWidth
		}
	}

	if layout.Pad {
		width  += padding.Horizontal()
		height += padding.Vertical()
	}
	return
}

func (layout Dialog) minimumSizeOfControlRow (
	entries []layouts.LayoutEntry,
	margin  image.Point,
	padding artist.Inset,
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
			width += margin.X
		}
	}
	return
}
