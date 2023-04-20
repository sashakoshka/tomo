package elements

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/artist/patterns"

// TODO: this element is lame need to make it better

// Image is an element capable of displaying an image.
type Image struct {
	entity tomo.Entity
	buffer canvas.Canvas
}

// NewImage creates a new image element.
func NewImage (image image.Image) (element *Image) {
	element = &Image { buffer: canvas.FromImage(image) }
	element.entity = tomo.NewEntity(element)
	bounds := element.buffer.Bounds()
	element.entity.SetMinimumSize(bounds.Dx(), bounds.Dy())
	return
}

// Entity returns this element's entity.
func (element *Image) Entity () tomo.Entity {
	return element.entity
}

// Draw causes the element to draw to the specified destination canvas.
func (element *Image) Draw (destination canvas.Canvas) {
	if element.entity == nil { return }
	(patterns.Texture { Canvas: element.buffer }).
		Draw(destination, element.entity.Bounds())
}
