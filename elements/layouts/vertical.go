package layouts

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"

// Vertical arranges elements vertically. Elements at the start of the entry
// list will be positioned at the top, and elements at the end of the entry list
// will positioned at the bottom. All elements have the same width.
type Vertical struct {
	// If Gap is true, a gap will be placed between each element.
	Gap bool

	// If Pad is true, there will be padding running along the inside of the
	// layout's border.
	Pad bool
}

// Arrange arranges a list of entries vertically.
func (layout Vertical) Arrange (entries []tomo.LayoutEntry, width, height int) {
	if layout.Pad {
		width  -= theme.Padding() * 2
		height -= theme.Padding() * 2
	}
	freeSpace := height
	expandingElements := 0

	// count the number of expanding elements and the amount of free space
	// for them to collectively occupy
	for _, entry := range entries {
		if entry.Expand {
			expandingElements ++
		} else {
			_, entryMinHeight := entry.MinimumSize()
			freeSpace -= entryMinHeight
		}
	}
	expandingElementHeight := 0
	if expandingElements > 0 {
		expandingElementHeight = freeSpace / expandingElements
	}
	
	x, y := 0, 0
	if layout.Pad {
		x += theme.Padding()
		y += theme.Padding()
	}

	// set the size and position of each element
	for index, entry := range entries {
		if index > 0 && layout.Gap { y += theme.Padding() }
		
		entries[index].Position = image.Pt(x, y)
		entryHeight := 0
		if entry.Expand {
			entryHeight = expandingElementHeight
		} else {
			_, entryHeight = entry.MinimumSize()
		}
		y += entryHeight
		entryBounds := entry.Bounds()
		if entryBounds.Dx() != width || entryBounds.Dy() != entryHeight {
			entry.Handle (tomo.EventResize {
				Width:  width,
				Height: entryHeight,
			})
		}
	}
}

// MinimumSize returns the minimum width and height will be needed to arrange
// the given list of entries.
func (layout Vertical) MinimumSize (entries []tomo.LayoutEntry) (width, height int) {
	for index, entry := range entries {
		entryWidth, entryHeight := entry.MinimumSize()
		if entryWidth > width {
			width = entryWidth
		}
		height += entryHeight
		if layout.Gap && index > 0 {
			height += theme.Padding()
		}
	}

	if layout.Pad {
		width  += theme.Padding() * 2
		height += theme.Padding() * 2
	}
	return
}
