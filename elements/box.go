package elements

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/shatter"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"

// Space is a list of spacing configurations that can be passed to some
// containers.
type Space int; const (
	SpaceNone    = 0
	SpacePadding = 1
	SpaceMargin  = 2
	SpaceBoth    = SpacePadding | SpaceMargin
)

// Includes returns whether a spacing value has been or'd with another spacing
// value.
func (space Space) Includes (sub Space) bool {
	return (space & sub) > 0
}

// Box is a container that lays out its children horizontally or vertically.
// Child elements can be set to contract to their minimum size, or expand to
// fill remaining space. Boxes can be nested and used together to create more
// complex layouts.
type Box struct {
	container
	theme    theme.Wrapped
	padding  bool
	margin   bool
	vertical bool
}

// NewHBox creates a new horizontal box.
func NewHBox (space Space, children ...tomo.Element) (element *Box) {
	element = &Box {
		padding: space.Includes(SpacePadding),
		margin:  space.Includes(SpaceMargin),
	}
	element.entity = tomo.NewEntity(element).(tomo.ContainerEntity)
	element.minimumSize = element.updateMinimumSize
	element.init()
	element.theme.Case = tomo.C("tomo", "box")
	element.Adopt(children...)
	return
}

// NewHBox creates a new vertical box.
func NewVBox (space Space, children ...tomo.Element) (element *Box) {
	element = &Box {
		padding:  space.Includes(SpacePadding),
		margin:   space.Includes(SpaceMargin),
		vertical: true,
	}
	element.entity = tomo.NewEntity(element).(tomo.ContainerEntity)
	element.minimumSize = element.updateMinimumSize
	element.init()
	element.theme.Case = tomo.C("tomo", "box")
	element.Adopt(children...)
	return
}

func (element *Box) Draw (destination canvas.Canvas) {
	rocks := make([]image.Rectangle, element.entity.CountChildren())
	for index := 0; index < element.entity.CountChildren(); index ++ {
		rocks[index] = element.entity.Child(index).Entity().Bounds()
	}

	tiles := shatter.Shatter(element.entity.Bounds(), rocks...)
	for _, tile := range tiles {
		element.entity.DrawBackground(canvas.Cut(destination, tile))
	}
}

func (element *Box) Layout () {
	margin  := element.theme.Margin(tomo.PatternBackground)
	padding := element.theme.Padding(tomo.PatternBackground)
	bounds  := element.entity.Bounds()
	if element.padding { bounds = padding.Apply(bounds) }

	var marginSize float64; if element.vertical {
		marginSize = float64(margin.Y)
	} else {
		marginSize = float64(margin.X)
	}

	freeSpace, nExpanding := element.freeSpace()
	expandingElementSize := freeSpace / nExpanding

	// set the size and position of each element
	x := float64(bounds.Min.X)
	y := float64(bounds.Min.Y)
	for index := 0; index < element.entity.CountChildren(); index ++ {
		entry := element.scratch[element.entity.Child(index)]
		
		var size float64; if entry.expand {
			size = expandingElementSize
		} else {
			size = entry.minSize
		}

		var childBounds image.Rectangle; if element.vertical {
			childBounds = tomo.Bounds(int(x), int(y), bounds.Dx(), int(size))
		} else {
			childBounds = tomo.Bounds(int(x), int(y), int(size), bounds.Dy())
		}
		element.entity.PlaceChild(index, childBounds)

		if element.vertical {
			y += size
			if element.margin { y += marginSize }
		} else {
			x += size
			if element.margin { x += marginSize }
		}
	}
}

func (element *Box) AdoptExpand (children ...tomo.Element) {
	element.adopt(true, children...)
}

func (element *Box) DrawBackground (destination canvas.Canvas) {
	element.entity.DrawBackground(destination)
}

// SetTheme sets the element's theme.
func (element *Box) SetTheme (theme tomo.Theme) {
	if theme == element.theme.Theme { return }
	element.theme.Theme = theme
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

func (element *Box) freeSpace () (space float64, nExpanding float64) {
	margin  := element.theme.Margin(tomo.PatternBackground)
	padding := element.theme.Padding(tomo.PatternBackground)

	var marginSize int; if element.vertical {
		marginSize = margin.Y
	} else {
		marginSize = margin.X
	}
	
	if element.vertical {
		space = float64(element.entity.Bounds().Dy())
	} else {
		space = float64(element.entity.Bounds().Dx())
	}

	for _, entry := range element.scratch {
		if entry.expand {
			nExpanding ++;
		} else {
			space -= float64(entry.minSize)
		}
	}

	if element.padding {
		space -= float64(padding.Vertical())
	}
	if element.margin {
		space -= float64(marginSize * (len(element.scratch) - 1))
	}

	return
}

func (element *Box) updateMinimumSize () {
	margin  := element.theme.Margin(tomo.PatternBackground)
	padding := element.theme.Padding(tomo.PatternBackground)
	var breadth, size int
	var marginSize int; if element.vertical {
		marginSize = margin.Y
	} else {
		marginSize = margin.X
	}
	
	for index := 0; index < element.entity.CountChildren(); index ++ {
		childWidth, childHeight := element.entity.ChildMinimumSize(index)
		var childBreadth, childSize int; if element.vertical {
			childBreadth, childSize = childWidth, childHeight
		} else {
			childBreadth, childSize = childHeight, childWidth
		}
		
		key   := element.entity.Child(index)
		entry := element.scratch[key]
		entry.minSize = float64(childSize)
		element.scratch[key] = entry
		
		if childBreadth > breadth {
			breadth = childBreadth
		}
		size += childSize
		if element.margin && index > 0 {
			size += marginSize
		}
	}

	var width, height int; if element.vertical {
		width, height = breadth, size
	} else {
		width, height = size, breadth
	}

	if element.padding {
		width  += padding.Horizontal()
		height += padding.Vertical()
	}
	
	element.entity.SetMinimumSize(width, height)
}
