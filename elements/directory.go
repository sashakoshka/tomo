package elements

import "image"
import "path/filepath"
import "tomo"
import "tomo/input"
import "art"
import "tomo/ability"
import "art/shatter"

// TODO: base on flow implementation of list. also be able to switch to a table
// variant for a more information dense view.

var directoryCase = tomo.C("tomo", "list")

type historyEntry struct {
	location string
	filesystem ReadDirStatFS
}

// Directory displays a list of files within a particular directory and
// file system.
type Directory struct {
	container
	entity tomo.Entity
	
	scroll        image.Point
	contentBounds image.Rectangle
	
	history      []historyEntry
	historyIndex int
	
	onChoose             func (file string)
	onScrollBoundsChange func ()
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
	element.entity = tomo.GetBackend().NewEntity(element)
	element.container.entity = element.entity
	element.minimumSize = element.updateMinimumSize
	element.init()
	err = element.SetLocation(location, within)
	return
}

func (element *Directory) Draw (destination art.Canvas) {
	rocks := make([]image.Rectangle, element.entity.CountChildren())
	for index := 0; index < element.entity.CountChildren(); index ++ {
		rocks[index] = element.entity.Child(index).Entity().Bounds()
	}

	tiles := shatter.Shatter(element.entity.Bounds(), rocks...)
	for _, tile := range tiles {
		element.DrawBackground(art.Cut(destination, tile))
	}
}

func (element *Directory) Layout () {
	if element.scroll.Y > element.maxScrollHeight() {
		element.scroll.Y = element.maxScrollHeight()
	}
	
	margin := element.entity.Theme().Margin(tomo.PatternPinboard, directoryCase)
	padding := element.entity.Theme().Padding(tomo.PatternPinboard, directoryCase)
	bounds := padding.Apply(element.entity.Bounds())
	element.contentBounds = image.Rectangle { }

	dot := bounds.Min.Sub(element.scroll)
	xStart := dot.X
	rowHeight := 0

	nextLine := func () {
		dot.X = xStart
		dot.Y += margin.Y
		dot.Y += rowHeight
		rowHeight = 0
	}
	
	for index := 0; index < element.entity.CountChildren(); index ++ {
		child := element.entity.Child(index)
		entry := element.scratch[child]
	
		width  := int(entry.minBreadth)
		height := int(entry.minSize)
		if width + dot.X > bounds.Max.X {
			nextLine()
		}
		if typedChild, ok := child.(ability.Flexible); ok {
			height = typedChild.FlexibleHeightFor(width)
		}
		if rowHeight < height {
			rowHeight = height
		}

		childBounds := tomo.Bounds (
			dot.X, dot.Y,
			width, height)
		element.entity.PlaceChild(index, childBounds)
		element.contentBounds = element.contentBounds.Union(childBounds)
		
		dot.X += width + margin.X
	}
	
	element.contentBounds =
		element.contentBounds.Sub(element.contentBounds.Min)
		
	element.entity.NotifyScrollBoundsChange()
	if element.onScrollBoundsChange != nil {
		element.onScrollBoundsChange()
	}
}

func (element *Directory) HandleMouseDown  (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) {
	element.selectNone()
}

func (element *Directory) HandleMouseUp  (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) { }

func (element *Directory) HandleChildMouseDown  (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
	child tomo.Element,
) {
	element.selectNone()
	if child, ok := child.(ability.Selectable); ok {
		index := element.entity.IndexOf(child)
		element.entity.SelectChild(index, true)
	}
}

func (element *Directory) HandleChildMouseUp  (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
	child tomo.Element,
) { }

func (element *Directory) HandleChildFlexibleHeightChange (child ability.Flexible) {
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

// ScrollContentBounds returns the full content size of the element.
func (element *Directory) ScrollContentBounds () image.Rectangle {
	return element.contentBounds
}

// ScrollViewportBounds returns the size and position of the element's
// viewport relative to ScrollBounds.
func (element *Directory) ScrollViewportBounds () image.Rectangle {
	padding := element.entity.Theme().Padding(tomo.PatternPinboard, directoryCase)
	bounds  := padding.Apply(element.entity.Bounds())
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
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
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

func (element *Directory) DrawBackground (destination art.Canvas) {
	element.entity.Theme().Pattern(tomo.PatternPinboard, tomo.State { }, directoryCase).
		Draw(destination, element.entity.Bounds())
}

func (element *Directory) HandleThemeChange () {
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
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

	children := make([]tomo.Element, len(entries))
	for index, entry := range entries {
		filePath := filepath.Join(location, entry.Name())
		file, _ := NewFile(filePath, filesystem)
		file.OnChoose (func () {
			if element.onChoose != nil {
				element.onChoose(filePath)
			}
		})
		
		children[index] = file
	}

	element.DisownAll()
	element.Adopt(children...)
	return err
}

// OnChoose sets a function to be called when the user double-clicks a file or
// sub-directory within the directory view.
func (element *Directory) OnChoose (callback func (file string)) {
	element.onChoose = callback
}

func (element *Directory) selectNone () {
	for index := 0; index < element.entity.CountChildren(); index ++ {
		element.entity.SelectChild(index, false)
	}
}

func (element *Directory) maxScrollHeight () (height int) {
	padding := element.entity.Theme().Padding(tomo.PatternSunken, directoryCase)
	viewportHeight := element.entity.Bounds().Dy() - padding.Vertical()
	height = element.contentBounds.Dy() - viewportHeight
	if height < 0 { height = 0 }
	return
}


func (element *Directory) updateMinimumSize () {
	padding := element.entity.Theme().Padding(tomo.PatternPinboard, directoryCase)
	minimumWidth := 0
	for index := 0; index < element.entity.CountChildren(); index ++ {
		width, height := element.entity.ChildMinimumSize(index)
		if width > minimumWidth {
			minimumWidth = width
		}
		
		key   := element.entity.Child(index)
		entry := element.scratch[key]
		entry.minSize    = float64(height)
		entry.minBreadth = float64(width)
		element.scratch[key] = entry
	}
	element.entity.SetMinimumSize (
		minimumWidth + padding.Horizontal(),
		padding.Vertical())
}
