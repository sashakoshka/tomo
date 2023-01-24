package basic

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"

// ListEntry is an item that can be added to a list.
type ListEntry struct {
	drawer artist.TextDrawer
	bounds image.Rectangle
	textPoint image.Point
	text string
	forcedMinimumWidth int
	onClick func ()
}

func NewListEntry (text string, onClick func ()) (entry ListEntry) {
	entry = ListEntry  {
		text:    text,
		onClick: onClick,
	}
	entry.drawer.SetText([]rune(text))
	entry.drawer.SetFace(theme.FontFaceRegular())
	entry.updateBounds()
	return
}

func (entry *ListEntry) Collapse (width int) {
	if entry.forcedMinimumWidth == width { return }
	entry.forcedMinimumWidth = width
	entry.updateBounds()
}

func (entry *ListEntry) updateBounds () {
	padding := theme.Padding()
	
	entry.bounds = image.Rectangle { }
	entry.bounds.Max.Y = entry.drawer.LineHeight().Round() + padding
	if entry.forcedMinimumWidth > 0 {
		entry.bounds.Max.X = entry.forcedMinimumWidth
	} else {
		entry.bounds.Max.X =
			entry.drawer.LayoutBounds().Dx() + padding * 2
	}
	
	entry.textPoint =
		image.Pt(padding, padding / 2).
		Sub(entry.drawer.LayoutBounds().Min)
}

func (entry *ListEntry) Draw (
	destination tomo.Canvas,
	offset image.Point,
	selected bool,
) (
	updatedRegion image.Rectangle,
) {
	return entry.drawer.Draw (
		destination,
		theme.ForegroundPattern(true),
		offset.Add(entry.textPoint))
}

func (entry *ListEntry) Bounds () (bounds image.Rectangle) {
	return entry.bounds
}
