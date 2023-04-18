package elements

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"

type listEntity interface {
	tomo.ContainerEntity
	tomo.ScrollableEntity
}

type List struct {
	entity listEntity
	
	scratch       map[tomo.Element] scratchEntry
	scroll        image.Point
	contentBounds image.Rectangle
	columnSizes   []int
	selected      int
	
	forcedMinimumWidth  int
	forcedMinimumHeight int

	theme theme.Wrapped
	
	onScrollBoundsChange func ()
}

func NewList (columns int, children ...tomo.Selectable) (element *List) {
	if columns < 1 { columns = 1 }
	element = &List { selected: -1 }
	element.scratch = make(map[tomo.Element] scratchEntry)
	element.columnSizes = make([]int, columns)
	element.theme.Case = tomo.C("tomo", "list")
	element.entity = tomo.NewEntity(element).(listEntity)

	for _, child := range children {
		element.Adopt(child)
	}
	return
}

func (element *List) Entity () tomo.Entity {
	return element.entity
}

func (element *List) Draw (destination canvas.Canvas) {
	rocks := make([]image.Rectangle, element.entity.CountChildren())
	for index := 0; index < element.entity.CountChildren(); index ++ {
		rocks[index] = element.entity.Child(index).Entity().Bounds()
	}

	pattern := element.theme.Pattern(tomo.PatternSunken, tomo.State { })
	artist.DrawShatter(destination, pattern, element.entity.Bounds(), rocks...)
}

func (element *List) Layout () {
	if element.scroll.Y > element.maxScrollHeight() {
		element.scroll.Y = element.maxScrollHeight()
	}
	
	margin := element.theme.Margin(tomo.PatternSunken)
	padding := element.theme.Padding(tomo.PatternSunken)
	bounds := padding.Apply(element.entity.Bounds())
	element.contentBounds = image.Rectangle { }

	dot         := bounds.Min.Sub(element.scroll)
	xStart      := dot.X
	rowHeight   := 0
	columnIndex := 0
	nextLine := func () {
		dot.X = xStart
		dot.Y += margin.Y
		dot.Y += rowHeight
		rowHeight   = 0
		columnIndex = 0
	}

	for index := 0; index < element.entity.CountChildren(); index ++ {
		child := element.entity.Child(index)
		entry := element.scratch[child]
	
		if columnIndex >= len(element.columnSizes) {
			nextLine()
		}
		width  := element.columnSizes[columnIndex]
		height := int(entry.minSize)

		if len(element.columnSizes) == 1 && width < bounds.Dx() {
			width = bounds.Dx()
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

		columnIndex ++
	}
	
	element.contentBounds =
		element.contentBounds.Sub(element.contentBounds.Min)
		
	element.entity.NotifyScrollBoundsChange()
	if element.onScrollBoundsChange != nil {
		element.onScrollBoundsChange()
	}
}

func (element *List) Adopt (child tomo.Element) {
	element.entity.Adopt(child)
	element.scratch[child] = scratchEntry { }
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

func (element *List) Disown (child tomo.Element) {
	index := element.entity.IndexOf(child)
	if index < 0 { return }
	if index == element.selected {
		element.selected = -1
		element.entity.SelectChild(index, false)
	}
	element.entity.Disown(index)
	delete(element.scratch, child)
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

func (element *List) DisownAll () {
	func () {
		for index := 0; index < element.entity.CountChildren(); index ++ {
			index := index
			defer element.entity.Disown(index)
		}
	} ()
	element.scratch = make(map[tomo.Element] scratchEntry)
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

func (element *List) HandleChildMouseDown (x, y int, button input.Button, child tomo.Element) {
	if child, ok := child.(tomo.Selectable); ok {
		index := element.entity.IndexOf(child)
		if element.selected == index { return }
		if element.selected >= 0 {
			element.entity.SelectChild(element.selected, false)
		}
		element.selected = index
		element.entity.SelectChild(index, true)
	}
}

func (element *List) HandleChildMouseUp (int, int, input.Button, tomo.Element) { }

func (element *List) HandleChildMinimumSizeChange (child tomo.Element) {
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

func (element *List) HandleChildFlexibleHeightChange (child tomo.Flexible) {
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

func (element *List) DrawBackground (destination canvas.Canvas) {
	element.entity.DrawBackground(destination)
}

// SetTheme sets the element's theme.
func (element *List) SetTheme (theme tomo.Theme) {
	if theme == element.theme.Theme { return }
	element.theme.Theme = theme
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

// Collapse forces a minimum width and height upon the list. If a zero value is
// given for a dimension, its minimum will be determined by the list's content.
// If the list's height goes beyond the forced size, it will need to be accessed
// via scrolling. If an entry's width goes beyond the forced size, its text will
// be truncated so that it fits.
func (element *List) Collapse (width, height int) {
	if
		element.forcedMinimumWidth == width &&
		element.forcedMinimumHeight == height {
		
		return
	}
	
	element.forcedMinimumWidth  = width
	element.forcedMinimumHeight = height
	
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

// ScrollContentBounds returns the full content size of the element.
func (element *List) ScrollContentBounds () image.Rectangle {
	return element.contentBounds
}

// ScrollViewportBounds returns the size and position of the element's
// viewport relative to ScrollBounds.
func (element *List) ScrollViewportBounds () image.Rectangle {
	padding := element.theme.Padding(tomo.PatternBackground)
	bounds  := padding.Apply(element.entity.Bounds())
	bounds   = bounds.Sub(bounds.Min).Add(element.scroll)
	return bounds
}

// ScrollTo scrolls the viewport to the specified point relative to
// ScrollBounds.
func (element *List) ScrollTo (position image.Point) {
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
func (element *List) OnScrollBoundsChange (callback func ()) {
	element.onScrollBoundsChange = callback
}

// ScrollAxes returns the supported axes for scrolling.
func (element *List) ScrollAxes () (horizontal, vertical bool) {
	return false, true
}

func (element *List) maxScrollHeight () (height int) {
	padding := element.theme.Padding(tomo.PatternSunken)
	viewportHeight := element.entity.Bounds().Dy() - padding.Vertical()
	height = element.contentBounds.Dy() - viewportHeight
	if height < 0 { height = 0 }
	return
}

func (element *List) updateMinimumSize () {
	margin := element.theme.Margin(tomo.PatternSunken)
	padding := element.theme.Padding(tomo.PatternSunken)

	for index := range element.columnSizes {
		element.columnSizes[index] = 0
	}

	height      := 0
	rowHeight   := 0
	columnIndex := 0
	nextLine := func () {
		height += rowHeight
		rowHeight   = 0
		columnIndex = 0
	}
	for index := 0; index < element.entity.CountChildren(); index ++ {
		if columnIndex >= len(element.columnSizes) {
			if index > 0 { height += margin.Y }
			nextLine()
		}

		child := element.entity.Child(index)
		entry := element.scratch[child]
		
		entryWidth, entryHeight := element.entity.ChildMinimumSize(index)
		entry.minBreadth = float64(entryWidth)
		entry.minSize    = float64(entryHeight)
		element.scratch[child] = entry
		
		if rowHeight < entryHeight {
			rowHeight = entryHeight
		}
		if element.columnSizes[columnIndex] < entryWidth {
			element.columnSizes[columnIndex] = entryWidth
		}

		columnIndex ++
	}
	nextLine()

	width := 0; for index, size := range element.columnSizes {
		width += size
		if index > 0 { width += margin.X }
	}
	width  += padding.Horizontal()
	height += padding.Vertical()

	if element.forcedMinimumHeight > 0 {
		height = element.forcedMinimumHeight
	}
	if element.forcedMinimumWidth > 0 {
		width = element.forcedMinimumWidth
	}

	element.entity.SetMinimumSize(width, height)
}
