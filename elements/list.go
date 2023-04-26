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
	tomo.FocusableEntity
}

type list struct {
	container
	entity listEntity

	enabled       bool
	scroll        image.Point
	contentBounds image.Rectangle
	selected      int
	
	forcedMinimumWidth  int
	forcedMinimumHeight int

	theme theme.Wrapped

	onClick func ()
	onSelectionChange func ()
	onScrollBoundsChange func ()
}

type List struct {
	list
}

type FlowList struct {
	list
}

func NewList (children ...tomo.Element) (element *List) {
	element = &List { }
	element.theme.Case       = tomo.C("tomo", "list")
	element.entity           = tomo.NewEntity(element).(listEntity)
	element.container.entity = element.entity
	element.minimumSize      = element.updateMinimumSize
	element.init(children...)
	return
}

func NewFlowList (children ...tomo.Element) (element *FlowList) {
	element = &FlowList { }
	element.theme.Case       = tomo.C("tomo", "flowList")
	element.entity           = tomo.NewEntity(element).(listEntity)
	element.container.entity = element.entity
	element.minimumSize      = element.updateMinimumSize
	element.init(children...)
	return
}

func (element *list) init (children ...tomo.Element) {
	element.selected = -1
	element.enabled  = true
	element.container.init()
	element.Adopt(children...)
}

func (element *list) Draw (destination canvas.Canvas) {
	rocks := make([]image.Rectangle, element.entity.CountChildren())
	for index := 0; index < element.entity.CountChildren(); index ++ {
		rocks[index] = element.entity.Child(index).Entity().Bounds()
	}

	pattern := element.theme.Pattern(tomo.PatternSunken, element.state())
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

	dot := bounds.Min.Sub(element.scroll)

	for index := 0; index < element.entity.CountChildren(); index ++ {
		child := element.entity.Child(index)
		entry := element.scratch[child]
	
		width  := bounds.Dx()
		height := int(entry.minSize)

		childBounds := tomo.Bounds (
			dot.X, dot.Y,
			width, height)
		element.entity.PlaceChild(index, childBounds)
		element.contentBounds = element.contentBounds.Union(childBounds)

		dot.Y += height
		dot.Y += margin.Y
	}
	
	element.contentBounds =
		element.contentBounds.Sub(element.contentBounds.Min)
		
	element.entity.NotifyScrollBoundsChange()
	if element.onScrollBoundsChange != nil {
		element.onScrollBoundsChange()
	}
}

func (element *FlowList) Layout () {
	if element.scroll.Y > element.maxScrollHeight() {
		element.scroll.Y = element.maxScrollHeight()
	}
	
	margin := element.theme.Margin(tomo.PatternSunken)
	padding := element.theme.Padding(tomo.PatternSunken)
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
		
		dot.X += width + margin.X
	}
	
	element.contentBounds =
		element.contentBounds.Sub(element.contentBounds.Min)
		
	element.entity.NotifyScrollBoundsChange()
	if element.onScrollBoundsChange != nil {
		element.onScrollBoundsChange()
	}
}

func (element *list) Selected () tomo.Selectable {
	if element.selected == -1 { return nil }
	child, ok := element.entity.Child(element.selected).(tomo.Selectable)
	if !ok { return nil }
	return child
}

func (element *list) Select (child tomo.Selectable) {
	index := element.entity.IndexOf(child)
	if element.selected == index { return }
	element.selectNone()
	element.selected = index
	element.entity.SelectChild(index, true)
	if element.onSelectionChange != nil {
		element.onSelectionChange()
	}
	element.scrollToSelected()
}

func (element *list) Enabled () bool {
	return element.enabled
}

func (element *list) SetEnabled (enabled bool) {
	if element.enabled == enabled { return }
	element.enabled = enabled
	element.entity.Invalidate()
}

func (element *list) Focus () {
	element.entity.Focus()
}

func (element *list) HandleFocusChange () {
	element.entity.Invalidate()
}

func (element *list) HandleMouseDown (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) {
	if !element.enabled { return }
	element.Focus()
	element.selectNone()
}

func (element *list) HandleMouseUp (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) { }

func (element *list) HandleChildMouseDown (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
	child tomo.Element,
) {
	if !element.enabled { return }
	element.Focus()
	if child, ok := child.(tomo.Selectable); ok {
		element.Select(child)
	}
}

func (element *list) HandleChildMouseUp  (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
	child tomo.Element,
) {
	if !position.In(child.Entity().Bounds()) { return }
	if element.onClick != nil {
		element.onClick()
	}
}

func (element *list) HandleChildFlexibleHeightChange (child tomo.Flexible) {
	element.minimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

func (element *list) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
	if !element.Enabled() { return }
	index := -1
	switch key {
	case input.KeyUp, input.KeyLeft:
		index = element.selected - 1
	case input.KeyDown, input.KeyRight:
		index = element.selected + 1
	case input.KeyEnter:
		if element.onClick != nil {
			element.onClick()
		}
	}
	if index >= 0 && index < element.entity.CountChildren() {
		element.selectNone()
		element.selected = index
		element.entity.SelectChild(index, true)
	if element.onSelectionChange != nil {
		element.onSelectionChange()
	}
		element.scrollToSelected()
	}
}

func (element *list) HandleKeyUp(key input.Key, modifiers input.Modifiers) { }

func (element *list) DrawBackground (destination canvas.Canvas) {
	element.entity.DrawBackground(destination)
}

// SetTheme sets the element's theme.
func (element *list) SetTheme (theme tomo.Theme) {
	if theme == element.theme.Theme { return }
	element.theme.Theme = theme
	element.minimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

// Collapse forces a minimum width and height upon the list. If a zero value is
// given for a dimension, its minimum will be determined by the list's content.
// If the list's height goes beyond the forced size, it will need to be accessed
// via scrolling. If an entry's width goes beyond the forced size, its text will
// be truncated so that it fits.
func (element *list) Collapse (width, height int) {
	if
		element.forcedMinimumWidth == width &&
		element.forcedMinimumHeight == height {
		
		return
	}
	
	element.forcedMinimumWidth  = width
	element.forcedMinimumHeight = height
	
	element.minimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

// ScrollContentBounds returns the full content size of the element.
func (element *list) ScrollContentBounds () image.Rectangle {
	return element.contentBounds
}

// ScrollViewportBounds returns the size and position of the element's
// viewport relative to ScrollBounds.
func (element *list) ScrollViewportBounds () image.Rectangle {
	padding := element.theme.Padding(tomo.PatternSunken)
	bounds  := padding.Apply(element.entity.Bounds())
	bounds   = bounds.Sub(bounds.Min).Add(element.scroll)
	return bounds
}

// ScrollTo scrolls the viewport to the specified point relative to
// ScrollBounds.
func (element *list) ScrollTo (position image.Point) {
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
func (element *list) OnScrollBoundsChange (callback func ()) {
	element.onScrollBoundsChange = callback
}

func (element *list) OnClick (callback func ()) {
	element.onClick = callback
}

func (element *list) OnSelectionChange (callback func ()) {
	element.onSelectionChange = callback
}

// ScrollAxes returns the supported axes for scrolling.
func (element *list) ScrollAxes () (horizontal, vertical bool) {
	return false, true
}

func (element *list) selectNone () {
	if element.selected >= 0 {
		element.entity.SelectChild(element.selected, false)
	if element.onSelectionChange != nil {
		element.onSelectionChange()
	}
	}
}

func (element *list) scrollToSelected () {
	if element.selected < 0 { return }
	target := element.entity.Child(element.selected).Entity().Bounds()
	padding := element.theme.Padding(tomo.PatternSunken)
	bounds  := padding.Apply(element.entity.Bounds())
	if target.Min.Y < bounds.Min.Y {
		element.scroll.Y -= bounds.Min.Y - target.Min.Y
		element.entity.Invalidate()
		element.entity.InvalidateLayout()
	} else if target.Max.Y > bounds.Max.Y {
		element.scroll.Y += target.Max.Y - bounds.Max.Y
		element.entity.Invalidate()
		element.entity.InvalidateLayout()
	} 
}

func (element *list) state () tomo.State {
	return tomo.State {
		Focused: element.entity.Focused(),
		Disabled: !element.enabled,
	}
}

func (element *list) maxScrollHeight () (height int) {
	padding := element.theme.Padding(tomo.PatternSunken)
	viewportHeight := element.entity.Bounds().Dy() - padding.Vertical()
	height = element.contentBounds.Dy() - viewportHeight
	if height < 0 { height = 0 }
	return
}

func (element *List) updateMinimumSize () {
	margin := element.theme.Margin(tomo.PatternSunken)
	padding := element.theme.Padding(tomo.PatternSunken)

	width  := 0
	height := 0
	for index := 0; index < element.entity.CountChildren(); index ++ {
		if index > 0 { height += margin.Y }

		child := element.entity.Child(index)
		entry := element.scratch[child]
		
		entryWidth, entryHeight := element.entity.ChildMinimumSize(index)
		entry.minBreadth = float64(entryWidth)
		entry.minSize    = float64(entryHeight)
		element.scratch[child] = entry

		height += entryHeight
		if width < entryWidth { width = entryWidth }
	}

	width  += padding.Horizontal()
	height += padding.Vertical()

	if element.forcedMinimumWidth > 0 {
		width = element.forcedMinimumWidth
	}
	if element.forcedMinimumHeight > 0 {
		height = element.forcedMinimumHeight
	}

	element.entity.SetMinimumSize(width, height)
}

func (element *FlowList) updateMinimumSize () {
	padding := element.theme.Padding(tomo.PatternSunken)
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
