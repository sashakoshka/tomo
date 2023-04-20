package elements

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
// import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"
import "git.tebibyte.media/sashakoshka/tomo/default/config"

type ScrollMode int; const (
	ScrollNeither ScrollMode = 0
	ScrollVertical           = 1
	ScrollHorizontal         = 2
	ScrollBoth               = ScrollVertical | ScrollHorizontal
)

// Includes returns whether a scroll mode has been or'd with another scroll
// mode.
func (mode ScrollMode) Includes (sub ScrollMode) bool {
	return (mode & sub) > 0
}

type Scroll struct {
	entity tomo.ContainerEntity
	
	child      tomo.Scrollable
	horizontal *ScrollBar
	vertical   *ScrollBar
	
	config config.Wrapped
	theme  theme.Wrapped
}

func NewScroll (mode ScrollMode, child tomo.Scrollable) (element *Scroll) {
	element = &Scroll { }
	element.theme.Case = tomo.C("tomo", "scroll")
	element.entity = tomo.NewEntity(element).(tomo.ContainerEntity)

	if mode.Includes(ScrollHorizontal) {
		element.horizontal = NewScrollBar(false)
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
		element.vertical = NewScrollBar(true)
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

func (element *Scroll) Entity () tomo.Entity {
	return element.entity
}

func (element *Scroll) Draw (destination canvas.Canvas) {
	if element.horizontal != nil && element.vertical != nil {
		bounds := element.entity.Bounds()
		bounds.Min = image.Pt (
			bounds.Max.X - element.vertical.Entity().Bounds().Dx(),
			bounds.Max.Y - element.horizontal.Entity().Bounds().Dy())
		state := tomo.State { }
		deadArea := element.theme.Pattern(tomo.PatternDead, state)
		deadArea.Draw(canvas.Cut(destination, bounds), bounds)
	}
}

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

func (element *Scroll) DrawBackground (destination canvas.Canvas) {
	element.entity.DrawBackground(destination)
}

func (element *Scroll) Adopt (child tomo.Scrollable) {
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

func (element *Scroll) HandleChildMinimumSizeChange (tomo.Element) {
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

func (element *Scroll) HandleChildScrollBoundsChange (tomo.Scrollable) {
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
	x, y int,
	deltaX, deltaY float64,
) {
	horizontal, vertical := element.child.ScrollAxes()
	if !horizontal { deltaX = 0 }
	if !vertical   { deltaY = 0 }
	element.scrollChildBy(int(deltaX), int(deltaY))
}

func (element *Scroll) SetTheme (theme tomo.Theme) {
	if theme == element.theme.Theme { return }
	element.theme.Theme = theme
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

func (element *Scroll) SetConfig (config tomo.Config) {
	element.config.Config = config
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