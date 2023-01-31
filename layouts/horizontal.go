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
func (layout Horizontal) Arrange (entries []tomo.LayoutEntry, bounds image.Rectangle) {
	if layout.Pad { bounds = bounds.Inset(theme.Margin()) }
	
	// get width of expanding elements
	expandingElementWidth := layout.expandingElementWidth(entries, bounds.Dx())

	// set the size and position of each element
	dot := bounds.Min
	for index, entry := range entries {
		if index > 0 && layout.Gap { dot.X += theme.Margin() }
		
		entry.Bounds.Min = dot
		entryWidth := 0
		if entry.Expand {
			entryWidth = expandingElementWidth
		} else {
			entryWidth, _ = entry.MinimumSize()
		}
		dot.X += entryWidth
		entryBounds := entry.Bounds
		if entryBounds.Dy() != bounds.Dy() || entryBounds.Dx() != entryWidth {
			entry.Bounds.Max = entryBounds.Min.Add (
				image.Pt(entryWidth, bounds.Dy()))
		}

		entries[index] = entry
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
			width += theme.Margin()
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
func (layout Horizontal) FlexibleHeightFor (
	entries []tomo.LayoutEntry,
	width int,
) (
	height int,
) {
	if layout.Pad {
		width -= theme.Margin() * 2
	}
	// get width of expanding elements
	expandingElementWidth := layout.expandingElementWidth(entries, width)
	
	x, y := 0, 0
	if layout.Pad {
		x += theme.Margin()
		y += theme.Margin()
	}

	// set the size and position of each element
	for index, entry := range entries {
		entryWidth, entryHeight := entry.MinimumSize()
		if entry.Expand {
			entryWidth = expandingElementWidth
		}
		if child, flexible := entry.Element.(tomo.Flexible); flexible {
			entryHeight = child.FlexibleHeightFor(entryWidth)
		}
		if entryHeight > height { height = entryHeight }
		
		x += entryWidth
		if index > 0 && layout.Gap { x += theme.Margin() }
	}

	if layout.Pad {
		height += theme.Margin() * 2
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
			freeSpace -= theme.Margin()
		}
	}
	
	if expandingElements > 0 {
		width = freeSpace / expandingElements
	}
	return
}
