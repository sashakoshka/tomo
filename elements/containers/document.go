package containers

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/shatter"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"

type documentEntity interface {
	tomo.ContainerEntity
	tomo.ScrollableEntity
}

type Document struct {
	entity documentEntity
	
	scratch       map[tomo.Element] scratchEntry
	scroll        image.Point
	contentBounds image.Rectangle

	theme theme.Wrapped
	
	onScrollBoundsChange func ()
}

func NewDocument () (element *Document) {
	element = &Document { }
	element.scratch = make(map[tomo.Element] scratchEntry)
	element.theme.Case = tomo.C("tomo", "document")
	element.entity = tomo.NewEntity(element).(documentEntity)
	return
}

func (element *Document) Entity () tomo.Entity {
	return element.entity
}

func (element *Document) Draw (destination canvas.Canvas) {
	rocks := make([]image.Rectangle, element.entity.CountChildren())
	for index := 0; index < element.entity.CountChildren(); index ++ {
		rocks[index] = element.entity.Child(index).Entity().Bounds()
	}

	tiles := shatter.Shatter(element.entity.Bounds(), rocks...)
	for _, tile := range tiles {
		element.entity.DrawBackground(canvas.Cut(destination, tile))
	}
}

func (element *Document) Layout () {
	margin := element.theme.Margin(tomo.PatternBackground)
	padding := element.theme.Padding(tomo.PatternBackground)
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
		
		if dot.X > xStart && entry.expand {
			nextLine()
		}
	
		width  := int(entry.minBreadth)
		height := int(entry.minSize)
		if width + dot.X > bounds.Dx() && !entry.expand {
			nextLine()
		}
		if width < bounds.Dx() && entry.expand {
			width = bounds.Dx()
		}
		if typedChild, ok := child.(tomo.Flexible); ok {
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
		
		if entry.expand {
			nextLine()
		} else {
			dot.X += width + margin.X
		}
	}
	
	element.contentBounds =
		element.contentBounds.Sub(element.contentBounds.Min)
		
	element.entity.NotifyScrollBoundsChange()
	if element.onScrollBoundsChange != nil {
		element.onScrollBoundsChange()
	}
}

func (element *Document) Adopt (child tomo.Element, expand bool) {
	element.entity.Adopt(child)
	element.scratch[child] = scratchEntry { expand: expand }
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

func (element *Document) Disown (child tomo.Element) {
	index := element.entity.IndexOf(child)
	if index < 0 { return }
	element.entity.Disown(index)
	delete(element.scratch, child)
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

func (element *Document) DisownAll () {
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

func (element *Document) HandleChildMinimumSizeChange (child tomo.Element) {
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

func (element *Document) HandleChildFlexibleHeightChange (child tomo.Flexible) {
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

func (element *Document) DrawBackground (destination canvas.Canvas) {
	element.entity.DrawBackground(destination)
}

// SetTheme sets the element's theme.
func (element *Document) SetTheme (theme tomo.Theme) {
	if theme == element.theme.Theme { return }
	element.theme.Theme = theme
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}


// ScrollContentBounds returns the full content size of the element.
func (element *Document) ScrollContentBounds () image.Rectangle {
	return element.contentBounds
}

// ScrollViewportBounds returns the size and position of the element's
// viewport relative to ScrollBounds.
func (element *Document) ScrollViewportBounds () image.Rectangle {
	padding := element.theme.Padding(tomo.PatternBackground)
	bounds  := padding.Apply(element.entity.Bounds())
	bounds   = bounds.Sub(bounds.Min).Add(element.scroll)
	return bounds
}

// ScrollTo scrolls the viewport to the specified point relative to
// ScrollBounds.
func (element *Document) ScrollTo (position image.Point) {
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
func (element *Document) OnScrollBoundsChange (callback func ()) {
	element.onScrollBoundsChange = callback
}

// ScrollAxes returns the supported axes for scrolling.
func (element *Document) ScrollAxes () (horizontal, vertical bool) {
	return false, true
}

func (element *Document) maxScrollHeight () (height int) {
	padding := element.theme.Padding(tomo.PatternSunken)
	viewportHeight := element.entity.Bounds().Dy() - padding.Vertical()
	height = element.contentBounds.Dy() - viewportHeight
	if height < 0 { height = 0 }
	return
}

func (element *Document) updateMinimumSize () {
	padding := element.theme.Padding(tomo.PatternBackground)
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
