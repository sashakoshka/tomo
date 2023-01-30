package basic

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"

var listEntryCase = theme.C("basic", "listEntry")

// ListEntry is an item that can be added to a list.
type ListEntry struct {
	drawer artist.TextDrawer
	bounds image.Rectangle
	textPoint image.Point
	text string
	forcedMinimumWidth int
	onSelect func ()
}

func NewListEntry (text string, onSelect func ()) (entry ListEntry) {
	entry = ListEntry  {
		text:     text,
		onSelect: onSelect,
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
	entry.bounds = image.Rectangle { }
	entry.bounds.Max.Y = entry.drawer.LineHeight().Round()
	if entry.forcedMinimumWidth > 0 {
		entry.bounds.Max.X = entry.forcedMinimumWidth
	} else {
		entry.bounds.Max.X = entry.drawer.LayoutBounds().Dx()
	}
	
	_, inset := theme.ItemPattern(theme.PatternState {
	})
	entry.bounds.Max.Y += inset[0] + inset[2]
	
	entry.textPoint =
		image.Pt(inset[3], inset[0]).
		Sub(entry.drawer.LayoutBounds().Min)
}

func (entry *ListEntry) Draw (
	destination tomo.Canvas,
	offset image.Point,
	focused bool,
	on bool,
) (
	updatedRegion image.Rectangle,
) {
	pattern, _ := theme.ItemPattern(theme.PatternState {
		Case: listEntryCase,
		Focused: focused,
		On: on,
	})
	artist.FillRectangle (
		destination,
		pattern,
		entry.Bounds().Add(offset))
	foreground, _ := theme.ForegroundPattern (theme.PatternState {
		Case: listEntryCase,
		Focused: focused,
		On: on,
	})
	return entry.drawer.Draw (
		destination,
		foreground,
		offset.Add(entry.textPoint))
}

func (entry *ListEntry) RunSelect () {
	if entry.onSelect != nil {
		entry.onSelect()
	}
}

func (entry *ListEntry) Bounds () (bounds image.Rectangle) {
	return entry.bounds
}
