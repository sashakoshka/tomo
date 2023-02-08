package basicElements

import "image"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/artist"

// ListEntry is an item that can be added to a list.
type ListEntry struct {
	drawer artist.TextDrawer
	bounds image.Rectangle
	textPoint image.Point
	text string
	forcedMinimumWidth int
	
	theme  theme.Theme
	config config.Config
	c theme.Case
	
	onSelect func ()
}

func NewListEntry (text string, onSelect func ()) (entry ListEntry) {
	entry = ListEntry  {
		text:     text,
		onSelect: onSelect,
		c: theme.C("basic", "listEntry"),
	}
	entry.drawer.SetText([]rune(text))
	entry.updateBounds()
	return
}

func (entry *ListEntry) Collapse (width int) {
	if entry.forcedMinimumWidth == width { return }
	entry.forcedMinimumWidth = width
	entry.updateBounds()
}

func (entry *ListEntry) SetTheme (new theme.Theme) {
	entry.theme = new
	entry.drawer.SetFace (entry.theme.FontFace (
		theme.FontStyleRegular,
		theme.FontSizeNormal,
		entry.c))
	entry.updateBounds()
}

func (entry *ListEntry) SetConfig (config config.Config) {
	entry.config = config
}

func (entry *ListEntry) updateBounds () {
	entry.bounds = image.Rectangle { }
	entry.bounds.Max.Y = entry.drawer.LineHeight().Round()
	if entry.forcedMinimumWidth > 0 {
		entry.bounds.Max.X = entry.forcedMinimumWidth
	} else {
		entry.bounds.Max.X = entry.drawer.LayoutBounds().Dx()
	}
	
	inset := entry.theme.Inset(theme.PatternRaised, entry.c)
	entry.bounds.Max.Y += inset[0] + inset[2]
	
	entry.textPoint =
		image.Pt(inset[3], inset[0]).
		Sub(entry.drawer.LayoutBounds().Min)
}

func (entry *ListEntry) Draw (
	destination canvas.Canvas,
	offset image.Point,
	focused bool,
	on bool,
) (
	updatedRegion image.Rectangle,
) {
	state := theme.PatternState {
		Focused: focused,
		On: on,
	}
	pattern := entry.theme.Pattern (theme.PatternRaised, entry.c, state)
	artist.FillRectangle (
		destination,
		pattern,
		entry.Bounds().Add(offset))
	foreground := entry.theme.Pattern (theme.PatternForeground, entry.c, state)
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
