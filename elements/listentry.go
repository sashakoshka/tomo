package elements

import "image"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/textdraw"

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
	entry.theme.Case = theme.C("tomo", "listEntry")
	entry.drawer.SetText([]rune(text))
	entry.updateBounds()
	return
}

func (entry *ListEntry) SetTheme (new theme.Theme) {
	if new == entry.theme.Theme { return }
	entry.theme.Theme = new
	entry.drawer.SetFace (entry.theme.FontFace (
		theme.FontStyleRegular,
		theme.FontSizeNormal))
	entry.updateBounds()
}

func (entry *ListEntry) SetConfig (new config.Config) {
	if new == entry.config.Config { return }
	entry.config.Config = new
}

func (entry *ListEntry) updateBounds () {
	padding := entry.theme.Padding(theme.PatternRaised)
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
	state := theme.State {
		Focused: focused,
		On: on,
	}

	pattern := entry.theme.Pattern(theme.PatternRaised, state)
	padding := entry.theme.Padding(theme.PatternRaised)
	bounds  := entry.Bounds().Add(offset)
	pattern.Draw(destination, bounds)
		
	foreground := entry.theme.Color (theme.ColorForeground, state)
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
