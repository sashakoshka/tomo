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
	return
}

// Bind binds this element to an entity.
func (element *Image) Bind (entity tomo.Entity) {
	if entity == nil { element.entity = nil; return }
	element.entity = entity
	bounds := element.buffer.Bounds()
	element.entity.SetMinimumSize(bounds.Dx(), bounds.Dy())
}

// Draw causes the element to draw to the specified destination canvas.
func (element *Image) Draw (destination canvas.Canvas) {
	if element.entity == nil { return }
	(patterns.Texture { Canvas: element.buffer }).
		Draw(destination, element.entity.Bounds())
}
