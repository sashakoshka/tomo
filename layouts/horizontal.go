package layouts

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"

// Horizontal arranges elements horizontally. Elements at the start of the entry
// list will be positioned on the left, and elements at the end of the entry
// list will positioned on the right. All elements have the same height.
type Horizontal struct {
	// If Gap is true, a gap will be placed between each element.
	Gap bool

	// If Pad is true, there will be padding running along the inside of the
	// layout's border.
	Pad bool
}

// Arrange arranges a list of entries horizontally.
func (layout Horizontal) Arrange (entries []tomo.LayoutEntry, width, height int) {
	if layout.Pad {
		width  -= theme.Padding() * 2
		height -= theme.Padding() * 2
	}
	// get width of expanding elements
	expandingElementWidth := layout.expandingElementWidth(entries, width)
	
	x, y := 0, 0
	if layout.Pad {
		x += theme.Padding()
		y += theme.Padding()
	}

	// set the size and position of each element
	for index, entry := range entries {
		if index > 0 && layout.Gap { x += theme.Padding() }
		
		entries[index].Position = image.Pt(x, y)
		entryWidth := 0
		if entry.Expand {
			entryWidth = expandingElementWidth
		} else {
			entryWidth, _ = entry.MinimumSize()
		}
		x += entryWidth
		entryBounds := entry.Bounds()
		if entryBounds.Dy() != height || entryBounds.Dx() != entryWidth {
			entry.Resize(entryWidth, height)
		}
	}
}

// MinimumSize returns the minimum width and height that will be needed to
// arrange the given list of entries.
func (layout Horizontal) MinimumSize (
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

	if layout.Pad {
		width  += theme.Padding() * 2
		height += theme.Padding() * 2
	}
	return
}

func (layout Horizontal) MinimumHeightFor (
	entries []tomo.LayoutEntry,
	width int,
) (
	height int,
) {
	if layout.Pad {
		width -= theme.Padding() * 2
	}
	// get width of expanding elements
	expandingElementWidth := layout.expandingElementWidth(entries, width)
	
	x, y := 0, 0
	if layout.Pad {
		x += theme.Padding()
		y += theme.Padding()
	}

	// set the size and position of each element
	for index, entry := range entries {
		entryWidth, entryHeight := entry.MinimumSize()
		if entry.Expand {
			entryWidth = expandingElementWidth
		}
		if child, flexible := entry.Element.(tomo.Flexible); flexible {
			entryHeight = child.MinimumHeightFor(entryWidth)
		}
		if entryHeight > height { height = entryHeight }
		
		x += entryWidth
		if index > 0 && layout.Gap { x += theme.Padding() }
	}

	if layout.Pad {
		height += theme.Padding() * 2
	}
	return
}

func (layout Horizontal) expandingElementWidth (
	entries []tomo.LayoutEntry,
	freeSpace int,
) (
	width int,
) {
	expandingElements := 0

	// count the number of expanding elements and the amount of free space
	// for them to collectively occupy
	for index, entry := range entries {
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
	
	if expandingElements > 0 {
		width = freeSpace / expandingElements
	}
	return
}
