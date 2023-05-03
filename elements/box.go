package elements

import "image"
import "tomo"
import "tomo/artist"
import "tomo/shatter"

var boxCase = tomo.C("tomo", "box") 

// Space is a list of spacing configurations that can be passed to some
// containers.
type Space int; const (
	SpaceNone    Space = 0
	SpacePadding Space = 1
	SpaceMargin  Space = 2
	SpaceBoth    Space = SpacePadding | SpaceMargin
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
	element.entity = tomo.GetBackend().NewEntity(element)
	element.minimumSize = element.updateMinimumSize
	element.init()
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
	element.entity = tomo.GetBackend().NewEntity(element)
	element.minimumSize = element.updateMinimumSize
	element.init()
	element.Adopt(children...)
	return
}

// Draw causes the element to draw to the specified destination canvas.
func (element *Box) Draw (destination artist.Canvas) {
	rocks := make([]image.Rectangle, element.entity.CountChildren())
	for index := 0; index < element.entity.CountChildren(); index ++ {
		rocks[index] = element.entity.Child(index).Entity().Bounds()
	}

	tiles := shatter.Shatter(element.entity.Bounds(), rocks...)
	for _, tile := range tiles {
		element.entity.DrawBackground(artist.Cut(destination, tile))
	}
}

// Layout causes this element to perform a layout operation.
func (element *Box) Layout () {
	margin  := element.entity.Theme().Margin(tomo.PatternBackground, boxCase)
	padding := element.entity.Theme().Padding(tomo.PatternBackground, boxCase)
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

// AdoptExpand adds one or more elements to the box. These elements will be
// expanded to fill in empty space.
func (element *Box) AdoptExpand (children ...tomo.Element) {
	element.adopt(true, children...)
}

// DrawBackground draws this element's background pattern to the specified
// destination canvas.
func (element *Box) DrawBackground (destination artist.Canvas) {
	element.entity.DrawBackground(destination)
}

func (element *Box) HandleThemeChange () {
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

func (element *Box) freeSpace () (space float64, nExpanding float64) {
	margin  := element.entity.Theme().Margin(tomo.PatternBackground, boxCase)
	padding := element.entity.Theme().Padding(tomo.PatternBackground, boxCase)

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
	margin  := element.entity.Theme().Margin(tomo.PatternBackground, boxCase)
	padding := element.entity.Theme().Padding(tomo.PatternBackground, boxCase)
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
