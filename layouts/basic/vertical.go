package basicLayouts

import "image"
import "git.tebibyte.media/sashakoshka/tomo/layouts"
import "git.tebibyte.media/sashakoshka/tomo/elements"

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
func (layout Vertical) Arrange (
	entries []layouts.LayoutEntry,
	margin int,
	bounds image.Rectangle,
) {
	if layout.Pad { bounds = bounds.Inset(margin) }

	// count the number of expanding elements and the amount of free space
	// for them to collectively occupy, while gathering minimum heights.
	freeSpace := bounds.Dy()
	minimumHeights := make([]int, len(entries))
	expandingElements := 0
	for index, entry := range entries {
		var entryMinHeight int

		if child, flexible := entry.Element.(elements.Flexible); flexible {
			entryMinHeight = child.FlexibleHeightFor(bounds.Dx())
		} else {
			_, entryMinHeight = entry.MinimumSize()
		}
		minimumHeights[index] = entryMinHeight
		
		if entry.Expand {
			expandingElements ++
		} else {
			freeSpace -= entryMinHeight
		}
		if index > 0 && layout.Gap {
			freeSpace -= margin
		}
	}
	
	expandingElementHeight := 0
	if expandingElements > 0 {
		expandingElementHeight = freeSpace / expandingElements
	}
	
	// set the size and position of each element
	dot := bounds.Min
	for index, entry := range entries {
		if index > 0 && layout.Gap { dot.Y += margin }
		
		entry.Bounds.Min = dot
		entryHeight := 0
		if entry.Expand {
			entryHeight = expandingElementHeight
		} else {
			entryHeight = minimumHeights[index]
		}
		dot.Y += entryHeight
		entryBounds := entry.Bounds
		entry.Bounds.Max = entryBounds.Min.Add(image.Pt(bounds.Dx(), entryHeight))
		entries[index] = entry
	}
}

// MinimumSize returns the minimum width and height that will be needed to
// arrange the given list of entries.
func (layout Vertical) MinimumSize (
	entries []layouts.LayoutEntry,
	margin int,
) (
	width, height int,
) {
	for index, entry := range entries {
		entryWidth, entryHeight := entry.MinimumSize()
		if entryWidth > width {
			width = entryWidth
		}
		height += entryHeight
		if layout.Gap && index > 0 {
			height += margin
		}
	}

	if layout.Pad {
		width  += margin * 2
		height += margin * 2
	}
	return
}

// FlexibleHeightFor Returns the minimum height the layout needs to lay out the
// specified elements at the given width, taking into account flexible elements.
func (layout Vertical) FlexibleHeightFor (
	entries []layouts.LayoutEntry,
	margin int,
	width int,
) (
	height int,
) {
	if layout.Pad {
		width -= margin * 2
		height += margin * 2
	}
	
	for index, entry := range entries {
		child, flexible := entry.Element.(elements.Flexible)
		if flexible {
			height += child.FlexibleHeightFor(width)
		} else {
			_, entryHeight := entry.MinimumSize()
			height += entryHeight
		}
		
		if layout.Gap && index > 0 {
			height += margin
		}
	}
	return
}