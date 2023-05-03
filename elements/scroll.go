package elements

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/ability"

var scrollCase = tomo.C("tomo", "scroll")

// ScrollMode specifies which sides of a Scroll have scroll bars.
type ScrollMode int; const (
	ScrollNeither    ScrollMode = 0
	ScrollVertical   ScrollMode = 1
	ScrollHorizontal ScrollMode = 2
	ScrollBoth       ScrollMode = ScrollVertical | ScrollHorizontal
)

// Includes returns whether a scroll mode has been or'd with another scroll
// mode.
func (mode ScrollMode) Includes (sub ScrollMode) bool {
	return (mode & sub) > 0
}

// Scroll adds scroll bars to any scrollable element. It also captures scroll
// wheel input.
type Scroll struct {
	entity tomo.Entity
	
	child      ability.Scrollable
	horizontal *ScrollBar
	vertical   *ScrollBar
}

// NewScroll creates a new scroll element.
func NewScroll (mode ScrollMode, child ability.Scrollable) (element *Scroll) {
	element = &Scroll { }
	element.entity = tomo.GetBackend().NewEntity(element)

	if mode.Includes(ScrollHorizontal) {
		element.horizontal = NewHScrollBar()
		element.horizontal.OnScroll (func (viewport image.Point) {
			if element.child != nil {
				element.child.ScrollTo(viewport)
			}
			if element.vertical != nil {
				element.vertical.SetBounds (
					element.child.ScrollContentBounds(),
					element.child.ScrollViewportBounds())
			}
		})
		element.entity.Adopt(element.horizontal)
	}
	if mode.Includes(ScrollVertical) {
		element.vertical = NewVScrollBar()
		element.vertical.OnScroll (func (viewport image.Point) {
			if element.child != nil {
				element.child.ScrollTo(viewport)
			}
			if element.horizontal != nil {
				element.horizontal.SetBounds (
					element.child.ScrollContentBounds(),
					element.child.ScrollViewportBounds())
			}
		})
		element.entity.Adopt(element.vertical)
	}

	element.Adopt(child)
	return
}

// Entity returns this element's entity.
func (element *Scroll) Entity () tomo.Entity {
	return element.entity
}

// Draw causes the element to draw to the specified destination canvas.
func (element *Scroll) Draw (destination artist.Canvas) {
	if element.horizontal != nil && element.vertical != nil {
		bounds := element.entity.Bounds()
		bounds.Min = image.Pt (
			bounds.Max.X - element.vertical.Entity().Bounds().Dx(),
			bounds.Max.Y - element.horizontal.Entity().Bounds().Dy())
		state := tomo.State { }
		deadArea := element.entity.Theme().Pattern(tomo.PatternDead, state, scrollCase)
		deadArea.Draw(artist.Cut(destination, bounds), bounds)
	}
}

// Layout causes this element to perform a layout operation.
func (element *Scroll) Layout () {
	bounds := element.entity.Bounds()
	child  := bounds

	iHorizontal := element.entity.IndexOf(element.horizontal)
	iVertical   := element.entity.IndexOf(element.vertical)
	iChild      := element.entity.IndexOf(element.child)

	var horizontal, vertical image.Rectangle

	if element.horizontal != nil {
		_, hMinHeight := element.entity.ChildMinimumSize(iHorizontal)
		child.Max.Y -= hMinHeight
	}
	if element.vertical != nil {
		vMinWidth, _  := element.entity.ChildMinimumSize(iVertical)
		child.Max.X -= vMinWidth
	}

	horizontal.Min.X = bounds.Min.X
	horizontal.Max.X = child.Max.X
	horizontal.Min.Y = child.Max.Y
	horizontal.Max.Y = bounds.Max.Y
	
	vertical.Min.X = child.Max.X
	vertical.Max.X = bounds.Max.X
	vertical.Min.Y = bounds.Min.Y
	vertical.Max.Y = child.Max.Y

	if element.horizontal != nil {
		element.entity.PlaceChild (iHorizontal, horizontal)
	}
	if element.vertical != nil {
		element.entity.PlaceChild(iVertical, vertical)
	}
	if element.child != nil {
		element.entity.PlaceChild(iChild, child)
	}
}

// DrawBackground draws this element's background pattern to the specified
// destination canvas.
func (element *Scroll) DrawBackground (destination artist.Canvas) {
	element.entity.DrawBackground(destination)
}

// Adopt sets this element's child. If nil is passed, any child is removed.
func (element *Scroll) Adopt (child ability.Scrollable) {
	if element.child != nil {
		element.entity.Disown(element.entity.IndexOf(element.child))
	}
	if child != nil {
		element.entity.Adopt(child)
	}
	element.child = child

	element.updateEnabled()
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

// Child returns this element's child. If there is no child, this method will
// return nil.
func (element *Scroll) Child () ability.Scrollable {
	return element.child
}

func (element *Scroll) HandleChildMinimumSizeChange (tomo.Element) {
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

func (element *Scroll) HandleChildScrollBoundsChange (ability.Scrollable) {
	element.updateEnabled()
	viewportBounds := element.child.ScrollViewportBounds()
	contentBounds  := element.child.ScrollContentBounds()
	if element.horizontal != nil {
		element.horizontal.SetBounds(contentBounds, viewportBounds)
	}
	if element.vertical != nil {
		element.vertical.SetBounds(contentBounds, viewportBounds)
	}
}

func (element *Scroll) HandleScroll (
	position image.Point,
	deltaX, deltaY float64,
	modifiers input.Modifiers,
) {
	horizontal, vertical := element.child.ScrollAxes()
	if !horizontal { deltaX = 0 }
	if !vertical   { deltaY = 0 }
	element.scrollChildBy(int(deltaX), int(deltaY))
}

func (element *Scroll) HandleThemeChange () {
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

func (element *Scroll) updateMinimumSize () {
	var width, height int

	if element.child != nil {
		width, height = element.entity.ChildMinimumSize (
			element.entity.IndexOf(element.child))
	}
	if element.horizontal != nil {
		hMinWidth, hMinHeight := element.entity.ChildMinimumSize (
			element.entity.IndexOf(element.horizontal))
		height += hMinHeight
		if hMinWidth > width {
			width = hMinWidth
		}
	}
	if element.vertical != nil {
		vMinWidth, vMinHeight := element.entity.ChildMinimumSize (
			element.entity.IndexOf(element.vertical))
		width += vMinWidth
		if vMinHeight > height {
			height = vMinHeight
		}
	}
	element.entity.SetMinimumSize(width, height)
}

func (element *Scroll) updateEnabled () {
	horizontal, vertical := false, false
	if element.child != nil {
		horizontal, vertical = element.child.ScrollAxes()
	}
	if element.horizontal != nil {
		element.horizontal.SetEnabled(horizontal)
	}
	if element.vertical != nil {
		element.vertical.SetEnabled(vertical)
	}
}

func (element *Scroll) scrollChildBy (x, y int) {
	if element.child == nil { return }
	scrollPoint :=
		element.child.ScrollViewportBounds().Min.
		Add(image.Pt(x, y))
	element.child.ScrollTo(scrollPoint)
}
