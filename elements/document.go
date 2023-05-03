package elements

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/ability"
import "git.tebibyte.media/sashakoshka/tomo/shatter"

// Document is a scrollable container capcable of laying out flexible child
// elements. Children can be added either inline (similar to an HTML/CSS inline
// element), or expanding (similar to an HTML/CSS block element).
type Document struct {
	container
	entity tomo.Entity
	
	scroll        image.Point
	contentBounds image.Rectangle
	
	onScrollBoundsChange func ()
}

// NewDocument creates a new document container.
func NewDocument (children ...tomo.Element) (element *Document) {
	element = &Document { }
	element.theme.Case = tomo.C("tomo", "document")
	element.entity = tomo.NewEntity(element)
	element.container.entity = element.entity
	element.minimumSize = element.updateMinimumSize
	element.init()
	element.Adopt(children...)
	return
}

// Draw causes the element to draw to the specified destination canvas.
func (element *Document) Draw (destination artist.Canvas) {
	rocks := make([]image.Rectangle, element.entity.CountChildren())
	for index := 0; index < element.entity.CountChildren(); index ++ {
		rocks[index] = element.entity.Child(index).Entity().Bounds()
	}

	tiles := shatter.Shatter(element.entity.Bounds(), rocks...)
	for _, tile := range tiles {
		element.entity.DrawBackground(canvas.Cut(destination, tile))
	}
}

// Layout causes this element to perform a layout operation.
func (element *Document) Layout () {
	if element.scroll.Y > element.maxScrollHeight() {
		element.scroll.Y = element.maxScrollHeight()
	}
	
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
		if width + dot.X > bounds.Max.X && !entry.expand {
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

// Adopt adds one or more elements to the container, placing each on its own
// line.
func (element *Document) Adopt (children ...tomo.Element) {
	element.adopt(true, children...)
}

// AdoptInline adds one or more elements to the container, packing multiple
// elements onto the same line(s).
func (element *Document) AdoptInline (children ...tomo.Element) {
	element.adopt(false, children...)
}

func (element *Document) HandleChildFlexibleHeightChange (child ability.Flexible) {
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

// DrawBackground draws this element's background pattern to the specified
// destination canvas.
func (element *Document) DrawBackground (destination artist.Canvas) {
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
