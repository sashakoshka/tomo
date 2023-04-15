package containers

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/shatter"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"

type scratchEntry struct {
	expand  bool
	minimum float64
}

type VBox struct {
	entity  tomo.ContainerEntity
	scratch map[tomo.Element] scratchEntry
	theme   theme.Wrapped
	padding bool
	margin  bool
}

func NewVBox (padding, margin bool) (element *VBox) {
	element = &VBox { padding: padding, margin: margin }
	element.scratch = make(map[tomo.Element] scratchEntry)
	element.theme.Case = tomo.C("tomo", "vBox")
	element.entity = tomo.NewEntity(element).(tomo.ContainerEntity)
	return
}

func (element *VBox) Entity () tomo.Entity {
	return element.entity
}

func (element *VBox) Draw (destination canvas.Canvas) {
	rocks := make([]image.Rectangle, element.entity.CountChildren())
	for index := 0; index < element.entity.CountChildren(); index ++ {
		rocks[index] = element.entity.Child(index).Entity().Bounds()
	}

	tiles := shatter.Shatter(element.entity.Bounds(), rocks...)
	for _, tile := range tiles {
		element.entity.DrawBackground(canvas.Cut(destination, tile))
	}
}

func (element *VBox) Layout () {
	margin  := element.theme.Margin(tomo.PatternBackground)
	padding := element.theme.Padding(tomo.PatternBackground)
	bounds  := element.entity.Bounds()
	if element.padding { bounds = padding.Apply(bounds) }

	freeSpace, nExpanding := element.freeSpace()
	expandingElementHeight := freeSpace / nExpanding

	// set the size and position of each element
	x := float64(bounds.Min.X)
	y := float64(bounds.Min.Y)
	for index := 0; index < element.entity.CountChildren(); index ++ {
		entry := element.scratch[element.entity.Child(index)]
		
		var height float64; if entry.expand {
			height = expandingElementHeight
		} else {
			height = entry.minimum
		}

		element.entity.PlaceChild (index, tomo.Bounds (
			int(x),      int(y),
			bounds.Dx(), int(height)))
			
		y += height
		if element.margin { y += float64(margin.Y) }
	}

}

func (element *VBox) Adopt (child tomo.Element, expand bool) {
	element.entity.Adopt(child)
	element.scratch[child] = scratchEntry { expand: expand }
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

func (element *VBox) Disown (child tomo.Element) {
	index := element.entity.IndexOf(child)
	if index < 0 { return }
	element.entity.Disown(index)
	delete(element.scratch, child)
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

func (element *VBox) DisownAll () {
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

func (element *VBox) HandleChildMinimumSizeChange (child tomo.Element) {
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

func (element *VBox) DrawBackground (destination canvas.Canvas) {
	element.entity.DrawBackground(destination)
}

// SetTheme sets the element's theme.
func (element *VBox) SetTheme (theme tomo.Theme) {
	if theme == element.theme.Theme { return }
	element.theme.Theme = theme
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

func (element *VBox) freeSpace () (space float64, nExpanding float64) {
	margin  := element.theme.Margin(tomo.PatternBackground)
	padding := element.theme.Padding(tomo.PatternBackground)
	space    = float64(element.entity.Bounds().Dy())

	for _, entry := range element.scratch {
		if entry.expand {
			nExpanding ++;
		} else {
			space -= float64(entry.minimum)
		}
	}

	if element.padding {
		space -= float64(padding.Vertical())
	}
	if element.margin {
		space -= float64(margin.Y * (len(element.scratch) - 1))
	}

	return
}

func (element *VBox) updateMinimumSize () {
	margin  := element.theme.Margin(tomo.PatternBackground)
	padding := element.theme.Padding(tomo.PatternBackground)
	var width, height int
	
	for index := 0; index < element.entity.CountChildren(); index ++ {
		childWidth, childHeight :=  element.entity.ChildMinimumSize(index)
		
		key   := element.entity.Child(index)
		entry := element.scratch[key]
		entry.minimum = float64(childHeight)
		element.scratch[key] = entry
		
		if childWidth > width {
			width = childWidth
		}
		height += childHeight
		if element.margin && index > 0 {
			height += margin.Y
		}
	}

	if element.padding {
		width  += padding.Horizontal()
		height += padding.Vertical()
	}
	
	element.entity.SetMinimumSize(width, height)
}
