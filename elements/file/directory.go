package fileElements

import "io/fs"
import "image"
import "path/filepath"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/textdraw"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

type fileLayoutEntry struct {
	*File
	fs.DirEntry
	Bounds image.Rectangle
	Drawer textdraw.Drawer
	TextPoint image.Point
}

type historyEntry struct {
	location string
	filesystem ReadDirStatFS
}

// Directory displays a list of files within a particular directory and
// file system.
type Directory struct {
	*core.Core
	*core.Propagator
	core core.CoreControl

	children []fileLayoutEntry
	scroll   image.Point
	contentBounds image.Rectangle
	
	config config.Wrapped
	theme  theme.Wrapped

	onScrollBoundsChange func ()

	history      []historyEntry
	historyIndex int
	onChoose func (file string)
}

// NewDirectory creates a new directory view. If within is nil, it will use
// the OS file system.
func NewDirectory (
	location string,
	within ReadDirStatFS,
) (
	element *Directory,
	err error,
) {
	element = &Directory { }
	element.theme.Case = theme.C("files", "directory")
	element.Core, element.core = core.NewCore(element, element.redoAll)
	element.Propagator = core.NewPropagator(element, element.core)
	err = element.SetLocation(location, within)
	return
}

// Location returns the directory's location and filesystem.
func (element *Directory) Location () (string, ReadDirStatFS) {
	if len(element.history) < 1 { return "", nil }
	current := element.history[element.historyIndex]
	return current.location, current.filesystem
}

// SetLocation sets the directory's location and filesystem. If within is nil,
// it will use the OS file system.
func (element *Directory) SetLocation (
	location string,
	within ReadDirStatFS,
) error {
	if within == nil {
		within = defaultFS { }
	}
	element.scroll = image.Point { }

	if element.history != nil {
		element.historyIndex ++
	}
	element.history = append (
		element.history[:element.historyIndex],
		historyEntry { location, within })
	return element.Update()
}

// Backward goes back a directory in history
func (element *Directory) Backward () (bool, error) {
	if element.historyIndex > 1 {
		element.historyIndex --
		return true, element.Update()
	} else {
		return false, nil
	}
}

// Forward goes forward a directory in history
func (element *Directory) Forward () (bool, error) {
	if element.historyIndex < len(element.history) - 1 {
		element.historyIndex ++
		return true, element.Update()
	} else {
		return false, nil
	}
}

// Update refreshes the directory's contents.
func (element *Directory) Update () error {
	location, filesystem := element.Location()
	entries, err := filesystem.ReadDir(location)

	// disown all entries
	for _, file := range element.children {
		file.DrawTo(nil, image.Rectangle { }, nil)
		file.SetParent(nil)
		
		if file.Focused() {
			file.HandleUnfocus()
		}
	}

	element.children = make([]fileLayoutEntry, len(entries))
	for index, entry := range entries {
		filePath := filepath.Join(location, entry.Name())
		file, err := NewFile(filePath, filesystem)
		if err != nil { continue }
		file.SetParent(element)
		file.OnChoose (func () {
			if element.onChoose != nil {
				element.onChoose(filePath)
			}
		})
		element.children[index].File = file
		element.children[index].DirEntry = entry
		element.children[index].Drawer.SetFace (element.theme.FontFace(
			theme.FontStyleRegular,
			theme.FontSizeNormal))
		element.children[index].Drawer.SetText([]rune(entry.Name()))
		element.children[index].Drawer.SetAlign(textdraw.AlignCenter)
	}
	
	if element.core.HasImage() {
		element.redoAll()
		element.core.DamageAll()
	}
	return err
}

// OnChoose sets a function to be called when the user double-clicks a file or
// sub-directory within the directory view.
func (element *Directory) OnChoose (callback func (file string)) {
	element.onChoose = callback
}

// CountChildren returns the amount of children contained within this element.
func (element *Directory) CountChildren () (count int) {
	return len(element.children)
}

// Child returns the child at the specified index. If the index is out of
// bounds, this method will return nil.
func (element *Directory) Child (index int) (child elements.Element) {
	if index < 0 || index > len(element.children) { return }
	return element.children[index].File
}

func (element *Directory) HandleMouseDown (x, y int, button input.Button) {
	if button == input.ButtonLeft {
		var file *File
		for _, entry := range element.children {
			if image.Pt(x, y).In(entry.Bounds) {
				file = entry.File
			}
		}
		if file != nil {
			file.SetSelected(!file.Selected())
		}
	}
	element.Propagator.HandleMouseDown(x, y, button)
}

func (element *Directory) redoAll () {
	if !element.core.HasImage() { return }
	
	// do a layout
	element.doLayout()
	
	maxScrollHeight := element.maxScrollHeight()
	if element.scroll.Y > maxScrollHeight {
		element.scroll.Y = maxScrollHeight
		element.doLayout()
	}

	// draw a background
	rocks := make([]image.Rectangle, len(element.children))
	for index, entry := range element.children {
		rocks[index] = entry.Bounds
	}
	pattern := element.theme.Pattern (
		theme.PatternPinboard,
		theme.State { })
	artist.DrawShatter(element.core, pattern, element.Bounds(), rocks...)

	element.partition()
	if parent, ok := element.core.Parent().(elements.ScrollableParent); ok {
		parent.NotifyScrollBoundsChange(element)
	}
	if element.onScrollBoundsChange != nil {
		element.onScrollBoundsChange()
	}

	// draw labels
	foreground := element.theme.Color(theme.ColorForeground, theme.State { })
	for _, entry := range element.children {
		entry.Drawer.Draw(element.core, foreground, entry.TextPoint)
	}
}

func (element *Directory) partition () {
	for _, entry := range element.children {
		entry.DrawTo(nil, entry.Bounds, nil)
	}

	// cut our canvas up and give peices to child elements
	for _, entry := range element.children {
		if entry.Bounds.Overlaps(element.Bounds()) {
			entry.DrawTo (	
				canvas.Cut(element.core, entry.Bounds),
				entry.Bounds, func (region image.Rectangle) {
					element.core.DamageRegion(region)
				})
		}
	}
}

// NotifyMinimumSizeChange notifies the container that the minimum size of a
// child element has changed.
func (element *Directory) NotifyMinimumSizeChange (child elements.Element) {
	element.redoAll()
	element.core.DamageAll()
}

// SetTheme sets the element's theme.
func (element *Directory) SetTheme (new theme.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.Propagator.SetTheme(new)
	element.redoAll()
}

// SetConfig sets the element's configuration.
func (element *Directory) SetConfig (new config.Config) {
	if new == element.config.Config { return }
	element.Propagator.SetConfig(new)
	element.redoAll()
}
// ScrollContentBounds returns the full content size of the element.
func (element *Directory) ScrollContentBounds () image.Rectangle {
	return element.contentBounds
}

// ScrollViewportBounds returns the size and position of the element's
// viewport relative to ScrollBounds.
func (element *Directory) ScrollViewportBounds () image.Rectangle {
	padding := element.theme.Padding(theme.PatternPinboard)
	bounds  := padding.Apply(element.Bounds())
	bounds   = bounds.Sub(bounds.Min).Add(element.scroll)
	return bounds
}

// ScrollTo scrolls the viewport to the specified point relative to
// ScrollBounds.
func (element *Directory) ScrollTo (position image.Point) {
	if position.Y < 0 {
		position.Y = 0
	}
	maxScrollHeight := element.maxScrollHeight()
	if position.Y > maxScrollHeight {
		position.Y = maxScrollHeight
	}
	element.scroll = position
	if element.core.HasImage() {
		element.redoAll()
		element.core.DamageAll()
	}
}

// OnScrollBoundsChange sets a function to be called when the element's viewport
// bounds, content bounds, or scroll axes change.
func (element *Directory) OnScrollBoundsChange (callback func ()) {
	element.onScrollBoundsChange = callback
}

// ScrollAxes returns the supported axes for scrolling.
func (element *Directory) ScrollAxes () (horizontal, vertical bool) {
	return false, true
}

func (element *Directory) maxScrollHeight () (height int) {
	padding := element.theme.Padding(theme.PatternSunken)
	viewportHeight := element.Bounds().Dy() - padding.Vertical()
	height = element.contentBounds.Dy() - viewportHeight
	if height < 0 { height = 0 }
	return
}

func (element *Directory) doLayout () {
	margin := element.theme.Margin(theme.PatternPinboard)
	padding := element.theme.Padding(theme.PatternPinboard)
	bounds := padding.Apply(element.Bounds())
	element.contentBounds = image.Rectangle { }

	beginningOfRow := true
	dot := bounds.Min.Sub(element.scroll)
	rowHeight := 0
	for index, entry := range element.children {
		width, height := entry.MinimumSize()
		
		if dot.X + width > bounds.Max.X {
			dot.X = bounds.Min.Sub(element.scroll).X
			dot.Y += rowHeight
			if index > 1 {
				dot.Y += margin.Y
			}
			beginningOfRow = true
		}
		
		if beginningOfRow {
			beginningOfRow = false
		} else {
			dot.X += margin.X
		}
	
		entry.Drawer.SetMaxWidth(width)
		bounds := image.Rect(dot.X, dot.Y, dot.X + width, dot.Y + height)
		entry.Bounds = bounds
		
		drawerHeight := entry.Drawer.ReccomendedHeightFor(width)
		entry.TextPoint =
			image.Pt(bounds.Min.X, bounds.Max.Y + margin.Y).
			Sub(entry.Drawer.LayoutBounds().Min)
		bounds.Max.Y += margin.Y + drawerHeight
		height += margin.Y + drawerHeight
		if rowHeight < height {
			rowHeight = height
		}
		
		element.contentBounds = element.contentBounds.Union(bounds)
		element.children[index] = entry
		dot.X += width
	}
	
	element.contentBounds =
		element.contentBounds.Sub(element.contentBounds.Min)
}

func (element *Directory) updateMinimumSize () {
	padding := element.theme.Padding(theme.PatternPinboard)
	minimumWidth := 0
	for _, entry := range element.children {
		width, _ := entry.MinimumSize()
		if width > minimumWidth {
			minimumWidth = width
		}
	}
	element.core.SetMinimumSize (
		minimumWidth + padding.Horizontal(),
		padding.Vertical())
}