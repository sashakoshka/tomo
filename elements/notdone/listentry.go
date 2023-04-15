package elements

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/textdraw"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"
import "git.tebibyte.media/sashakoshka/tomo/default/config"

// ListEntry is an item that can be added to a list.
type ListEntry struct {
	drawer textdraw.Drawer
	bounds image.Rectangle
	text string
	width int
	minimumWidth int
	
	config config.Wrapped
	theme  theme.Wrapped
	
	onSelect func ()
}

func NewListEntry (text string, onSelect func ()) (entry ListEntry) {
	entry = ListEntry  {
		text:     text,
		onSelect: onSelect,
	}
	entry.theme.Case = tomo.C("tomo", "listEntry")
	entry.drawer.SetFace (entry.theme.FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal))
	entry.drawer.SetText([]rune(text))
	entry.updateBounds()
	return
}

func (entry *ListEntry) SetTheme (new tomo.Theme) {
	if new == entry.theme.Theme { return }
	entry.theme.Theme = new
	entry.drawer.SetFace (entry.theme.FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal))
	entry.updateBounds()
}

func (entry *ListEntry) SetConfig (new tomo.Config) {
	if new == entry.config.Config { return }
	entry.config.Config = new
}

func (entry *ListEntry) updateBounds () {
	padding := entry.theme.Padding(tomo.PatternRaised)
	entry.bounds = padding.Inverse().Apply(entry.drawer.LayoutBounds())
	entry.bounds = entry.bounds.Sub(entry.bounds.Min)
	entry.minimumWidth = entry.bounds.Dx()
	entry.bounds.Max.X = entry.width
}

func (entry *ListEntry) Draw (
	destination canvas.Canvas,
	offset image.Point,
	focused bool,
	on bool,
) (
	updatedRegion image.Rectangle,
) {
	state := tomo.State {
		Focused: focused,
		On: on,
	}

	pattern := entry.theme.Pattern(tomo.PatternRaised, state)
	padding := entry.theme.Padding(tomo.PatternRaised)
	bounds  := entry.Bounds().Add(offset)
	pattern.Draw(destination, bounds)
		
	foreground := entry.theme.Color (tomo.ColorForeground, state)
	return entry.drawer.Draw (
		destination,
		foreground,
		offset.Add(image.Pt(padding[artist.SideLeft], padding[artist.SideTop])).
		Sub(entry.drawer.LayoutBounds().Min))
}

func (entry *ListEntry) RunSelect () {
	if entry.onSelect != nil {
		entry.onSelect()
	}
}

func (entry *ListEntry) Bounds () (bounds image.Rectangle) {
	return entry.bounds
}

func (entry *ListEntry) Resize (width int) {
	entry.width = width
	entry.updateBounds()
}

func (entry *ListEntry) MinimumWidth () (width int) {
	return entry.minimumWidth
} 
